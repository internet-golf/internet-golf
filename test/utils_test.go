// this has have a name ending in _test to be interpreted as part of this
// package apparently even though it has no actual tests

package internetgolf_test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

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
