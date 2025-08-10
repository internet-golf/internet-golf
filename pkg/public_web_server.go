package internetgolf

import (
	"encoding/json"
	"fmt"
	"log"
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

func getCaddyStaticRoute(d Deployment) (caddyhttp.Route, error) {
	// decompose matcher into host and path
	matcherComps := strings.Split(d.Matcher, "/")
	host := matcherComps[0]
	if len(host) == 0 || !strings.Contains(host, ".") {
		return caddyhttp.Route{}, fmt.Errorf("\"%v\" is not a valid matcher: does not start with valid host", d.Matcher)
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

	return route, nil
}

func (c CaddyServer) DeployAll(deployments []Deployment) error {
	httpApp := caddyhttp.App{
		Servers: map[string]*caddyhttp.Server{
			"internetgolf": {
				Listen: []string{"localhost:80"},
				AutoHTTPS: &caddyhttp.AutoHTTPSConfig{
					// just for local testing...
					Disabled: true,
				},
				Routes: caddyhttp.RouteList{},
			},
		},
	}

	for _, deployment := range deployments {
		if deployment.LocalResourceType == Files {
			if route, err := getCaddyStaticRoute(deployment); err != nil {
				log.Printf("encountered error: %v", err)
			} else {
				httpApp.Servers["internetgolf"].Routes = append(
					httpApp.Servers["internetgolf"].Routes,
					route,
				)
			}
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
