package internetgolf_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	internetgolf "github.com/toBeOfUse/internet-golf/pkg"
)

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

type MockApiServer struct {
	port     string
	server   *http.Server
	reqQueue chan http.Request
}

func (m *MockApiServer) Init() {
	sm := http.NewServeMux()
	// "/" matches every path. for some reason
	sm.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("request to %s\n", r.RequestURI)
		// TODO: we are blocking here bc the client command never completes and
		// thus the test function can't proceed to receive this message
		// find a way to send an early 200 back?
		m.reqQueue <- *r

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

func TestBasicDeploymentCreate(t *testing.T) {
	testCases := []CliApiTestCase[internetgolf.DeploymentCreateInput]{
		{
			name:       "Create basic deployment",
			cliCommand: "create-deployment example.com",
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
			// TODO: test that the new deployment can be retrieved via the api i
			// guess?
			deploymentTest: func(t *testing.T) {},
		},
	}

	// create mock intermediary server
	m := MockApiServer{}
	m.Init()

	// create real server
	bus := createBus()
	defer bus.Stop()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			runClientCliCommand(testCase.cliCommand, m.port)
			req := <-m.reqQueue
			if req.URL.Path != testCase.apiPath {
				t.Fatalf("expected %s, got %s\n", testCase.apiPath, req.URL.Path)
			}
			// TODO: assert that request body matches expectations

			// TODO: forward the api req to the real server; run
			// deploymentTest()
		})
	}
}
