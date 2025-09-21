package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mholt/archives"
	"github.com/spf13/cobra"
	golfsdk "github.com/toBeOfUse/internet-golf/client-sdk"
)

var rootCmd = &cobra.Command{
	Use:   "golf",
	Short: "Deploy stuff to a server",
}

// TODO: should all the subcommands run the health check and wait for it to work
// for a few seconds?

var apiUrl string

func createClient() *golfsdk.APIClient {
	return golfsdk.NewAPIClient(&golfsdk.Configuration{
		UserAgent: "InternetGolfClient",
		Servers: golfsdk.ServerConfigurations{
			{URL: apiUrl},
		},
	})
}

func createDeploymentCommand() *cobra.Command {

	var github string
	var path string
	var name string

	createDeployment := cobra.Command{
		Use:     "create-deployment domain",
		Example: "create-deployment example.com --github repoOwner/repoName",
		Short:   "Creates a deployment",
		// TODO: actually, allow multiple domains to be specified (but then how
		// to give them paths? maybe the format should be example.com/path for
		// the args)
		Args: cobra.ExactArgs(1),
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

			client := createClient()
			fmt.Printf("created client with config %+v\n", client.GetConfig())

			body, resp, respError := client.
				DefaultAPI.PostDeployNew(context.TODO()).
				DeploymentCreateInputBody(golfsdk.DeploymentCreateInputBody{
					Urls:               []golfsdk.Url{{Domain: args[0]}},
					ExternalSourceType: &externalSourceType,
					ExternalSource:     &externalSource,
					Name:               &name,
				}).
				Execute()

			if respError != nil {
				panic(respError.Error())
			}
			fmt.Printf("status: %v\n", resp.Status)
			fmt.Printf("response body: %+v\n", body)
		},
	}

	createDeployment.Flags().StringVar(
		&github, "github", "", "Specify a Github Repo: repoOwner/repoName",
	)
	createDeployment.Flags().StringVar(
		&path, "path", "", "Specify a path, like \"/my-page\"",
	)
	createDeployment.Flags().StringVar(
		&name, "name", "", "Specify a name for the deployment.",
	)

	return &createDeployment
}

func deployContentCommand() *cobra.Command {
	var name string
	var files string

	deployContent := cobra.Command{
		Use:   "deploy-content",
		Short: "Deploys content",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.TODO()
			fileTree, err := archives.FilesFromDisk(ctx, nil, map[string]string{
				files: "",
			})
			if err != nil {
				panic(err.Error())
			}
			format := archives.CompressedArchive{
				Compression: archives.Gz{},
				Archival:    archives.Tar{},
			}
			tempFile, tempFileErr := os.CreateTemp("", "files-to-deploy")
			if tempFileErr != nil {
				panic(tempFileErr.Error())
			}
			defer os.Remove(tempFile.Name())

			archiveErr := format.Archive(ctx, tempFile, fileTree)
			if archiveErr != nil {
				panic(archiveErr.Error())
			}

			tempFile.Seek(0, 0)

			client := createClient()

			body, resp, respError := client.
				DefaultAPI.PutDeployFiles(context.TODO()).
				Name(name).
				Contents(tempFile).
				Execute()

			if respError != nil {
				panic(respError.Error())
			}
			fmt.Printf("status: %v\n", resp.Status)
			fmt.Printf("response body: %+v\n", body)
			if body == nil || !body.Success {
				panic("Did not get success status back from server")
			}
		},
	}

	deployContent.Flags().StringVar(
		&name, "name", "",
		"Specify the name of the deployment you wish to update.",
	)

	deployContent.Flags().StringVar(
		&files, "files", "",
		"Supply a path to a directory with the content you wish to deploy.",
	)

	return &deployContent
}

func main() {
	rootCmd.AddCommand(createDeploymentCommand())
	rootCmd.AddCommand(deployContentCommand())
	// TODO: the default should actually depend on the passed-in url arg(s).
	// also, the /internet--golf--admin path should be added for non-local
	// requests
	rootCmd.PersistentFlags().StringVar(
		&apiUrl, "apiUrl", "http://localhost:8888", "Specify the API URL",
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
