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

func urlToMatcher(url string, requireHost bool) (caddyhttp.RawMatcherSets, error) {
	matcherComps := strings.Split(url, "/")
	host := matcherComps[0]
	if (len(host) == 0 || !strings.Contains(host, ".")) && requireHost {
		return caddyhttp.RawMatcherSets{}, fmt.Errorf(
			"\"%v\" is not a valid URL: does not start with valid host",
			url,
		)
	}
	var path string
	if len(matcherComps) == 1 || len(matcherComps[1]) == 0 {
		path = ""
	} else {
		path = "/" + strings.Join(matcherComps[1:], "/")
	}
	matcher := caddyhttp.RawMatcherSets{
		caddy.ModuleMap{
			"host": jsonOrPanic([]string{host}),
		},
	}

	if path != "" {
		matcher[0]["path"] = jsonOrPanic([]string{path})
	}

	return matcher, nil

}

// returns a caddy route that corresponds to a static file server for the
// Deployment d.
//
// this function is composed mainly of terrifying json soup but it's unclear how
// else to do it since caddy.Run() expects json-serializable stuff (and the doc
// comment for the caddy Config struct says "Go programs which are directly
// building a Config struct value should take care to populate the
// JSON-encodable fields of the struct")
func getCaddyStaticRoute(d Deployment) (caddyhttp.Route, error) {
	if d.ServedThingType != StaticFiles {
		return caddyhttp.Route{}, fmt.Errorf(
			"deployment with URL %s passed to getCaddyStaticRoute despite having resource type %s",
			d.Url, d.ServedThingType,
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

	// TODO: #2
	// if d.Settings.CacheControlMode != Default {
	// 	var cacheControlMatcher []string
	// 	switch d.Settings.CacheControlMode {
	// 	case AllButHtml:
	// 		// */ is supposed to match index routes (i.e. those ending in /,
	// 		// like thing.com/whatever/)
	// 		cacheControlMatcher = []string{"*/", "*.html"}
	// 	case Nothing:
	// 		cacheControlMatcher = []string{"*"}
	// 	}
	// 	initialSubroutes = append(initialSubroutes, jsonObj{
	// 		"match": []jsonObj{
	// 			{"path": cacheControlMatcher},
	// 		},
	// 		"handle": []jsonObj{
	// 			{
	// 				"handler": "headers",
	// 				"response": jsonObj{
	// 					"set": jsonObj{
	// 						"Cache-Control": []string{"max-age=0,no-store"},
	// 					},
	// 				}},
	// 		},
	// 	})
	// }

	matcher, matcherErr := urlToMatcher(d.Url, true)
	if matcherErr != nil {
		return caddyhttp.Route{}, matcherErr
	}

	route := caddyhttp.Route{
		MatcherSetsRaw: matcher,
		HandlersRaw: []json.RawMessage{
			jsonOrPanic(jsonObj{
				"handler": "subroute",
				"routes":  slices.Concat(initialSubroutes, []jsonObj{standardSubroute}),
			}),
		},
	}

	return route, nil
}

// low-level internal deployment type that probably should only be used for the
// admin api
func getCaddyReverseProxyRoute(d Deployment) (caddyhttp.Route, error) {
	if d.ServedThingType != ReverseProxy {
		return caddyhttp.Route{}, fmt.Errorf(
			"deployment with URL %s passed to getCaddyReverseProxyRoute despite having resource type %s",
			d.Url, d.ServedThingType,
		)
	}

	// not requiring a host here bc this deployment type is for meta-deployments
	matcher, matcherErr := urlToMatcher(d.Url, false)
	if matcherErr != nil {
		return caddyhttp.Route{}, matcherErr
	}

	return caddyhttp.Route{
		MatcherSetsRaw: matcher,
		HandlersRaw: []json.RawMessage{
			jsonOrPanic(jsonObj{
				"handler": "subroute",
				"routes": []jsonObj{
					{
						"handle": []jsonObj{
							// TODO: control strip_path_prefix with a setting
							{"handler": "rewrite", "strip_path_prefix": d.Url},
							{
								"handler": "reverse_proxy",
								// TODO: control this with a setting? maybe?
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
		},
	}, nil
}

func (c *CaddyServer) DeployAll(deployments []Deployment) error {
	var listen []string
	if c.Settings.LocalOnly {
		listen = []string{"localhost:80"}
	} else {
		listen = []string{":80", ":443"}
	}
	httpApp := caddyhttp.App{
		Servers: map[string]*caddyhttp.Server{
			"internetgolf": {
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

		var getCaddyRoute func(Deployment) (caddyhttp.Route, error)

		switch deployment.ServedThingType {
		case StaticFiles:
			getCaddyRoute = getCaddyStaticRoute
		case ReverseProxy:
			getCaddyRoute = getCaddyReverseProxyRoute
		default:
			fmt.Printf("could not process deployment with type %s\n", deployment.ServedThingType)
		}

		if route, err := getCaddyRoute(deployment); err != nil {
			log.Printf("encountered error: %v", err)
		} else {
			httpApp.Servers["internetgolf"].Routes = append(
				httpApp.Servers["internetgolf"].Routes,
				route,
			)
		}
	}

	// TODO: docker cases

	httpJson, err := json.Marshal(httpApp)
	if err != nil {
		panic(err)
	}

	caddyConfig := caddy.Config{
		AppsRaw: caddy.ModuleMap{
			"http": httpJson,
		},
	}

	err = caddy.Run(&caddyConfig)
	if err != nil {
		panic(err)
	}

	return nil
}
