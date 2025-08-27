// this is meant to be called from a github action with `go run misc/oidc_sink.go eyJsdfkjsdfkalj`

package main

import (
	"fmt"
	"os"

	golf "github.com/toBeOfUse/internet-golf/pkg"
)

func main() {
	oidcArg := os.Args[1]
	result, err := golf.ParseGithubOidcToken(oidcArg)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%+v\n", result)
	}
}
