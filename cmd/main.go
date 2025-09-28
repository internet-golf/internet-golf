package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	internetgolf "github.com/toBeOfUse/internet-golf/pkg"
)

func main() {
	var adminApiPort string
	var localOnly bool
	var adminApiUrl string
	var dataDirectory string
	var verbose bool

	var rootCmd = &cobra.Command{
		Use:   "golf-server",
		Short: "A server to which you can deploy stuff",
		Long: "An instance of Internet Golf that you can use to deploy websites. " +
			"You probably don't need to worry about the CLI flags.",
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {

			// the core architecture of this app consists of these top-level actors:

			// 0. the db module
			storageSettings := internetgolf.StorageSettings{}
			storageSettings.Init(dataDirectory)
			db := internetgolf.StormDb{}
			if err := db.Init(storageSettings); err != nil {
				panic(err)
			}

			// 1. interface to the web server that actually deploys the deployments
			deploymentServer := internetgolf.CaddyServer{}
			deploymentServer.Settings.LocalOnly = localOnly
			deploymentServer.Settings.Verbose = verbose

			// 2. object that receives the active deployments and broadcasts
			// them to the deploymentServer when necessary
			deploymentBus := internetgolf.DeploymentBus{
				Server: &deploymentServer,
				Db:     &db,
			}
			deploymentBus.Init()

			// 3. api server that receives admin API requests and updates the active
			// deployments in response to them
			adminApi := internetgolf.AdminApi{
				Web:  deploymentBus,
				Auth: internetgolf.AuthManager{Db: &db},
				Port: adminApiPort,
			}

			// putting the pieces together:

			// create a deployment for the admin api (slightly premature, but
			// that's fine as long as the health check endpoint is used)
			adminApiName := "__internet__golf__admin__"
			deploymentBus.SetupDeployment(
				internetgolf.DeploymentMetadata{
					Name: adminApiName,
					Urls: []internetgolf.Url{
						internetgolf.Url{Path: adminApiUrl},
					},
					DontPersist: true,
				})
			deploymentBus.PutDeploymentContentByName(
				adminApiName,
				internetgolf.DeploymentContent{
					ServedThingType: internetgolf.ReverseProxy,
					ServedThing:     "localhost:" + adminApiPort,
				})

			// start the admin api
			server := adminApi.CreateServer()
			server.ListenAndServe()
		},
	}

	// commented out bc it makes it hard for the client to find the admin API to
	// make requests to it from the same machine
	// openPort, openPortErr := internetgolf.GetFreePort()
	// if openPortErr != nil {
	// 	panic(openPortErr.Error())
	// }
	rootCmd.Flags().StringVar(
		&adminApiPort, "admin-api-port", "8888", // strconv.Itoa(openPort),
		"Specify a port for the internal admin API.\n"+
			"This is only really useful for testing and to avoid port conflicts.",
	)
	rootCmd.Flags().BoolVar(
		&localOnly, "local", false,
		"Run in local-only mode, so that deployments are only available at localhost:80.",
	)
	rootCmd.Flags().StringVar(
		&adminApiUrl, "admin-api-path", "/internet--golf--admin",
		"Path prefix for the Admin API endpoints.",
	)
	rootCmd.Flags().StringVar(
		&dataDirectory, "data-dir", "$HOME/.internetgolf",
		"Location on disk where deployments will be stored. "+
			"Separate from Caddy's data directory.\n",
	)
	rootCmd.Flags().BoolVarP(
		&verbose, "verbose", "v", false,
		"Output all internal logs",
	)

	var openapiOutputPath string

	outputOpenapiCommand := &cobra.Command{
		Use:  "openapi",
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			adminApi := internetgolf.AdminApi{}
			adminApi.OutputOpenApiSpec(openapiOutputPath)
		},
	}
	outputOpenapiCommand.Flags().StringVar(
		&openapiOutputPath, "--output", "golf-openapi.yaml",
		"Path to the YAML file to output to",
	)

	rootCmd.AddCommand(outputOpenapiCommand)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
