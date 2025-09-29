package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mholt/archives"
	"github.com/spf13/cobra"
	golfsdk "github.com/toBeOfUse/internet-golf/client-sdk"
)

var rootCmd = &cobra.Command{
	Use:   "golf",
	Short: "Deploy stuff to a server",
}

var apiUrl string
var auth string

func createClient(hostToTry string) *golfsdk.APIClient {
	resolvedApiUrl := apiUrl
	if len(resolvedApiUrl) == 0 {
		if len(auth) == 0 {
			// if no auth setting is specified, assume localhost
			resolvedApiUrl = "http://localhost:8888"
		} else if len(hostToTry) > 0 {
			resolvedApiUrl = "https://" + hostToTry + "/internet--golf--admin"
		} else {
			panic("could not resolve API URL")
		}
	}
	// TODO: run health check, wait for it to pass
	return golfsdk.NewAPIClient(&golfsdk.Configuration{
		UserAgent: "InternetGolfClient",
		Servers: golfsdk.ServerConfigurations{
			{URL: resolvedApiUrl},
		},
	})
}

func createDeploymentCommand() *cobra.Command {

	var github string

	createDeployment := cobra.Command{
		Use:     "create-deployment domain",
		Example: "create-deployment example.com --github repoOwner/repoName",
		Short:   "Creates a deployment",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var externalSourceType string
			var externalSource string

			if len(github) > 0 {
				// TODO: would be nice to use ExternalSourceType somehow
				externalSourceType = "GithubRepo"
				externalSource = github
			}

			client := createClient(args[0])
			fmt.Printf("created client with config %+v\n", client.GetConfig())

			var urls []golfsdk.Url
			for _, arg := range args {
				firstSlash := strings.Index(arg, "/")
				if firstSlash == -1 {
					urls = append(urls, golfsdk.Url{Domain: arg})
				} else {
					path := arg[firstSlash:]
					urls = append(urls, golfsdk.Url{Domain: arg[0:firstSlash], Path: &path})
				}
			}

			body, _, respError := client.
				DefaultAPI.PostDeployNew(context.TODO()).
				DeploymentCreateInputBody(golfsdk.DeploymentCreateInputBody{
					Url:                args[0],
					ExternalSourceType: &externalSourceType,
					ExternalSource:     &externalSource,
				}).
				Execute()

			if respError != nil {
				panic(respError.Error())
			}
			fmt.Println(body.Message)
		},
	}

	createDeployment.Flags().StringVar(
		&github, "github", "", "Specify a Github Repo: repoOwner/repoName",
	)

	return &createDeployment
}

func deployContentCommand() *cobra.Command {
	var files string

	deployContent := cobra.Command{
		Use:     "deploy-content [deployment-name]",
		Example: "deploy-content thing.net --files ./dist",
		Short:   "Deploys content",
		Args:    cobra.ExactArgs(1),
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

			client := createClient(args[0])

			body, _, respError := client.
				DefaultAPI.PutDeployFiles(context.TODO()).
				Url(args[0]).
				Contents(tempFile).
				Execute()

			if respError != nil {
				panic(respError.Error())
			}
			fmt.Println(body.Message)
			if body == nil || !body.Success {
				panic("Did not get success status back from server")
			}
		},
	}

	deployContent.Flags().StringVar(
		&files, "files", "",
		"Supply a path to a directory with the content you wish to deploy.",
	)

	return &deployContent
}

func registerExternalUserCommand() *cobra.Command {
	var source string
	var handle string
	var id string

	registerUser := cobra.Command{
		Use:   "register-user",
		Short: "Registers an external user from (currently, only) Github",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			client := createClient("")
			body, _, respError := client.DefaultAPI.
				PutUserRegister(context.TODO()).
				AddExternalUserInputBody(golfsdk.AddExternalUserInputBody{
					ExternalUserHandle: &handle,
					ExternalUserId:     &id,
					ExternalUserSource: source,
				}).Execute()
			if respError != nil {
				panic(respError.Error())
			}
			fmt.Println(body.Message)
			if body == nil || !body.Success {
				panic("Did not get success status back from server")
			}
		},
	}

	registerUser.Flags().StringVar(&source, "source", "Github", "The place where the external user lives")
	registerUser.Flags().StringVar(&handle, "handle", "", "User's username (not needed if --id is specified)")
	registerUser.Flags().StringVar(&id, "id", "", "User's ID (not needed if --handle is specified)")

	return &registerUser
}

func main() {
	// group the real commands away from the help commands - i think it looks
	// better that way
	golfGroup := cobra.Group{
		Title: "Main Commands:",
		ID:    "IG",
	}
	rootCmd.AddGroup(&golfGroup)

	golfCmds := [](*cobra.Command){
		createDeploymentCommand(), deployContentCommand(), registerExternalUserCommand(),
	}
	for _, cmd := range golfCmds {
		cmd.GroupID = "IG"
		rootCmd.AddCommand(cmd)
	}

	// TODO: the default should actually depend on the passed-in url arg(s).
	// also, the /internet--golf--admin path should be added for non-local
	// requests
	rootCmd.PersistentFlags().StringVar(
		&apiUrl, "apiUrl", "", "Specify the API URL. Will be smartly guessed if not present.",
	)
	rootCmd.PersistentFlags().StringVar(
		&auth, "auth", "", "Specify an auth source, like \"github-oidc\"",
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
