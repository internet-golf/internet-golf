package internetgolf_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	golfsdk "github.com/toBeOfUse/internet-golf/client-sdk"
	internetgolf "github.com/toBeOfUse/internet-golf/pkg"
)

func createClient(url string) *golfsdk.APIClient {
	return golfsdk.NewAPIClient(&golfsdk.Configuration{
		UserAgent: "InternetGolfClient",
		Servers: golfsdk.ServerConfigurations{
			{URL: url},
		},
	})
}

type CliApiTestCase[ReqBody interface {
	// this union will be expanded
	internetgolf.DeploymentCreateInput | internetgolf.DeployFilesInput
}] struct {
	name           string
	cliCommand     string
	apiPath        string
	apiMethod      string
	apiReqBody     ReqBody
	deploymentTest func(*testing.T)
}

type InterceptedRequest struct {
	Req  http.Request
	Body []byte
}

type MockApiServer struct {
	port     string
	server   *http.Server
	reqQueue chan InterceptedRequest
}

func (m *MockApiServer) Init() {
	// making this have a buffer of size 2 so that sends and receives don't
	// instantly block
	m.reqQueue = make(chan InterceptedRequest, 2)
	sm := http.NewServeMux()
	// "/" matches every path. for some reason
	sm.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// the body has to be read before this handler finishes because the
		// server automatically closes it; hence, reading it here and then
		// sending it to the test function in an InterceptedRequest, instead of
		// just sending `r` to the test runner
		body, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err.Error())
		}
		m.reqQueue <- InterceptedRequest{
			Req:  *r,
			Body: body,
		}

	})
	port, err := internetgolf.GetFreePort()
	if err != nil {
		panic(err)
	}
	m.port = strconv.Itoa(port)
	m.server = &http.Server{
		Addr:    "localhost:" + m.port,
		Handler: sm,
	}

	fmt.Printf("starting mock server at localhost:%s\n", m.port)
	go func() {
		// always returns error. ErrServerClosed on graceful close
		if err := m.server.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	time.Sleep(2 * time.Second)
}

func (m *MockApiServer) Stop() {
	if err := m.server.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}

func runClientCliCommand(command string, apiPort string) {
	cmd := exec.Command(
		"go",
		strings.Split(
			"run ../client-cmd --apiUrl http://localhost:"+apiPort+" "+command, " ",
		)...,
	)
	fmt.Printf("running %v\n", cmd.Args)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("%v", err.Error())
		panic(err)
	}
}

var realServerPort = "9999"
var client = createClient("http://127.0.0.1:" + realServerPort)

func startRealServer() func() {
	deploymentServer := internetgolf.CaddyServer{}
	deploymentServer.Settings.LocalOnly = true

	tempDir, tempDirError := os.MkdirTemp("", "internet-golf-test")
	if tempDirError != nil {
		panic(tempDirError)
	}
	tempDirs = append(tempDirs, tempDir)

	settings := internetgolf.StorageSettings{}
	settings.Init(tempDir)

	deploymentBus := internetgolf.DeploymentBus{
		Server:          &deploymentServer,
		StorageSettings: settings,
	}
	deploymentBus.Init()
	adminApi := internetgolf.AdminApi{
		Web:             deploymentBus,
		StorageSettings: settings,
		Port:            realServerPort,
	}

	server := adminApi.CreateServer()
	go func() {
		err := server.ListenAndServe()
		if err != nil && err.Error() != "http: Server closed" {
			panic(err)
		}
	}()

	for i := 0; i < 100; i++ {
		body, res, err := client.DefaultAPI.GetAlive(context.TODO()).Execute()
		if err == nil && res.StatusCode == 200 && body.Ok {
			break
		}
		if i == 99 {
			panic("real server failed health check")
		}
		time.Sleep(100 * time.Millisecond)
	}

	return func() {
		server.Shutdown(context.TODO())
		deploymentBus.Stop()
	}
}

// this tests deployment creation (with no content added yet)
func TestCreateDeployment(t *testing.T) {
	// create mock intermediary server
	m := MockApiServer{}
	m.Init()
	defer m.Stop()

	// create real server
	stopRealServer := startRealServer()
	defer stopRealServer()

	testCases := []CliApiTestCase[internetgolf.DeploymentCreateInput]{
		{
			name:       "Create basic deployment",
			cliCommand: "create-deployment --name example-com example.com",
			apiPath:    "/deploy/new",
			apiReqBody: internetgolf.DeploymentCreateInput{
				Body: struct {
					internetgolf.DeploymentMetadata
				}{
					DeploymentMetadata: internetgolf.DeploymentMetadata{
						Name: "example-com",
						Urls: []internetgolf.Url{{Domain: "example.com", Path: ""}},
					},
				},
			},
			deploymentTest: func(t *testing.T) {
				output, res, err := client.DefaultAPI.GetDeploymentByName(context.TODO(), "example-com").Execute()
				if err != nil {
					fmt.Printf("deploymentTest error: %+v\n", err)
					t.Fail()
					return
				}
				if res.StatusCode != 200 {
					fmt.Printf("deploymentTest error: status was %s", res.Status)
					t.Fail()
					return
				}
				fmt.Printf("OUTPUT: %+v\n", output)
				if output.Urls[0].Domain != "example.com" {
					t.Fail()
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			runClientCliCommand(testCase.cliCommand, m.port)
			intercepted := <-m.reqQueue
			req := intercepted.Req
			if req.URL.Path != testCase.apiPath {
				t.Fatalf("expected %s, got %s\n", testCase.apiPath, req.URL.Path)
			}

			fmt.Printf("%+v\n", req)

			var contents internetgolf.DeploymentMetadata
			jsonErr := json.Unmarshal(intercepted.Body, &contents)
			if jsonErr != nil {
				t.Fatal(jsonErr.Error())
			}

			if !contents.Equals(&testCase.apiReqBody.Body.DeploymentMetadata) {
				t.Fatalf("%s failed; incorrect struct: %+v\n", testCase.name, contents)
			}

			if _, err := http.Post("http://127.0.0.1:"+realServerPort+testCase.apiPath, "application/json", bytes.NewReader(intercepted.Body)); err != nil {
				t.Fatalf("%s", err.Error())
			}

			testCase.deploymentTest(t)

		})
	}
}
