package internetgolf

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"slices"
	"strconv"
	"strings"

	"github.com/caddyserver/caddy/v2"
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

type PublicWebServer interface {
	Init() error
	DeployAll([]Deployment) error
	Stop() error
}

// implements the PublicWebServer interface
type CaddyServer struct {
	Settings struct {
		LocalOnly bool
		Verbose   bool
	}
	StorageSettings        StorageSettings
	placeholderContentPath string
}

// this is apparently how you have to do this
type jsonObj map[string]any

// utility function to turn a value into json without possibly returning an
// error. should only really be used if it seems incredibly unlikely that
// json.Marshal will panic when given v.
func jsonOrPanic(v any) []byte {
	result, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("Could not JSON-serialize value: %v", v))
	}
	return result
}

// TODO: remove requireDomain argument. if enforced, that should be validated at
// the api call/deployment creation level
func urlToMatcher(url Url, requireDomain bool, makePathCatchAll bool) (caddy.ModuleMap, error) {
	if (len(url.Domain) == 0 || !strings.Contains(url.Domain, ".")) && requireDomain {
		return caddy.ModuleMap{}, fmt.Errorf(
			"\"%v\" is not a valid URL: does not start with valid host",
			url,
		)
	}
	matcher := caddy.ModuleMap{}
	if url.Domain != "" {
		matcher["host"] = jsonOrPanic([]string{url.Domain})
	}
	if url.Path != "" {
		if makePathCatchAll {
			url.Path += "*"
		}
		matcher["path"] = jsonOrPanic([]string{url.Path})
	}

	return matcher, nil
}

// returns a slice of caddy Route struct instances: one caddy route that
// corresponds to a static file server for each URL in d.Urls.
func getCaddyStaticRoutes(d Deployment) ([]caddyhttp.Route, error) {
	if d.ServedThingType != StaticFiles {
		return []caddyhttp.Route{}, fmt.Errorf(
			"deployment with URL %s passed to getCaddyStaticRoute despite having resource type %s",
			d.Url, d.ServedThingType,
		)
	}

	routes := []caddyhttp.Route{}

	// TODO: control the "makePathCatchAll" argument with setting on deployment
	matcher, matcherErr := urlToMatcher(d.Url, false, true)
	if matcherErr != nil {
		return nil, matcherErr
	}

	standardSubroute := jsonObj{
		"handle": []jsonObj{
			{
				"handler": "vars",
				"root":    d.ServedThing,
			},
			{
				"handler": "encode",
				"encodings": jsonObj{
					"gzip": jsonObj{},
					"zstd": jsonObj{},
				},
				"prefer": []string{"zstd", "gzip"},
			},
			{
				"handler": "file_server",
			},
		},
	}

	var initialSubroutes []jsonObj
	if !d.DeploymentMetadata.PreserveExternalPath {
		cleanPath, _ := strings.CutSuffix(d.Url.Path, "*")
		initialSubroutes = append(initialSubroutes,
			// TODO: does this work with asterisks?
			jsonObj{
				"handle": []jsonObj{
					jsonObj{"handler": "rewrite", "strip_path_prefix": cleanPath},
				},
			},
		)
	}

	handlers := []json.RawMessage{
		jsonOrPanic(jsonObj{
			"handler": "subroute",
			"routes":  slices.Concat(initialSubroutes, []jsonObj{standardSubroute}),
		}),
	}

	routes = append(routes, caddyhttp.Route{
		MatcherSetsRaw: caddyhttp.RawMatcherSets{matcher},
		HandlersRaw:    handlers,
	})

	return routes, nil
}

// low-level internal deployment type that probably should only be used for the
// admin api
func getCaddyReverseProxyRoute(d Deployment) ([]caddyhttp.Route, error) {
	if d.ServedThingType != ReverseProxy {
		return []caddyhttp.Route{}, fmt.Errorf(
			"deployment with name %s passed to getCaddyReverseProxyRoute despite having resource type %s",
			d.Url, d.ServedThingType,
		)
	}

	handlers := []json.RawMessage{
		jsonOrPanic(jsonObj{
			"handler": "subroute",
			"routes": []jsonObj{
				{
					"handle": []jsonObj{
						{"handler": "rewrite", "strip_path_prefix": d.Url.Path},
						{
							"handler": "reverse_proxy",
							// TODO: someday, control this with a setting? maybe?
							"headers": jsonObj{
								"request": jsonObj{
									"set": jsonObj{
										"Host":            []string{"{http.request.host}"},
										"X-Forwarded-For": []string{"{http.request.remote}"},
									},
								},
							},
							"upstreams": []jsonObj{{"dial": d.ServedThing}},
						},
					},
				},
			},
		}),
	}

	// not requiring a host here bc this deployment type is for meta-deployments
	matcher, matcherErr := urlToMatcher(d.Url, false, true)

	if matcherErr != nil {
		return []caddyhttp.Route{}, matcherErr
	}
	return []caddyhttp.Route{
		caddyhttp.Route{
			MatcherSetsRaw: caddyhttp.RawMatcherSets{matcher},
			HandlersRaw:    handlers,
		},
	}, nil
}

const httpAppServerName = "internetgolf"

//go:embed content/placeholder.html
var placeholderContent []byte

func (c *CaddyServer) Init() error {
	c.placeholderContentPath = path.Join(c.StorageSettings.DataDirectory, "placeholder-content")
	os.MkdirAll(c.placeholderContentPath, 0644)

	return os.WriteFile(path.Join(c.placeholderContentPath, "index.html"), placeholderContent, 0644)

	// caddy seems to start itself?
}

// puts all the deployments on the public internet. prioritizes more specific
// urls over less specific urls
func (c *CaddyServer) DeployAll(deployments []Deployment) error {
	var listen []string
	if c.Settings.LocalOnly {
		listen = []string{"localhost:80"}
	} else {
		listen = []string{":80", ":443"}
	}
	httpApp := caddyhttp.App{
		Servers: map[string]*caddyhttp.Server{
			httpAppServerName: {
				Listen: listen,
				AutoHTTPS: &caddyhttp.AutoHTTPSConfig{
					Disabled: c.Settings.LocalOnly,
				},
				Routes: caddyhttp.RouteList{{
					// match all
					MatcherSetsRaw: caddyhttp.RawMatcherSets{},
					HandlersRaw: []json.RawMessage{
						jsonOrPanic(jsonObj{
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
		var getCaddyRoute func(Deployment) ([]caddyhttp.Route, error)

		if !deployment.DeploymentContent.HasContent {
			getCaddyRoute = func(d Deployment) ([]caddyhttp.Route, error) {
				return getCaddyStaticRoutes(
					Deployment{
						DeploymentMetadata: d.DeploymentMetadata,
						DeploymentContent: DeploymentContent{
							HasContent:      true,
							ServedThingType: StaticFiles,
							ServedThing:     c.placeholderContentPath,
						},
					},
				)
			}
		} else {
			switch deployment.ServedThingType {
			case StaticFiles:
				getCaddyRoute = getCaddyStaticRoutes
			case ReverseProxy:
				getCaddyRoute = getCaddyReverseProxyRoute
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
	slices.SortFunc(
		httpApp.Servers[httpAppServerName].Routes,
		func(a caddyhttp.Route, b caddyhttp.Route) int {
			// catch-all "middleware"
			if len(a.MatcherSetsRaw) == 0 {
				return -1
			}
			if len(b.MatcherSets) == 0 {
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

	// starting the caddy admin api at a random port that is only known within
	// this program might make it slightly harder to reach and exploit 🤞
	caddyAdminApiPort, _ := GetFreePort()
	logLevel := "ERROR"
	if c.Settings.Verbose {
		logLevel = "DEBUG"
	}
	caddyConfig := caddy.Config{
		AppsRaw: caddy.ModuleMap{"http": httpJson},
		Admin: &caddy.AdminConfig{
			Listen: "localhost:" + strconv.Itoa(caddyAdminApiPort),
		},
		StorageRaw: jsonOrPanic(map[string]string{
			"module": "file_system",
			"root":   path.Join(c.StorageSettings.DataDirectory, "caddy"),
		}),
		Logging: &caddy.Logging{
			Logs: map[string]*caddy.CustomLog{
				"default": &caddy.CustomLog{
					BaseLog: caddy.BaseLog{Level: logLevel},
				},
			},
		},
	}

	err = caddy.Run(&caddyConfig)
	if err != nil {
		panic(err)
	}

	return nil
}

func (c *CaddyServer) Stop() error {
	return caddy.Stop()
}
