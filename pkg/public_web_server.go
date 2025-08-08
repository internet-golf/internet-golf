package internetgolf

import (
	"encoding/json"
	"fmt"
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
)

type PublicWebServer interface {
	Deploy([]Deployment)
}

type CaddyServer struct{}

// TODO: receive the whole deployment object for to take into account its settings
func getCaddyStaticRouteBoilerplate(publicUrl string, filePath string) caddyhttp.Route {
	route := caddyhttp.Route{
		MatcherSetsRaw: caddyhttp.RawMatcherSets{
			caddy.ModuleMap{"host": ([]byte)(`["` + publicUrl + `"]`)},
		},
		// TODO: handlers
	}

	return route

	// return fmt.Sprintf(`{
	// 	"match": [{ "host": ["%s"] }],
	// 	"handle": [
	// 	{
	// 		"handler": "subroute",
	// 		"routes": [
	// 		{
	// 			"handle": [
	// 			{
	// 				"handler": "vars",
	// 				"root": "%s"
	// 			}
	// 			]
	// 		},
	// 		{
	// 			"handle": [
	// 			{
	// 				"handler": "headers",
	// 				"response": { "set": { "Cache-Control": ["max-age=0,no-store"] } }
	// 			}
	// 			],
	// 			"match": [{ "path": ["*/"] }]
	// 		},
	// 		{
	// 			"handle": [
	// 			{
	// 				"handler": "encode",
	// 				"encodings": { "gzip": {}, "zstd": {} },
	// 				"prefer": ["zstd", "gzip"]
	// 			},
	// 			{"browse": {"file_limit": 1000}, "handler": "file_server" }
	// 			]
	// 		}
	// 		]
	// 	}
	// 	],
	// 	"terminal": true
	// }`, publicUrl, filePath)
}

func (c CaddyServer) Deploy(deployments []Deployment) error {
	fmt.Printf("deploying %+v\n", deployments)
	// TODO: implement https://caddyserver.com/docs/api#concurrent-config-changes

	httpApp := caddyhttp.App{
		Servers: map[string]*caddyhttp.Server{
			"internetgolf": {
				Listen: []string{":80", ":443"},
				AutoHTTPS: &caddyhttp.AutoHTTPSConfig{
					// just for local testing...
					Disabled: true,
				},
				Routes: caddyhttp.RouteList{},
			},
		},
	}
	for _, deployment := range deployments {
		uriParts := strings.Split(deployment.ResourceUri, "://")
		if uriParts[0] == "file" {
			httpApp.Servers["internetgolf"].Routes = append(
				httpApp.Servers["internetgolf"].Routes,
				getCaddyStaticRouteBoilerplate(deployment.Matcher, uriParts[1]),
			)
		}
		// TODO: docker cases
	}

	httpJson, err := json.Marshal(httpApp)
	if err != nil {
		panic(err)
	}

	caddyConfig := caddy.Config{AppsRaw: caddy.ModuleMap{
		"http": httpJson,
	}}

	err = caddy.Run(&caddyConfig)
	if err != nil {
		panic(err)
	}

	return nil
}
