// tests for the DeploymentBus type. this testing is kind of white-box (since it
// uses the Deployment type directly) and mostly just exists as a sanity check
// and to aid TDD.

package internetgolf_test

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/toBeOfUse/internet-golf/pkg/db"
	"github.com/toBeOfUse/internet-golf/pkg/web"
)

var tempDirs []string

func createBus() web.DeploymentBus {

	tempDir, tempDirError := os.MkdirTemp("", "internet-golf-test")
	if tempDirError != nil {
		panic(tempDirError)
	}
	tempDirs = append(tempDirs, tempDir)

	settings := db.StorageSettings{}
	settings.Init(tempDir)
	db := db.StormDb{}
	db.Init(settings)

	// interface to the web server that actually deploys the deployments
	deploymentServer := web.CaddyServer{StorageSettings: settings}
	deploymentServer.Settings.LocalOnly = true

	// object that receives the active deployments and broadcasts
	// them to the deploymentServer when necessary
	deploymentBus := web.DeploymentBus{
		Server: &deploymentServer,
		Db:     &db,
	}
	deploymentBus.Init()

	return deploymentBus
}

// this is called by TestMain which lives in utils_test.go
func busCleanup() {
	for _, tempDir := range tempDirs {
		os.RemoveAll(tempDir)
	}
}

func assertUrlEmpty(url string, t *testing.T) {
	// TODO: should this return an error? the caddy Routes type says "By
	// default, all unrouted requests receive a 200 OK response to indicate the
	// server is working." maybe i should add a catch-all non-persisted route
	// with default content? but with, like, a 404 status code?
	blankResp := urlToPageContent(url, t)
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

	deploymentBus := createBus()
	defer deploymentBus.Stop()

	// create a deployment that serves the static-site fixture at
	// http://internet-golf-test.local

	url := "http://" + BasicTestHost
	assertUrlEmpty(url, t)

	deploymentBus.SetupDeployment(db.DeploymentMetadata{
		Url: db.Url{Domain: BasicTestHost},
	})

	deploymentBus.PutDeploymentContentByUrl(
		db.Url{Domain: BasicTestHost},
		db.DeploymentContent{
			ServedThingType: db.StaticFiles,
			ServedThing:     getFixturePath("static-site"),
		})

	bodyStr := urlToPageContent(url, t)
	if bodyStr != "stuff\n" {
		t.Fatalf("expected stuff\\n, got %v", []byte(bodyStr))
	}
}

func TestBasicStaticDeploymentPersistence(t *testing.T) {
	deploymentBus := createBus()

	url := "http://" + BasicTestHost

	deploymentBus.SetupDeployment(db.DeploymentMetadata{
		Url: db.Url{Domain: BasicTestHost},
	})

	deploymentBus.PutDeploymentContentByUrl(
		db.Url{Domain: BasicTestHost},
		db.DeploymentContent{
			ServedThingType: db.StaticFiles,
			ServedThing:     getFixturePath("static-site"),
		})

	bodyStr := urlToPageContent(url, t)
	if bodyStr != "stuff\n" {
		t.Fatalf("expected stuff\\n, got %v", []byte(bodyStr))
	}

	deploymentBus.Stop()

	deploymentBus.Init()

	bodyStr = urlToPageContent(url, t)
	if bodyStr != "stuff\n" {
		t.Fatalf("expected stuff\\n, got %v", []byte(bodyStr))
	}
}

func TestStaticDeploymentWithPath(t *testing.T) {

	deploymentBus := createBus()
	defer deploymentBus.Stop()

	// create a deployment that serves the static-site-2 fixture at
	// http://internet-golf-test.local/stuff/

	url := "http://" + BasicTestHost + "/stuff/"
	assertUrlEmpty(url, t)

	parsedUrl := db.Url{Domain: BasicTestHost, Path: "/stuff/*"}

	deploymentBus.SetupDeployment(db.DeploymentMetadata{
		// TODO: decide how asterisks work. are they implied? how would you turn
		// them off? i feel like if your path ends in a slash, you almost
		// definitely want an asterisk
		Url: parsedUrl,
	})

	deploymentBus.PutDeploymentContentByUrl(parsedUrl, db.DeploymentContent{
		ServedThingType: db.StaticFiles,
		ServedThing:     getFixturePath("static-site-2"),
	})

	bodyStr := urlToPageContent(url, t)
	if bodyStr != "stuff 2\n" {
		t.Fatalf("expected stuff 2\\n, got %v", bodyStr)
	}

	bodyStr = urlToPageContent(url+"thing.txt", t)
	if strings.Trim(bodyStr, " \n\r") != "whatever 2" {
		t.Fatalf("expected whatever 2, got %v", bodyStr)
	}
}

func TestStaticDeploymentWithPreservedPath(t *testing.T) {

	deploymentBus := createBus()
	defer deploymentBus.Stop()

	// create a deployment that serves the static-site-2 fixture at
	// http://internet-golf-test.local/stuff/

	url := "http://" + BasicTestHost + "/other-stuff/"
	assertUrlEmpty(url, t)

	parsedUrl := db.Url{Domain: BasicTestHost, Path: "/other-stuff/"}

	deploymentBus.SetupDeployment(db.DeploymentMetadata{
		Url:                  parsedUrl,
		PreserveExternalPath: true,
	})

	deploymentBus.PutDeploymentContentByUrl(parsedUrl, db.DeploymentContent{
		ServedThingType: db.StaticFiles,
		ServedThing:     getFixturePath("static-site-3"),
	})

	bodyStr := urlToPageContent(url, t)
	if bodyStr != "stuff 3\n" {
		t.Fatalf("expected stuff 3\\n, got %v", []byte(bodyStr))
	}
}
