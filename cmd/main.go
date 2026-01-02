package main

import (
	"fmt"
	"os"

	"github.com/internet-golf/internet-golf/pkg/admin_api"
	database "github.com/internet-golf/internet-golf/pkg/db"
	"github.com/internet-golf/internet-golf/pkg/public_web_server"
	"github.com/internet-golf/internet-golf/pkg/settings"
	"github.com/spf13/cobra"
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
		// TODO: can this function be pulled out to use in tests? there's a
		// re-implementation of it in utils_test.go
		Run: func(cmd *cobra.Command, args []string) {

			config := settings.NewConfig(dataDirectory, localOnly, verbose, adminApiPort)

			db, err := database.NewDb(config)
			if err != nil {
				panic(err)
			}

			deploymentServer, err := public_web_server.NewPublicWebServer(config)
			if err != nil {
				panic(err)
			}

			deploymentBus, err := admin_api.NewDeploymentBus(deploymentServer, db)
			if err != nil {
				panic(err)
			}

			adminApi := admin_api.NewAdminApi(deploymentBus, db, config)

			// create a deployment for the admin api (slightly premature, but
			// that's fine as long as the health check endpoint is used)
			adminApiUrl := database.Url{Path: adminApiUrl}

			deploymentBus.SetupDeployment(
				database.DeploymentMetadata{
					Url:         adminApiUrl,
					DontPersist: true,
				})
			deploymentBus.PutDeploymentContentByUrl(
				adminApiUrl,
				database.DeploymentContent{
					ServedThingType: database.ReverseProxy,
					ServedThing:     "127.0.0.1:" + adminApiPort,
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
		&adminApiUrl, "admin-api-path", "/_golf",
		"Path prefix for the Admin API endpoints.",
	)
	rootCmd.Flags().StringVar(
		&dataDirectory, "data-dir", "$HOME/.internetgolf",
		"Location on disk where deployment content and configuration will be stored.",
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
			adminApi := admin_api.AdminApi{}
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
