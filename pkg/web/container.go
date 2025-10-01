package web

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/toBeOfUse/internet-golf/pkg/db"
	"github.com/toBeOfUse/internet-golf/pkg/utils"
)

// deployment type for containers, should pull the image, run the container and
// proxy it behind Caddy
func GetCaddyContainerRoute(d db.Deployment) ([]caddyhttp.Route, error) {
	if d.ServedThingType != db.DockerContainer {
		return []caddyhttp.Route{}, fmt.Errorf(
			"deployment with name %s passed to getCaddyContainerRoute despite having resource type %s",
			d.Url, d.ServedThingType,
		)
	}

	cli, err := client.NewClientWithOpts()
	if err != nil {
		fmt.Println("Unable to create docker client")
		return []caddyhttp.Route{}, err
	}

	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return []caddyhttp.Route{}, cwdErr
	}

	// TODO: move this kind of stuff to content/containers.go
	ctx := context.Background()
	cont, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:        "openapitools/openapi-generator-cli:latest",
			AttachStdout: false,
			AttachStderr: false,
			Cmd: []string{
				"generate", "-i", "/local/golf-openapi.yaml", "-g", "go", "-o", "/local/client-sdk",
				"--additional-properties=packageName=golfsdk,withGoMod=false",
			},
			// TODO: what is new(int) doing here
			StopTimeout:     new(int),
			NetworkDisabled: true,
		},
		&container.HostConfig{
			Binds: []string{
				cwd + ":/local",
			},
		},
		nil,
		nil,
		"",
	)
	if err != nil {
		return []caddyhttp.Route{}, err
	}

	if err := cli.ContainerStart(ctx, cont.ID, container.StartOptions{}); err != nil {
		return []caddyhttp.Route{}, err
	}

	statusCh, errCh := cli.ContainerWait(ctx, cont.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return []caddyhttp.Route{}, err
		}
	case <-statusCh:
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
										"Host":      []string{"{http.request.host}"},
										"X-Real-Ip": []string{"{http.request.remote}"},
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
