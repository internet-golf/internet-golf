package public

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/internet-golf/internet-golf/pkg/db"
	"github.com/internet-golf/internet-golf/pkg/utils"
)

func GetCaddyReverseProxyRoute(d db.Deployment) ([]caddyhttp.Route, error) {
	if d.ServedThingType != db.ReverseProxy {
		return []caddyhttp.Route{}, fmt.Errorf(
			"deployment with name %s passed to "+
				"getCaddyReverseProxyRoute despite having resource type %s",
			d.Url, d.ServedThingType,
		)
	}

	handlers := []json.RawMessage{
		utils.JsonOrPanic(utils.JsonObj{
			"handler": "subroute",
			"routes": []utils.JsonObj{
				{
					"handle": []utils.JsonObj{
						{"handler": "rewrite", "strip_path_prefix": d.Url.Path},
						{
							"handler": "reverse_proxy",
							// TODO: someday, control this with a setting? maybe?
							"headers": utils.JsonObj{
								"request": utils.JsonObj{
									"set": utils.JsonObj{
										"Host":            []string{"{http.request.host}"},
										"X-Forwarded-For": []string{"{http.request.remote}"},
									},
								},
							},
							"upstreams": []utils.JsonObj{{"dial": d.ServedThing}},
						},
					},
				},
			},
		}),
	}

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

// returns a slice of caddy Route struct instances: one caddy route that
// corresponds to a static file server for each URL in d.Urls.
func GetCaddyStaticRoutes(d db.Deployment) ([]caddyhttp.Route, error) {
	if d.ServedThingType != db.StaticFiles {
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

	initialSubroutes := []utils.JsonObj{
		utils.JsonObj{
			"handle": []utils.JsonObj{
				{
					"handler": "vars",
					"root":    d.ServedThing,
				},
			},
		},
	}

	cleanPath, _ := strings.CutSuffix(d.Url.Path, "*")
	if len(cleanPath) > 0 && !d.DeploymentMetadata.PreserveExternalPath {
		initialSubroutes = append(initialSubroutes,
			utils.JsonObj{
				"handle": []utils.JsonObj{
					utils.JsonObj{"handler": "rewrite", "strip_path_prefix": cleanPath},
				},
			},
		)
	}

	if d.DeploymentContent.SpaMode {
		initialSubroutes = append(initialSubroutes,
			utils.JsonObj{
				"handle": []utils.JsonObj{
					{

						"handler": "rewrite", "uri": "{http.matchers.file.relative}",
					},
				},
				"match": []utils.JsonObj{
					utils.JsonObj{
						"file": utils.JsonObj{"try_files": []string{"{http.request.uri.path}", "/index.html"}},
					},
				},
			},
		)
	}

	finalSubroute := utils.JsonObj{
		"handle": []utils.JsonObj{
			{
				"handler": "encode",
				"encodings": utils.JsonObj{
					"gzip": utils.JsonObj{},
					"zstd": utils.JsonObj{},
				},
				"prefer": []string{"zstd", "gzip"},
			},
			{
				"handler": "file_server",
			},
		},
	}

	handlers := []json.RawMessage{
		utils.JsonOrPanic(utils.JsonObj{
			"handler": "subroute",
			"routes":  slices.Concat(initialSubroutes, []utils.JsonObj{finalSubroute}),
		}),
	}

	routes = append(routes, caddyhttp.Route{
		MatcherSetsRaw: caddyhttp.RawMatcherSets{matcher},
		HandlersRaw:    handlers,
	})

	return routes, nil
}

// get a route that will respond with basic text content. this does not look at
// anything in the deployment that's passed in except the URL.
func GetCaddyTextContentRoute(d db.Deployment) ([]caddyhttp.Route, error) {
	matcher, matcherErr := urlToMatcher(d.Url, false, true)
	if matcherErr != nil {
		return []caddyhttp.Route{}, matcherErr
	}

	return []caddyhttp.Route{{
		MatcherSetsRaw: caddyhttp.RawMatcherSets{matcher},
		HandlersRaw: []json.RawMessage{
			utils.JsonOrPanic(utils.JsonObj{
				"handler":     "static_response",
				"status_code": 200,
				"body":        "server initialized",
			}),
		},
	}}, nil
}

// TODO: remove requireDomain argument. if enforced, that should be validated at
// the api call/deployment creation level
// utility function used by route creator functions above
func urlToMatcher(url db.Url, requireDomain bool, makePathCatchAll bool) (caddy.ModuleMap, error) {
	if (len(url.Domain) == 0 || !strings.Contains(url.Domain, ".")) && requireDomain {
		return caddy.ModuleMap{}, fmt.Errorf(
			"\"%v\" is not a valid URL: does not start with valid host",
			url,
		)
	}
	matcher := caddy.ModuleMap{}
	if url.Domain != "" {
		matcher["host"] = utils.JsonOrPanic([]string{url.Domain})
	}
	if url.Path != "" {
		if makePathCatchAll {
			url.Path += "*"
		}
		matcher["path"] = utils.JsonOrPanic([]string{url.Path})
	}

	return matcher, nil
}
