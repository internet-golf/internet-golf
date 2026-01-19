// this file runs integration tests across the whole system: the client cli, the
// admin API, and the deployment server.
//
// for each test case, the client cli is called from the shell against a real
// server, and a function is run to verify the deployment was created correctly.
//
// these tests lean on automatically authenticating requests by making them
// through localhost. auth tests live in auth_integration_test.go.

package internetgolf_test

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"testing"

	golfsdk "github.com/internet-golf/internet-golf/client-sdk"
	"github.com/internet-golf/internet-golf/pkg/api"
	database "github.com/internet-golf/internet-golf/pkg/db"
	"github.com/internet-golf/internet-golf/pkg/public"
	"github.com/internet-golf/internet-golf/pkg/resources"
	"github.com/internet-golf/internet-golf/pkg/utils"
)

// test case stuff =======================================================

type TestCase struct {
	// name for the test case: used for logging
	name string
	// client cli commands that should be run against the real api to prepare
	// for the test
	setupCommands []string
	// client cli command that should be run for the test case
	cliCommand string
	// function that should be run after the client cli is called to check the
	// state of the server and verify that it was updated correctly
	deploymentTest func(*testing.T, *golfsdk.APIClient)
}

var testCases = []TestCase{
	{
		name:       "Create basic deployment",
		cliCommand: "create-deployment example.com",
		deploymentTest: func(t *testing.T, client *golfsdk.APIClient) {
			output, _, err := client.DefaultAPI.GetDeployment(context.TODO(), "example.com").Execute()
			if err != nil {
				t.Fatal(err)
			}
			if (output.EmptyDeployment.Type != "Empty") || output.EmptyDeployment.Url != "example.com" {
				t.Fail()
			}
		},
	},
	// TODO: URLs with paths, other settings
	{
		name:       "Upload some files",
		cliCommand: "deploy-content internet-golf-test.local --files ./fixtures/static-site",
		deploymentTest: func(t *testing.T, _ *golfsdk.APIClient) {
			if content := urlToPageContent("http://internet-golf-test.local", t); content != "stuff\n" {
				t.Fatalf("expected stuff\\n, got %v", []byte(content))
			}
			if content := urlToPageContent("http://internet-golf-test.local/nested/concept.txt", t); content != "fnord" {
				t.Fatalf("expected fnord, got %v", []byte(content))
			}
		},
	},
}

// server setup ===========================================================

func startFullServer(port string) func() {
	tempDir, tempDirError := os.MkdirTemp("", "internet-golf-test")
	if tempDirError != nil {
		panic(tempDirError)
	}
	tempDirs = append(tempDirs, tempDir)

	config := utils.NewConfig(tempDir, true, true, port)

	fileManager := resources.NewFileManager(config)

	db, err := database.NewDb(config, fileManager)
	if err != nil {
		panic(err)
	}

	deploymentServer, err := public.NewPublicWebServer(config, fileManager)
	if err != nil {
		panic(err)
	}

	deploymentBus, err := api.NewDeploymentBus(deploymentServer, db, fileManager)
	if err != nil {
		panic(err)
	}

	adminApi := api.NewAdminApi(deploymentBus, db, config)

	// TODO: this default admin API path needs to be a global constant somewhere
	adminApiUrl := database.Url{Path: "/_golf"}
	deploymentBus.SetupDeployment(
		database.DeploymentMetadata{
			Url:         adminApiUrl,
			DontPersist: true,
		})
	deploymentBus.PutDeploymentContentByUrl(
		adminApiUrl,
		database.DeploymentContent{
			ServedThingType: database.ReverseProxy,
			ServedThing:     "127.0.0.1:" + port,
		})

	server := adminApi.CreateServer()
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		panic(err)
	}
	go func() {
		// always returns error. ErrServerClosed on graceful close
		if err := server.Serve(listener); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("Serve(): %v", err)
		}
	}()

	return func() {
		server.Shutdown(context.TODO())
		deploymentBus.Stop()
	}
}

func TestClientCli(t *testing.T) {
	serverPortInt, portErr := utils.GetFreePort()
	if portErr != nil {
		panic(portErr)
	}
	serverPort := strconv.Itoa(serverPortInt)
	client := createClient("http://127.0.0.1:" + serverPort)

	stopServer := startFullServer(serverPort)
	defer stopServer()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			for _, command := range testCase.setupCommands {
				runClientCliCommand(command, serverPort, t)
			}

			runClientCliCommand(testCase.cliCommand, serverPort, t)

			testCase.deploymentTest(t, client)
		})
	}
}
