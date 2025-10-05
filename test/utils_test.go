// this has have a name ending in _test to be interpreted as part of this
// package apparently even though it has no actual tests

package internetgolf_test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"

	golfsdk "github.com/internet-golf/internet-golf/client-sdk"
	"github.com/txn2/txeh"
)

const (
	BasicTestHost = "internet-golf-test.local"
)

func TestMain(m *testing.M) {
	setupHosts()
	code := m.Run()
	busCleanup()
	os.Exit(code)
}

func urlToPageContent(url string, t *testing.T) string {
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("request to %s failed", url)
		return ""
	}
	body, bodyErr := io.ReadAll(resp.Body)
	if bodyErr != nil {
		t.Fatal(bodyErr)
	}
	bodyStr := string(body)
	return bodyStr
}

func createClient(url string) *golfsdk.APIClient {
	return golfsdk.NewAPIClient(&golfsdk.Configuration{
		UserAgent: "InternetGolfClient",
		Servers: golfsdk.ServerConfigurations{
			{URL: url},
		},
	})
}

type CmdTeeWriter struct {
	savedOutput []byte
}

func (so *CmdTeeWriter) Write(p []byte) (n int, err error) {
	so.savedOutput = append(so.savedOutput, p...)
	return os.Stdout.Write(p)
}

// runs the provided command; its output is both sent to the terminal and
// returned as a string
func execWithTeedOutput(cmd *exec.Cmd) (string, error) {
	// https://stackoverflow.com/a/72809770/3962267
	var so CmdTeeWriter
	cmd.Stdout = &so
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	} else {
		return string(so.savedOutput), nil
	}
}

// runs the client cli. return the standard output of the command so it can be
// inspected.
//
// if apiPort is a non-empty string, localhost:[apiPort] will be specified as
// the api url; otherwise, if apiPort is empty, no api url will be automatically
// specified by this function.
func runClientCliCommand(command string, apiPort string, t *testing.T) string {
	var fullCommand string
	if len(apiPort) > 0 {
		fullCommand = "run ../client-cmd --api-url http://localhost:" + apiPort + " " + command
	} else {
		fullCommand = "run ../client-cmd " + command
	}
	commandSplitOnQuotes := strings.Split(fullCommand, "\"")
	commandParts := []string{}
	for i, quotePart := range commandSplitOnQuotes {
		if i%2 == 0 {
			// we are outside of the quotes
			partSplitOnSpaces := strings.Split(quotePart, " ")
			for _, spacePart := range partSplitOnSpaces {
				if len(spacePart) > 0 {
					commandParts = append(commandParts, spacePart)
				}
			}
		} else {
			// we are inside the quotes
			commandParts = append(commandParts, quotePart)
		}
	}
	cmd := exec.Command("go", commandParts...)
	fmt.Printf("Running client command: %s\n", cmd.String())
	output, err := execWithTeedOutput(cmd)
	if err != nil {
		t.Fatal(err)
	}
	return output
}

func setupHosts() {
	// this does not actually appear to create a new hosts file but rather
	// loads the existing one
	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		panic(err)
	}

	requiredHosts := []string{BasicTestHost}

	// add each required host that does not already exist
	for _, host := range requiredHosts {
		exists, ipAddress, _ := hosts.HostAddressLookup(host, txeh.IPFamilyV4)
		if !exists || ipAddress != "127.0.0.1" {
			hosts.AddHost("127.0.0.1", host)
			if saveErr := hosts.Save(); saveErr != nil {
				fmt.Printf("Could not add host %s to hosts file.\n", host)
				fmt.Println("Please either add it and 127.0.0.1 to your hosts file yourself, or ")
				fmt.Println("run this command with admin privileges.")
				panic(saveErr)
			}
		}
	}
}
