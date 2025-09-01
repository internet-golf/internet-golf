package internetgolf

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"
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
)

type PublicWebServer interface {
	DeployAll([]Deployment) error
}

// implements the PublicWebServer interface
type CaddyServer struct {
	Settings struct {
		LocalOnly bool
	}
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

// creates routes for each of the urls, all of which are handled by handlers.
// this is used to turn groups of urls into separate routes with one matcher
// each, so that the routes can then be sorted in order of matcher specificity.
func urlsToRoutes(urls []Url, requireDomain bool, handlers []json.RawMessage) ([]caddyhttp.Route, error) {
	routes := []caddyhttp.Route{}
	for _, url := range urls {
		if (len(url.Domain) == 0 || !strings.Contains(url.Domain, ".")) && requireDomain {
			return []caddyhttp.Route{}, fmt.Errorf(
				"\"%v\" is not a valid URL: does not start with valid host",
				url,
			)
		}
		matcher := caddy.ModuleMap{}
		if url.Domain != "" {
			matcher["host"] = jsonOrPanic([]string{url.Domain})
		}
		if url.Path != "" {
			matcher["path"] = jsonOrPanic([]string{url.Path})
		}
		routes = append(routes, caddyhttp.Route{
			MatcherSetsRaw: []caddy.ModuleMap{matcher},
			HandlersRaw:    handlers,
		})
	}
	return routes, nil
}

// returns a caddy route that corresponds to a static file server for each URL
// in d.Urls.
func getCaddyStaticRoutes(d Deployment) ([]caddyhttp.Route, error) {
	if d.ServedThingType != StaticFiles {
		return []caddyhttp.Route{}, fmt.Errorf(
			"deployment with name %s passed to getCaddyStaticRoute despite having resource type %s",
			d.Name, d.ServedThingType,
		)
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

	handlers := []json.RawMessage{
		jsonOrPanic(jsonObj{
			"handler": "subroute",
			"routes":  slices.Concat(initialSubroutes, []jsonObj{standardSubroute}),
		}),
	}

	routes, routesErr := urlsToRoutes(d.Urls, true, handlers)
	if routesErr != nil {
		return []caddyhttp.Route{}, routesErr
	}

	return routes, nil
}

// low-level internal deployment type that probably should only be used for the
// admin api
func getCaddyReverseProxyRoute(d Deployment) ([]caddyhttp.Route, error) {
	if d.ServedThingType != ReverseProxy {
		return []caddyhttp.Route{}, fmt.Errorf(
			"deployment with name %s passed to getCaddyReverseProxyRoute despite having resource type %s",
			d.Name, d.ServedThingType,
		)
	}

	handlers := []json.RawMessage{
		jsonOrPanic(jsonObj{
			"handler": "subroute",
			"routes": []jsonObj{
				{
					"handle": []jsonObj{
						// TODO: if this deployment type becomes public, find a
						// way to use the relevant URL from Urls instead of
						// assuming (as this does) that there will only be one
						// URL. either something can be done with variables or
						// `handlers` in urlsToRoutes should become a callback
						{"handler": "rewrite", "strip_path_prefix": d.Urls[0].Path},
						{
							"handler": "reverse_proxy",
							// TODO: someday, control this with a setting? maybe?
							"headers": jsonObj{
								"request": jsonObj{
									"set": jsonObj{
										"Host":      []string{"{http.request.host}"},
										"X-Real-Ip": []string{"{http.request.remote}"},
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
	routes, routesErr := urlsToRoutes(d.Urls, false, handlers)

	if routesErr != nil {
		return []caddyhttp.Route{}, routesErr
	}
	return routes, nil
}

const httpAppServerName = "internetgolf"

// puts all the deployments on the public internet. prioritizes more specific
// urls over less specific urls;
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
					Disabled: !c.Settings.LocalOnly,
				},
				Routes: caddyhttp.RouteList{},
			},
		},
	}

	for _, deployment := range deployments {
		if !deployment.DeploymentContent.HasContent {
			continue
		}

		var getCaddyRoute func(Deployment) ([]caddyhttp.Route, error)

		switch deployment.ServedThingType {
		case StaticFiles:
			getCaddyRoute = getCaddyStaticRoutes
		case ReverseProxy:
			getCaddyRoute = getCaddyReverseProxyRoute
		// TODO: more cases
		default:
			fmt.Printf("could not process deployment with type %s\n", deployment.ServedThingType)
		}

		if routes, err := getCaddyRoute(deployment); err != nil {
			log.Printf("encountered error: %v", err)
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

	caddyConfig := caddy.Config{
		AppsRaw: caddy.ModuleMap{"http": httpJson},
	}

	err = caddy.Run(&caddyConfig)
	if err != nil {
		panic(err)
	}

	return nil
}
