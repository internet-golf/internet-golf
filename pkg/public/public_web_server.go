package public

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"

	"github.com/caddyserver/caddy/v2"
	"github.com/internet-golf/internet-golf/pkg/db"
	"github.com/internet-golf/internet-golf/pkg/resources"

	"github.com/internet-golf/internet-golf/pkg/utils"

	// ??? these modules appear to register themselves with the main caddy
	// module as side effects of being imported. is there a better way to do
	// this?
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/encode"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/encode/gzip"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/encode/zstd"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/fileserver"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/headers"
	_ "github.com/caddyserver/caddy/v2/modules/caddyhttp/reverseproxy"
	_ "github.com/caddyserver/caddy/v2/modules/filestorage"
)

// the primary interface for this whole package
type PublicWebServer interface {
	DeployAll([]db.Deployment) error
	Stop() error
}

// implements the PublicWebServer interface. the primary struct for this whole
// package
type CaddyServer struct {
	config      *utils.Config
	dataPath    string
	onDemandTls onDemandTls
}

const httpAppServerName = "internetgolf"

func NewPublicWebServer(config *utils.Config, files *resources.FileManager) (PublicWebServer, error) {

	return &CaddyServer{config: config, dataPath: files.CaddyDataPath}, nil
}

// puts all the deployments on the public internet. prioritizes more specific
// urls over less specific urls
func (c *CaddyServer) DeployAll(deployments []db.Deployment) error {
	var listen []string
	if c.config.LocalOnly {
		listen = []string{"localhost:80"}
	} else {
		listen = []string{":80", ":443"}
	}
	httpApp := caddyhttp.App{
		Servers: map[string]*caddyhttp.Server{
			httpAppServerName: {
				Listen: listen,
				AutoHTTPS: &caddyhttp.AutoHTTPSConfig{
					Disabled: c.config.LocalOnly,
				},
				Routes: caddyhttp.RouteList{{
					// this matches everything (apparently)
					MatcherSetsRaw: caddyhttp.RawMatcherSets{},
					HandlersRaw: []json.RawMessage{
						utils.JsonOrPanic(utils.JsonObj{
							"handler": "headers",
							"response": map[string]any{
								"add": map[string][]string{
									"X-Deployed-By": []string{"Internet-Golf"},
								},
							},
						}),
					},
				}},
			},
		},
	}

	for _, deployment := range deployments {
		type deploymentToCaddyRouteConverter = func(d db.Deployment) ([]caddyhttp.Route, error)
		var getCaddyRoute deploymentToCaddyRouteConverter

		if !deployment.DeploymentContent.HasContent {
			// if the deployment has no content, substitute in this placeholder
			// content. this is actually load-bearing if not using auto tls,
			// since in that case, caddy will not generate an https cert for
			// this url until it's serving *something*, and until it sets up
			// https, we can't deploy actual content to this url from
			// non-localhost places
			getCaddyRoute = func(d db.Deployment) ([]caddyhttp.Route, error) {
				return GetCaddyTextContentRoute(
					db.Deployment{
						DeploymentMetadata: d.DeploymentMetadata,
					},
				)
			}
		} else {
			switch deployment.ServedThingType {
			case db.StaticFiles:
				getCaddyRoute = GetCaddyStaticRoutes
			case db.ReverseProxy:
				getCaddyRoute = GetCaddyReverseProxyRoute
			default:
				fmt.Printf("could not process deployment with type %s\n", deployment.ServedThingType)
				continue
			}
		}

		if routes, err := getCaddyRoute(deployment); err != nil {
			fmt.Printf("encountered error: %v", err)
		} else {
			httpApp.Servers[httpAppServerName].Routes = append(
				httpApp.Servers[httpAppServerName].Routes,
				routes...,
			)
		}
	}

	// sort more specific routes to the beginning of the slice so that they'll
	// get matched with higher precedence than the less specific routes; i.e.
	// mitch.website/thing needs to be sorted before mitch.website or else
	// mitch.website will always be matched and mitch.website/thing will never
	// be matched
	// TODO: test, with asterisks. need config to get admin api route from (#32)
	slices.SortFunc(
		httpApp.Servers[httpAppServerName].Routes,
		func(a caddyhttp.Route, b caddyhttp.Route) int {
			// catch-all "middleware"
			if len(a.MatcherSetsRaw) == 0 {
				return -1
			}
			if len(b.MatcherSetsRaw) == 0 {
				return 1
			}
			// TODO: make sure admin API route is always first, somehow.
			// terrible hack:
			if string(a.MatcherSetsRaw[0]["path"]) == "/_golf*" {
				return -1
			} else if string(b.MatcherSetsRaw[0]["path"]) == "/_golf*" {
				return 1
			}

			// these routes are guaranteed to have only one matcher set because of
			// how urlsToRoutes works
			if len(a.MatcherSetsRaw[0]["path"]) == 0 && len(b.MatcherSetsRaw[0]["path"]) == 0 {
				// if they both just have a host and no path, then they're equal
				return 0
			} else if len(b.MatcherSetsRaw[0]["path"]) == 0 {
				// if only a has a path, then a is more specific and should be first
				return -1
			} else if len(a.MatcherSetsRaw[0]["path"]) == 0 {
				// if only b has a path, then b is more specific and should be first
				return 1
			} else {
				// otherwise, assume the longer path is more specific. which i think
				// will give good results?
				// TODO: account for asterisks? needs testing
				return len(b.MatcherSetsRaw[0]["path"]) - len(a.MatcherSetsRaw[0]["path"])
			}
		},
	)

	httpJson, err := json.Marshal(httpApp)
	if err != nil {
		panic(err)
	}

	// TODO: like the auto tls server setup, the admin api port stuff is called
	// every time DeployAll is called, when it should be more like an initial
	// setup thing

	// starting the caddy admin api at a random port that is only known within
	// this program might make it slightly harder to reach and exploit 🤞
	caddyAdminApiPort, _ := utils.GetFreePort()
	logLevel := "ERROR"
	if c.config.Verbose {
		logLevel = "DEBUG"
	}

	fmt.Printf("Caddy API running at localhost:%v\n", caddyAdminApiPort)

	caddyConfig := caddy.Config{
		AppsRaw: caddy.ModuleMap{"http": httpJson},
		Admin: &caddy.AdminConfig{
			Listen: "localhost:" + strconv.Itoa(caddyAdminApiPort),
		},
		StorageRaw: utils.JsonOrPanic(map[string]string{
			"module": "file_system",
			"root":   c.dataPath,
		}),
		Logging: &caddy.Logging{
			Logs: map[string]*caddy.CustomLog{
				"default": &caddy.CustomLog{
					BaseLog: caddy.BaseLog{Level: logLevel},
				},
			},
		},
	}

	// TODO: wait, how is the tlsApprovalServer set up every time that DeployAll
	// is called? shouldn't this be in NewPublicWebServer?
	if !c.config.LocalOnly {
		tlsConfig, err := getOnDemandTls()
		if err != nil {
			panic(err)
		}
		c.onDemandTls = tlsConfig
		caddyConfig.AppsRaw["tls"] = utils.JsonOrPanic(tlsConfig.caddyTlsConfig)
		go c.onDemandTls.tlsApprovalServer.ListenAndServe()
	}

	err = caddy.Run(&caddyConfig)
	if err != nil {
		panic(err)
	}

	return nil
}

func (c *CaddyServer) Stop() error {
	if !c.config.LocalOnly {
		err := c.onDemandTls.tlsApprovalServer.Shutdown(context.TODO())
		if err != nil {
			caddy.Stop()
			return err
		}
	}
	return caddy.Stop()
}
