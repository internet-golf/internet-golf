package internetgolf_test

import (
	"strconv"
	"strings"
	"testing"

	internetgolf "github.com/toBeOfUse/internet-golf/pkg"
)

func TestBasicBearerTokenFlow(t *testing.T) {
	// 1. start server
	portInt, portErr := internetgolf.GetFreePort()
	if portErr != nil {
		panic(portErr)
	}

	port := strconv.Itoa(portInt)
	stopServer := startFullServer(port)
	defer stopServer()

	// 2. get token and create deployment on localhost (as if making a request
	// from the same machine as the server)

	output := runClientCliCommand("create-token", port, t)
	if !strings.Contains(output, "Generated token:") {
		t.Fatal()
	}
	token := strings.Split(output, "\n")[2]

	runClientCliCommand("create-deployment internet-golf-test.local", port, t)

	// 3. deploy content with an api request that is routed through the host
	// we're deploying to (as if making a request from a remote machine)

	// how to check that localhost isn't being used here???
	runClientCliCommand(
		"deploy-content internet-golf-test.local --files ./fixtures/static-site --auth "+token,
		"", t,
	)

	if content := urlToPageContent("http://internet-golf-test.local", t); content != "stuff\n" {
		t.Fatalf("expected stuff\\n, got %v", []byte(content))
	}
	if content := urlToPageContent("http://internet-golf-test.local/nested/concept.txt", t); content != "fnord" {
		t.Fatalf("expected fnord, got %v", []byte(content))
	}
}
