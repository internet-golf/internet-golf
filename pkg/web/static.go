package web

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/toBeOfUse/internet-golf/pkg/db"
	"github.com/toBeOfUse/internet-golf/pkg/utils"
)

// returns a slice of caddy Route struct instances: one caddy route that
// corresponds to a static file server for each URL in d.Urls.
func GetCaddyStaticRoutes(d db.Deployment) ([]caddyhttp.Route, error) {
	if d.ServedThingType != db.StaticFiles {
		return []caddyhttp.Route{}, fmt.Errorf(
			"deployment with URL %s passed to getCaddyStaticRoute despite having resource type %s",
			d.Url, d.ServedThingType,
		)
	}

	routes := []caddyhttp.Route{
		{
			MatcherSetsRaw: caddyhttp.RawMatcherSets{
				// {"host": jsonOrPanic([]string{""})},
			},
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
		},
	}

	// TODO: control the "makePathCatchAll" argument with setting on deployment
	matcher, matcherErr := urlToMatcher(d.Url, false, true)
	if matcherErr != nil {
		return nil, matcherErr
	}

	standardSubroute := utils.JsonObj{
		"handle": []utils.JsonObj{
			{
				"handler": "vars",
				"root":    d.ServedThing,
			},
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

	var initialSubroutes []utils.JsonObj
	if !d.DeploymentMetadata.PreserveExternalPath {
		cleanPath, _ := strings.CutSuffix(d.Url.Path, "*")
		initialSubroutes = append(initialSubroutes,
			// TODO: does this work with asterisks?
			utils.JsonObj{
				"handle": []utils.JsonObj{
					utils.JsonObj{"handler": "rewrite", "strip_path_prefix": cleanPath},
				},
			},
		)
	}

	handlers := []json.RawMessage{
		utils.JsonOrPanic(utils.JsonObj{
			"handler": "subroute",
			"routes":  slices.Concat(initialSubroutes, []utils.JsonObj{standardSubroute}),
		}),
	}

	routes = append(routes, caddyhttp.Route{
		MatcherSetsRaw: caddyhttp.RawMatcherSets{matcher},
		HandlersRaw:    handlers,
	})

	return routes, nil
}
