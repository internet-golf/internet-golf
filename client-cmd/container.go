package main

import (
	"context"
	"fmt"
	"io"

	golfsdk "github.com/internet-golf/internet-golf/client-sdk"
	"github.com/spf13/cobra"
)

func deployContainerCommand() *cobra.Command {
	var name string
	var registry string

	deployContent := cobra.Command{
		Use:     "deploy-container [deployment-name]",
		Example: "deploy-container example.com --name my-container --registry docker.io",
		Short:   "Deploys container",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := createClient(args[0])

			body, resp, respError := client.
				DefaultAPI.PutDeployContainer(context.TODO()).
				DeployContainerInputBody(golfsdk.DeployContainerInputBody{
					Url:         args[0],
					ImageName:   name,
					RegistryUrl: registry,
				}).
				Execute()

			// TODO: handle responses uniformly across commands
			if respError != nil {
				panic(respError.Error())
			}

			if resp.StatusCode != 200 {
				body, _ := io.ReadAll(resp.Body)
				panic("[error from server]: " + string(body))
			}
			if body == nil || !body.Success {
				panic("Did not get success status back from server. Request was to " + resp.Request.URL.String())
			}

			fmt.Println(body.Message)
		},
	}

	deployContent.Flags().StringVar(
		&name, "name", "",
		"Supply the name of the image you wish to deploy.",
	)

	deployContent.Flags().StringVar(
		&registry, "registry", "",
		"Supply the registry of the image you wish to deploy.",
	)

	return &deployContent
}
