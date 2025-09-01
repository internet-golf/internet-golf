package internetgolf_test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/txn2/txeh"

	internetgolf "github.com/toBeOfUse/internet-golf/pkg"
)

var deploymentBus internetgolf.DeploymentBus
var busCreated bool = false

func createBus() {
	if !busCreated {

		tempDir, tempDirError := os.MkdirTemp("", "internet-golf-test")
		if tempDirError != nil {
			panic(tempDirError)
		}

		settings := internetgolf.StorageSettings{}
		settings.Init(tempDir)

		// interface to the web server that actually deploys the deployments
		deploymentServer := internetgolf.CaddyServer{}
		deploymentServer.Settings.LocalOnly = true

		// object that (persistently) stores the active deployments and broadcasts
		// them to the deploymentServer when necessary
		deploymentBus = internetgolf.DeploymentBus{
			Server:          &deploymentServer,
			StorageSettings: settings,
		}
		deploymentBus.Init()

		busCreated = true
	}
}

const (
	BasicTestHost = "internet-golf-test.local"
	CacheTestHost = "internet-golf-cache-test.local"
)

// TODO: this requires elevated permissions to run and only needs to run once...
func setupHosts() {
	// this does not actually appear to create a new hosts file but rather
	// loads the existing one
	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		panic(err)
	}
	fmt.Println(hosts)
	hosts.AddHost("127.0.0.1", BasicTestHost)
	hosts.AddHost("127.0.0.1", CacheTestHost)
	if saveErr := hosts.Save(); saveErr != nil {
		panic(saveErr)
	}
}

func TestMain(m *testing.M) {
	createBus()
	// setupHosts()
	code := m.Run()

	os.RemoveAll(deploymentBus.StorageSettings.DataDirectory)
	os.Exit(code)
}

func urlToPageContent(url string, t *testing.T) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, bodyErr := io.ReadAll(resp.Body)
	if bodyErr != nil {
		panic(bodyErr)
	}
	bodyStr := string(body)
	return bodyStr, nil
}

func assertUrlEmpty(url string, t *testing.T) {
	// TODO: should this return an error? the caddy Routes type says "By
	// default, all unrouted requests receive a 200 OK response to indicate the
	// server is working." maybe i should add a catch-all non-persisted route
	// with default content? but with, like, a 404 status code?
	blankResp, _ := urlToPageContent(url, t)
	if len(blankResp) > 0 {
		t.Fatal("deployment already existed at beginning of test")
	}
}

func getFixturePath(fixture string) string {
	cwd, wdErr := os.Getwd()
	if wdErr != nil {
		panic(wdErr)
	}
	return path.Join(
		// for some reason the cwd already includes /test/
		strings.ReplaceAll(cwd, "\\", "/"), "fixtures", fixture)
}

func TestBasicStaticDeployment(t *testing.T) {

	// create a deployment that serves the static-site fixture at
	// http://internet-golf-test.local

	url := "http://" + BasicTestHost
	assertUrlEmpty(url, t)

	deploymentBus.SetupDeployment(internetgolf.DeploymentMetadata{
		Urls: []internetgolf.Url{internetgolf.Url{Domain: BasicTestHost}},
		Name: "test-1",
	})

	deploymentBus.PutDeploymentContentByName("test-1", internetgolf.DeploymentContent{
		ServedThingType: internetgolf.StaticFiles,
		ServedThing:     getFixturePath("static-site"),
	})

	bodyStr, _ := urlToPageContent(url, t)
	if bodyStr != "stuff\n" {
		t.Fatalf("expected stuff\\n, got %v", []byte(bodyStr))
	}
}

func TestStaticDeploymentWithPath(t *testing.T) {

	// create a deployment that serves the static-site-2 fixture at
	// http://internet-golf-test.local/stuff/

	url := "http://" + BasicTestHost + "/stuff/"
	assertUrlEmpty(url, t)

	deploymentBus.SetupDeployment(internetgolf.DeploymentMetadata{
		Urls: []internetgolf.Url{internetgolf.Url{Domain: BasicTestHost, Path: "/stuff/"}},
		Name: "test-2",
	})

	deploymentBus.PutDeploymentContentByName("test-2", internetgolf.DeploymentContent{
		ServedThingType: internetgolf.StaticFiles,
		ServedThing:     getFixturePath("static-site-2"),
	})

	configStr, _ := urlToPageContent("http://localhost:2019/config", t)
	fmt.Println(configStr)

	bodyStr, _ := urlToPageContent(url, t)
	if bodyStr != "stuff 2\n" {
		t.Fatalf("expected stuff 2\\n, got %v", []byte(bodyStr))
	}
}

func TestStaticDeploymentWithPreservedPath(t *testing.T) {

	// create a deployment that serves the static-site-2 fixture at
	// http://internet-golf-test.local/stuff/

	url := "http://" + BasicTestHost + "/other-stuff/"
	assertUrlEmpty(url, t)

	deploymentBus.SetupDeployment(internetgolf.DeploymentMetadata{
		Urls:                 []internetgolf.Url{internetgolf.Url{Domain: BasicTestHost, Path: "/other-stuff/"}},
		Name:                 "test-3",
		PreserveExternalPath: true,
	})

	deploymentBus.PutDeploymentContentByName("test-3", internetgolf.DeploymentContent{
		ServedThingType: internetgolf.StaticFiles,
		ServedThing:     getFixturePath("static-site-3"),
	})

	configStr, _ := urlToPageContent("http://localhost:2019/config", t)
	fmt.Println(configStr)

	bodyStr, _ := urlToPageContent(url, t)
	if bodyStr != "stuff 3\n" {
		t.Fatalf("expected stuff 3\\n, got %v", []byte(bodyStr))
	}
}
