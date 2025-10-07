package web

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/internet-golf/internet-golf/pkg/utils"
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

		addr := net.ParseIP(domain)
		if addr != nil {
			resp(400, fmt.Sprintf("cannot use ip %s as domain", domain))
			return
		}

		myIpResp, err := http.Get("https://ipv4.icanhazip.com/")
		if err != nil {
			resp(
				500,
				fmt.Sprintf(
					"Could not access icanhazip to get own IP address: %s",
					err.Error(),
				),
			)
			return
		}
		myIpBody, err := io.ReadAll(myIpResp.Body)
		myIp := strings.Trim(string(myIpBody), " \n\r")
		if err != nil || len(myIp) == 0 {
			var errorMessage string
			if err != nil {
				errorMessage = " " + err.Error()
			}
			resp(
				500,
				fmt.Sprintf(
					"Could not read own IP address from icanhazip response.%s",
					errorMessage,
				),
			)
			return
		}
		myParsedIp := net.ParseIP(myIp)

		domainIpAddresses, err := net.LookupIP(domain)
		if err != nil || len(domainIpAddresses) == 0 {
			resp(400, fmt.Sprintf("Could not lookup IP address for domain \"%s\"", domain))
			return
		}

		match := false
		for _, ip := range domainIpAddresses {
			if ip.Equal(myParsedIp) {
				match = true
				break
			}
		}

		if !match {
			resp(400, fmt.Sprintf(
				"Could not find server's IP address (%s) in addresses of \"%s\" (%+v)",
				myParsedIp, domain, domainIpAddresses,
			))
			return
		} else {
			resp(200, "OK")
		}
	})

	strPort := strconv.Itoa(port)
	server := http.Server{Addr: "localhost:" + strPort, Handler: router}
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
