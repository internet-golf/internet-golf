package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "golf",
	Short: "A server",
}

func createDeploymentCommand() *cobra.Command {

	var github string

	createDeployment := cobra.Command{
		Use:     "create-deployment [deployment-url]",
		Example: "create-deployment example.com --github repoOwner/repoName",
		Short:   "Creates a deployment",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hello i am creating a deployment")
			fmt.Printf("the deployment is for the url %s\n", args[0])
			fmt.Printf("the github arg is %v\n", github)

			var externalSourceType string
			var externalSource string

			if len(github) > 0 {
				// TODO: would be nice to use ExternalSourceType somehow
				externalSourceType = "GithubRepo"
				externalSource = github
			}
			client, clientErr := NewClientWithResponses("http://localhost:8888")
			if clientErr != nil {
				panic(clientErr.Error())
			}
			resp, err := client.PostDeployNewWithResponse(
				context.TODO(), PostDeployNewJSONRequestBody{
					Url:                args[0],
					ExternalSourceType: &externalSourceType,
					ExternalSource:     &externalSource,
				},
			)
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("success: %v\n", resp.JSON200.Success)
		},
	}

	createDeployment.Flags().StringVar(
		&github, "github", "", "Specify a Github Repo: repoOwner/repoName",
	)

	return &createDeployment
}

func main() {
	rootCmd.AddCommand(createDeploymentCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
