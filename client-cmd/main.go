package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	golfsdk "github.com/internet-golf/internet-golf/client-sdk"
	"github.com/mholt/archives"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "golf",
	Short: "Deploy stuff to a server",
}

var apiUrl string
var auth string

func exit1(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}

// TODO: standard response handling function that does roughly this instead of
// ever panicking

// if body == nil || respError != nil {
// // failure:
// if respError != nil {
// 	fmt.Println(respError.Error())
// }
// responseBody, responseBodyErr := io.ReadAll(resp.Body)
// if responseBodyErr != nil || len(responseBody) == 0 {
// 	fmt.Println("ERROR: Could not get response body")
// } else {
// 	fmt.Println(string(responseBody))
// }
// } else {
// // success:
// fmt.Println(body.Message)
// }

func createClient(hostnameFromTargetDeployment string) *golfsdk.APIClient {
	// determine the base URL of the API server:
	// 1. If apiUrl is set by the command line option, use apiUrl
	// 2. Alternatively, if no auth method is set, there's nothing that could
	// work other than localhost, so use localhost with the default API port (8888)
	// 3. otherwise, use the host that is passed into this function as targetHost,
	// which is the hostname intended to be deployed to.
	resolvedApiUrl := apiUrl
	if len(resolvedApiUrl) == 0 {
		if len(auth) == 0 {
			// if no auth setting is specified, assume localhost
			resolvedApiUrl = "http://localhost:8888"
		} else if len(hostnameFromTargetDeployment) > 0 {
			protocol := "https"
			ips, err := net.LookupIP(hostnameFromTargetDeployment)
			if err == nil && ips[0].String() == "127.0.0.1" {
				fmt.Fprintf(os.Stderr, "WARNING: connecting to local host %s Without HTTPS", hostnameFromTargetDeployment)
				protocol = "http"
			}
			resolvedApiUrl = protocol + "://" + hostnameFromTargetDeployment + "/_golf"
		} else {
			exit1("could not resolve API URL")
		}
	}

	authHeader := ""
	if auth == "github-oidc" {
		reqUrl, found := os.LookupEnv("ACTIONS_ID_TOKEN_REQUEST_URL")
		if !found {
			exit1("environment variable ACTIONS_ID_TOKEN_REQUEST_URL not found")
		}
		reqToken, found := os.LookupEnv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")
		if !found {
			exit1("environment variable ACTIONS_ID_TOKEN_REQUEST_TOKEN not found")
		}
		githubOidcReq, err := http.NewRequest("GET", reqUrl+"&audience=internet-golf", nil)
		if err != nil {
			exit1(err.Error())
		}
		githubOidcReq.Header["Authorization"] = []string{"Bearer " + reqToken}
		resp, err := http.DefaultClient.Do(githubOidcReq)
		if err != nil {
			exit1(err.Error())
		}
		oidcTokenJson, err := io.ReadAll(resp.Body)
		var oidcTokenData struct {
			Value string `json:"value"`
		}
		json.Unmarshal(oidcTokenJson, &oidcTokenData)
		authHeader = "GithubOIDC " + strings.Trim(string(oidcTokenData.Value), " \n\r")
	} else if len(auth) > 0 {
		authHeader = "Bearer " + auth
	}

	client := golfsdk.NewAPIClient(&golfsdk.Configuration{
		UserAgent: "InternetGolfClient",
		DefaultHeader: map[string]string{
			"Authorization": authHeader,
		},
		Servers: golfsdk.ServerConfigurations{
			{URL: resolvedApiUrl},
		},
	})

	// perform health check against the API URL that was determined above. (the
	// auth header doesn't actually matter for this part). try 20 times (with a
	// half-second pause between tries) in case the server is just starting up
	retries := 20
	for i := range retries {
		body, _, err := client.DefaultAPI.HealthCheck(context.Background()).Execute()

		if err == nil && body.Ok {
			break
		}
		if i == retries-1 {
			exit1("Health check of " + resolvedApiUrl + " failed :(")
		}
		time.Sleep(500 * time.Millisecond)
	}
	return client
}

func createDeploymentCommand() *cobra.Command {

	var github string
	// TODO: preserve external path option

	createDeployment := cobra.Command{
		Use:     "create-deployment domain",
		Example: "create-deployment example.com --github repoOwner/repoName",
		Short:   "Creates a deployment",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var externalSourceType string
			var externalSource string

			if len(github) > 0 {
				// TODO: would be nice to use the ExternalSourceType enum-ish
				// thing somehow, instead of this string literal
				externalSourceType = "GithubRepo"
				externalSource = github
			}

			client := createClient(args[0])

			body, resp, respError := client.
				DefaultAPI.CreateDeployment(context.TODO()).
				DeploymentCreateInputBody(golfsdk.DeploymentCreateInputBody{
					Url:                args[0],
					ExternalSourceType: &externalSourceType,
					ExternalSource:     &externalSource,
					Tags:               []string{},
				}).
				Execute()

			if body == nil || respError != nil {
				responseBody, responseBodyErr := io.ReadAll(resp.Body)
				if responseBodyErr != nil || len(responseBody) == 0 {
					if respError != nil {
						fmt.Println(respError.Error())
					} else {
						fmt.Println("ERROR: Could not get response body")
					}
					return
				}
				fmt.Println(string(responseBody))
				return
			}
			fmt.Println(body.Message)
		},
	}

	createDeployment.Flags().StringVar(
		&github, "github", "", "Specify a Github Repo: repoOwner/repoName",
	)

	return &createDeployment
}

func deployAdminDash() *cobra.Command {
	createDeployment := cobra.Command{
		Use:     "deploy-admin-dash domain",
		Example: "deploy-admin-dash dash.example.com",
		Short:   "Deploys the Internet Golf Admin Dashboard to a specific URL",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := createClient(args[0])

			body, resp, respError := client.
				DefaultAPI.DeployAdminDash(context.TODO()).
				DeployAdminDashBody(golfsdk.DeployAdminDashBody{Url: args[0]}).
				Execute()

			if body == nil || respError != nil {
				responseBody, responseBodyErr := io.ReadAll(resp.Body)
				if responseBodyErr != nil || len(responseBody) == 0 {
					if respError != nil {
						fmt.Println(respError.Error())
					} else {
						fmt.Println("ERROR: Could not get response body")
					}
					return
				}
				fmt.Println(string(responseBody))
				return
			}
			fmt.Println(body.Message)
		},
	}

	return &createDeployment
}

func deployAliasCommand() *cobra.Command {
	var redirect bool

	deployAlias := cobra.Command{
		Use:     "create-alias from to",
		Example: "create-alias from.example.com to.example.com",
		Short:   "Creates an alias from one deployment to another. Optionally, redirect the visitor. If the \"to\" URL is not a known deployment, it will always redirect.",
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := createClient(args[0])

			body, resp, respError := client.
				DefaultAPI.CreateDeployment(context.TODO()).
				DeploymentCreateInputBody(golfsdk.DeploymentCreateInputBody{
					Url: args[0],
				}).Execute()

			if body == nil || respError != nil {
				responseBody, responseBodyErr := io.ReadAll(resp.Body)
				if responseBodyErr != nil || len(responseBody) == 0 {
					if respError != nil {
						fmt.Println(respError.Error())
					} else {
						fmt.Println("ERROR: Could not get response body")
					}
					return
				}
				fmt.Println(string(responseBody))
				return
			}

			body, resp, respError = client.
				DefaultAPI.CreateAlias(context.TODO()).
				DeployAliasBody(golfsdk.DeployAliasBody{
					Url: args[0], AliasedTo: &args[1], Redirect: &redirect,
				}).Execute()

			if body == nil || respError != nil {
				responseBody, responseBodyErr := io.ReadAll(resp.Body)
				if responseBodyErr != nil || len(responseBody) == 0 {
					if respError != nil {
						fmt.Println(respError.Error())
					} else {
						fmt.Println("ERROR: Could not get response body")
					}
					return
				}
				fmt.Println(string(responseBody))
				return
			}
			fmt.Println(body.Message)
		},
	}

	deployAlias.Flags().BoolVar(
		&redirect, "redirect", false,
		"Redirect the visitor to the aliased-to URL.",
	)

	return &deployAlias
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

			body, resp, respError := client.
				DefaultAPI.DeployFiles(context.TODO()).
				Url(args[0]).
				Contents(tempFile).
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

func createBearerTokenCommand() *cobra.Command {
	createToken := cobra.Command{
		Use:   "create-token",
		Short: "Create a bearer token that can be used to authenticate API requests",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			client := createClient("")
			body, _, err := client.DefaultAPI.
				PostTokenGenerate(context.TODO()).
				CreateBearerTokenInputBody(
					// TODO: granular permissions
					golfsdk.CreateBearerTokenInputBody{FullPermissions: true},
				).Execute()
			if err != nil || len(body.Token) == 0 {
				panic(err)
			}
			fmt.Println("Generated token:")
			fmt.Println(body.Token)
		},
	}

	return &createToken
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
		createDeploymentCommand(), deployContentCommand(),
		registerExternalUserCommand(), createBearerTokenCommand(),
		deployAdminDash(), deployAliasCommand(),
	}
	for _, cmd := range golfCmds {
		cmd.GroupID = "IG"
		rootCmd.AddCommand(cmd)
	}

	rootCmd.PersistentFlags().StringVar(
		&apiUrl, "api-url", "", "Specify the API URL. Will be smartly guessed if not present.",
	)
	rootCmd.PersistentFlags().StringVar(
		&auth, "auth", "", "Specify a bearer token or give the value \"github-oidc\".",
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
