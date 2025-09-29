package internetgolf_test

import (
	"strconv"
	"strings"
	"testing"

	internetgolf "github.com/toBeOfUse/internet-golf/pkg"
)

func TestBasicBearerTokenFlow(t *testing.T) {
	portInt, portErr := internetgolf.GetFreePort()
	if portErr != nil {
		panic(portErr)
	}

	port := strconv.Itoa(portInt)
	stopServer := startFullServer(port)
	defer stopServer()

	output := runClientCliCommand("create-token", port, t)
	if !strings.Contains(output, "Generated token:") {
		t.Fatal()
	}

	// TODO:
	// create deployment through cli request to localhost
	// deploy content to it using token; request should be routed through the
	// actual deployment host
}
