// this file runs integration tests across the whole system: the client cli, the
// admin API, and the deployment server.
//
// for each test case, the client cli is called from the shell; it makes a
// request to a mock http server that intercepts the request; and the request is
// inspected to make sure that the client cli is making the correct request.
// then, a copy of the request is forwarded to a real version of the admin API,
// and a function is run to make sure that the deployment was created
// accordingly.
//
// these tests lean on automatically authenticating requests by making them
// through localhost. auth tests live in auth_integration_test.go.
//
// there's a lot of plumbing going on here. all you should really have to worry
// about when creating new tests are the tests cases at the beginning.

package internetgolf_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	golfsdk "github.com/internet-golf/internet-golf/client-sdk"
	"github.com/internet-golf/internet-golf/pkg/api"
	database "github.com/internet-golf/internet-golf/pkg/db"
	"github.com/internet-golf/internet-golf/pkg/public"
	"github.com/internet-golf/internet-golf/pkg/resources"
	"github.com/internet-golf/internet-golf/pkg/utils"
)

// test case stuff =======================================================

type CliApiTestCase struct {
	// name for the test case: used for logging
	name string
	// client cli commands that should be run against the real api to prepare
	// for the test
	setupCommands []string
	// client cli command that should be run for the test case; it will be given
	// the url of the mock server to send requests to so that they can be
	// intercepted and checked to make sure they match the fields below
	cliCommand string
	// path of the api request that the client cli should make as a result of
	// the command
	apiPath string
	// method of the api request that the client cli should make as a result of
	// the command
	apiMethod string
	// function that should be run after the client cli is called to check the
	// state of the server and verify that it was updated correctly
	deploymentTest func(*testing.T, *golfsdk.APIClient)
}

type NewDeploymentTestCase struct {
	CliApiTestCase
	// expected body for the API request that the client CLI will make to the
	// server
	apiBody api.DeploymentCreateBody
}

var deploymentCreateTestCases = []NewDeploymentTestCase{
	{
		CliApiTestCase: CliApiTestCase{
			name:       "Create basic deployment",
			cliCommand: "create-deployment example.com",
			apiPath:    "/deploy/new",
			apiMethod:  "PUT",
			deploymentTest: func(t *testing.T, client *golfsdk.APIClient) {
				output, _, _ := client.DefaultAPI.GetDeploymentByUrl(context.TODO(), "example.com").Execute()
				if output.Url.Domain != "example.com" {
					t.Fail()
				}
			},
		},
		apiBody: api.DeploymentCreateBody{
			Url: "example.com",
		},
	},
	// TODO: URLs with paths, other settings
}

type UserAddTestCase struct {
	CliApiTestCase
	apiBody api.AddExternalUserBody
}

var addUserTestCases = []UserAddTestCase{
	{
		CliApiTestCase: CliApiTestCase{
			name:       "Register external user",
			cliCommand: "register-user --handle internet-golf",
			apiPath:    "/user/register",
			apiMethod:  "PUT",
		},
		apiBody: api.AddExternalUserBody{
			ExternalUserHandle: "internet-golf",
			ExternalUserSource: "Github",
		},
	},
}

type DeployFilesTestCase struct {
	CliApiTestCase
	formData map[string][]string
}

var deployFilesTestCases = []DeployFilesTestCase{
	{
		CliApiTestCase: CliApiTestCase{
			name:          "Upload some files",
			setupCommands: []string{"create-deployment internet-golf-test.local"},
			cliCommand:    "deploy-content internet-golf-test.local --files ./fixtures/static-site",
			apiPath:       "/deploy/files",
			apiMethod:     "PUT",
			deploymentTest: func(t *testing.T, _ *golfsdk.APIClient) {
				if content := urlToPageContent("http://internet-golf-test.local", t); content != "stuff\n" {
					t.Fatalf("expected stuff\\n, got %v", []byte(content))
				}
				if content := urlToPageContent("http://internet-golf-test.local/nested/concept.txt", t); content != "fnord" {
					t.Fatalf("expected fnord, got %v", []byte(content))
				}
			},
		},
		formData: map[string][]string{
			"url": []string{"internet-golf-test.local"},
		},
	},
}

// meaningless plumbing ==================================================

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
		if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			r.Body = io.NopCloser(bytes.NewReader(body))
			r.ParseMultipartForm(1000000)
		}
		m.reqQueue <- InterceptedRequest{
			Req:  *r,
			Body: body,
		}

		w.Header().Add("Content-Type", "application/json")
		if r.URL.Path == "/alive" {
			w.Write([]byte("{\"$schema\": \"whatever\", \"ok\": true}"))
		} else {
			w.Write([]byte("{\"$schema\": \"whatever\", \"success\": true, \"message\": \"Request to mock API received\"}"))
		}
	})
	port, err := utils.GetFreePort()
	if err != nil {
		panic(err)
	}
	m.port = strconv.Itoa(port)
	m.server = &http.Server{
		Addr:    "localhost:" + m.port,
		Handler: sm,
	}

	fmt.Printf("starting mock server at localhost:%s\n", m.port)
	listener, err := net.Listen("tcp", "localhost:"+m.port)
	if err != nil {
		panic(err)
	}
	go func() {
		// always returns error. ErrServerClosed on graceful close
		if err := m.server.Serve(listener); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()
}

func (m *MockApiServer) Stop() {
	if err := m.server.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}

func startFullServer(port string) func() {
	tempDir, tempDirError := os.MkdirTemp("", "internet-golf-test")
	if tempDirError != nil {
		panic(tempDirError)
	}
	tempDirs = append(tempDirs, tempDir)

	config := utils.NewConfig(tempDir, true, true, port)

	db, err := database.NewDb(config)
	if err != nil {
		panic(err)
	}

	deploymentServer, err := public.NewPublicWebServer(config)
	if err != nil {
		panic(err)
	}

	deploymentBus, err := api.NewDeploymentBus(deploymentServer, db, resources.NewFileManager(config))
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
	// create mock intermediary server
	m := MockApiServer{}
	m.Init()
	defer m.Stop()

	realServerPortInt, portErr := utils.GetFreePort()
	if portErr != nil {
		panic(portErr)
	}
	realServerPort := strconv.Itoa(realServerPortInt)
	realClient := createClient("http://127.0.0.1:" + realServerPort)

	// create real server
	stopRealServer := startFullServer(realServerPort)
	defer stopRealServer()

	// TODO: deduplicate these two for loops somehow. some parts are the same
	// and some parts are different

	for _, testCase := range deploymentCreateTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			for _, command := range testCase.setupCommands {
				runClientCliCommand(command, realServerPort, t)
			}

			// run client cli command
			runClientCliCommand(testCase.cliCommand, m.port, t)

			// get the requests that were sent to the mock server as a result of
			// the client cli command

			// 1. standard health check request
			health := <-m.reqQueue
			if !strings.HasSuffix(health.Req.URL.Path, "/alive") {
				t.Fatalf("health check not performed in test case %s", testCase.name)
			}

			// 2. actual contentful request
			intercepted := <-m.reqQueue
			req := intercepted.Req

			// do assertions about the intercepted request

			if req.URL.Path != testCase.apiPath {
				t.Fatalf("expected %s, got %s\n", testCase.apiPath, req.URL.Path)
			}

			if req.Method != testCase.apiMethod {
				t.Fatalf("expected %s, got %s\n", testCase.apiMethod, req.Method)
			}

			var contents api.DeploymentCreateBody
			jsonErr := json.Unmarshal(intercepted.Body, &contents)
			if jsonErr != nil {
				t.Fatal(jsonErr.Error())
			}

			if !reflect.DeepEqual(contents, testCase.apiBody) {
				t.Fatalf("%s failed; expected %+v, got %+v", testCase.name, testCase.apiBody, contents)
			}

			// forward the intercepted request to the real server
			realUrl, err := url.Parse("http://127.0.0.1:" + realServerPort + testCase.apiPath)
			if err != nil {
				panic(err)
			}
			if _, err := http.DefaultClient.Do(
				&http.Request{
					Method: testCase.apiMethod,
					URL:    realUrl,
					Body:   io.NopCloser(bytes.NewReader(intercepted.Body)),
					Header: req.Header,
				},
			); err != nil {
				t.Fatalf("%s", err.Error())
			}

			// run the given deployment test function that should verify that
			// the real server's state has been updated correctly
			testCase.deploymentTest(t, realClient)
		})
	}

	for _, testCase := range deployFilesTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			for _, command := range testCase.setupCommands {
				runClientCliCommand(command, realServerPort, t)
			}

			// run client cli command
			runClientCliCommand(testCase.cliCommand, m.port, t)

			// get the requests that were sent to the mock server as a result of
			// the client cli command

			// 1. standard health check request
			health := <-m.reqQueue
			if !strings.HasSuffix(health.Req.URL.Path, "/alive") {
				t.Fatalf("health check not performed in test case %s", testCase.name)
			}

			// 2. actual contentful request
			intercepted := <-m.reqQueue
			req := intercepted.Req

			// do assertions about the intercepted request

			if req.URL.Path != testCase.apiPath {
				t.Fatalf("expected %s, got %s\n", testCase.apiPath, req.URL.Path)
			}

			if req.Method != testCase.apiMethod {
				t.Fatalf("expected %s, got %s\n", testCase.apiMethod, req.Method)
			}

			// note that we just compare the form values and not the files, bc
			// comparing the files is hard and should be covered by the
			// deploymentTest
			if !reflect.DeepEqual(req.MultipartForm.Value, testCase.formData) {
				t.Fatalf(
					"invalid form values: expected %+v, got %+v\n",
					testCase.formData,
					req.MultipartForm.Value,
				)
			}

			// forward the intercepted request to the real server
			realUrl, err := url.Parse("http://127.0.0.1:" + realServerPort + testCase.apiPath)
			if err != nil {
				panic(err)
			}

			if res, err := http.DefaultClient.Do(
				&http.Request{
					Method:        req.Method,
					URL:           realUrl,
					Body:          io.NopCloser(bytes.NewReader(intercepted.Body)),
					Header:        req.Header,
					ContentLength: req.ContentLength,
				},
			); err != nil {
				t.Fatal(err.Error())
			} else if res.StatusCode != 200 {
				body, _ := io.ReadAll(res.Body)
				t.Fatalf("Received error when forwarding request: %s: %s", res.Status, string(body))
			}

			// run the given deployment test function that should verify that
			// the real server's state has been updated correctly
			testCase.deploymentTest(t, realClient)
		})
	}

	for _, testCase := range addUserTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			for _, command := range testCase.setupCommands {
				runClientCliCommand(command, realServerPort, t)
			}

			// run client cli command
			runClientCliCommand(testCase.cliCommand, m.port, t)

			// get the requests that were sent to the mock server as a result of
			// the client cli command

			// 1. standard health check request
			health := <-m.reqQueue
			if !strings.HasSuffix(health.Req.URL.Path, "/alive") {
				t.Fatalf("health check not performed in test case %s", testCase.name)
			}

			// 2. actual contentful request
			intercepted := <-m.reqQueue
			req := intercepted.Req

			// do assertions about the intercepted request

			if req.URL.Path != testCase.apiPath {
				t.Fatalf("expected %s, got %s\n", testCase.apiPath, req.URL.Path)
			}

			if req.Method != testCase.apiMethod {
				t.Fatalf("expected %s, got %s\n", testCase.apiMethod, req.Method)
			}

			fmt.Printf("intercepted body %v\n", string(intercepted.Body))

			var contents api.AddExternalUserBody
			jsonErr := json.Unmarshal(intercepted.Body, &contents)
			if jsonErr != nil {
				t.Fatal(jsonErr.Error())
			}
			if !reflect.DeepEqual(contents, testCase.apiBody) {
				t.Fatalf("invalid api body: expect %+v, got %+v", testCase.apiBody, contents)
			}

			// forward the intercepted request to the real server
			realUrl, err := url.Parse("http://127.0.0.1:" + realServerPort + testCase.apiPath)
			if err != nil {
				panic(err)
			}

			if res, err := http.DefaultClient.Do(
				&http.Request{
					Method:        req.Method,
					URL:           realUrl,
					Body:          io.NopCloser(bytes.NewReader(intercepted.Body)),
					Header:        req.Header,
					ContentLength: req.ContentLength,
				},
			); err != nil {
				t.Fatal(err.Error())
			} else if res.StatusCode != 200 {
				body, _ := io.ReadAll(res.Body)
				t.Fatalf("Received error when forwarding request: %s: %s", res.Status, string(body))
			}

			// run the given deployment test function that should verify that
			// the real server's state has been updated correctly
			if testCase.deploymentTest != nil {
				testCase.deploymentTest(t, realClient)
			}
		})
	}
}
