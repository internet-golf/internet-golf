package handlers

import (
	"fmt"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/toBeOfUse/internet-golf/pkg/db"
	"github.com/toBeOfUse/internet-golf/pkg/utils"
)

// TODO: remove requireDomain argument. if enforced, that should be validated at
// the api call/deployment creation level
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

type Handler = func(d db.Deployment) ([]caddyhttp.Route, error)
