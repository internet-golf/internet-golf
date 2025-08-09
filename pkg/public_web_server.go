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

func jsonOrPanic(v any) []byte {
	result, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("Could not JSON-serialize value: %v", v))
	}
	return result
}

func getCaddyStaticRoute(d Deployment) caddyhttp.Route {
	// decompose matcher into host and path
	matcherComps := strings.Split(d.Matcher, "/")
	host := matcherComps[0]
	if len(host) == 0 || !strings.Contains(host, ".") {
		panic(fmt.Sprintf("%v is not a valid matcher: does not start with valid host", d.Matcher))
	}
	var path string
	if len(matcherComps) == 1 || len(matcherComps[1]) == 0 {
		path = "/*"
	} else {
		path = "/" + strings.Join(matcherComps[1:], "/")
	}

	route := caddyhttp.Route{
		MatcherSetsRaw: caddyhttp.RawMatcherSets{
			caddy.ModuleMap{"host": jsonOrPanic([]string{host})},
			caddy.ModuleMap{"path": jsonOrPanic([]string{path})},
		},
		HandlersRaw: []json.RawMessage{
			[]byte(fmt.Sprintf(`{"handler": "vars", "root": "%s"}`, d.LocalResourceLocator)),
			[]byte(`{"handler": "file_server"}`),
		},
	}

	return route
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
		// TODO: recover from panics and ignore the deployment with a warning
		if deployment.LocalResourceType == Files {
			httpApp.Servers["internetgolf"].Routes = append(
				httpApp.Servers["internetgolf"].Routes,
				getCaddyStaticRoute(deployment),
			)
		}
		// TODO: docker cases
	}

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
