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

func urlToPageContent(url string, t *testing.T) string {
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	body, bodyErr := io.ReadAll(resp.Body)
	if bodyErr != nil {
		panic(bodyErr)
	}
	bodyStr := string(body)
	return bodyStr
}

func TestBasicStaticDeployment(t *testing.T) {
	cwd, wdErr := os.Getwd()
	if wdErr != nil {
		panic(wdErr)
	}

	// create a deployment that serves the static-site fixture at
	// http://internet-golf-test.local

	deploymentBus.SetupDeployment(internetgolf.DeploymentMetadata{
		Url: BasicTestHost,
		// Settings: internetgolf.DeploymentSettings{},
	})

	deploymentBus.PutDeploymentContentByUrl(BasicTestHost, internetgolf.DeploymentContent{
		ServedThingType: internetgolf.StaticFiles,
		ServedThing: path.Join(
			// for some reason the cwd already includes /test/
			strings.ReplaceAll(cwd, "\\", "/"), "fixtures", "static-site"),
	})

	// configStr := urlToString("http://localhost:2019/config", t)
	// fmt.Println(configStr)

	bodyStr := urlToPageContent("http://"+BasicTestHost, t)
	if bodyStr != "stuff\n" {
		t.Fatalf("expected stuff\\n, got %v", []byte(bodyStr))
	}
}
