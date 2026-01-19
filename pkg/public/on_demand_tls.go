package public

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/internet-golf/internet-golf/pkg/utils"
	"golang.org/x/net/publicsuffix"
)

var onDemandTlsEndpointPath = "/approve-tls"

// returns server, server port (as string), error
func createTlsApprovalServer() (*http.Server, string, error) {
	port, portErr := utils.GetFreePort()
	if portErr != nil {
		return nil, "", portErr
	}

	router := http.NewServeMux()
	router.HandleFunc(onDemandTlsEndpointPath, func(w http.ResponseWriter, r *http.Request) {

		resp := func(status int, message string) {
			w.WriteHeader(status)
			w.Header().Add("Content-Type", "text/plain")
			w.Write([]byte(message))
		}

		domain := r.URL.Query().Get("domain")
		if len(domain) == 0 {
			resp(400, "malformed request; no \"domain\" query")
			return
		}

		apex, apexErr := publicsuffix.EffectiveTLDPlusOne(domain)
		if apexErr != nil || domain != apex {
			resp(400, fmt.Sprintf("domain %q is not an apex domain", domain))
			return
		}

		resp(200, "OK")
	})

	strPort := strconv.Itoa(port)
	server := http.Server{Addr: "localhost:" + strPort, Handler: router}
	fmt.Printf("on-demand tls server at localhost:%s\n", strPort)
	return &server, strPort, nil
}

func getTlsConfig(approvalServerPort string) utils.JsonObj {
	return utils.JsonObj{
		"automation": utils.JsonObj{
			"policies": []utils.JsonObj{
				{"on_demand": true},
			},
			"on_demand": utils.JsonObj{
				"permission": utils.JsonObj{
					"module":   "http",
					"endpoint": "http://localhost:" + approvalServerPort + onDemandTlsEndpointPath,
				},
			},
		},
	}
}

func getOnDemandTls() (onDemandTls, error) {
	server, port, err := createTlsApprovalServer()
	if err != nil {
		return onDemandTls{}, err
	}
	return onDemandTls{
		tlsApprovalServer:  server,
		caddyTlsConfig:     getTlsConfig(port),
		approvalServerPort: port,
	}, nil
}

type onDemandTls struct {
	caddyTlsConfig     utils.JsonObj
	tlsApprovalServer  *http.Server
	approvalServerPort string
}
