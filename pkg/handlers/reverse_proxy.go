package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/toBeOfUse/internet-golf/pkg/db"
	"github.com/toBeOfUse/internet-golf/pkg/utils"
)

// low-level internal deployment type that probably should only be used for the
// admin api
func GetCaddyReverseProxyRoute(d db.Deployment) ([]caddyhttp.Route, error) {
	if d.ServedThingType != db.ReverseProxy {
		return []caddyhttp.Route{}, fmt.Errorf(
			"deployment with name %s passed to getCaddyReverseProxyRoute despite having resource type %s",
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
