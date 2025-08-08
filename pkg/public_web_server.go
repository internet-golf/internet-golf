package deppy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type PublicWebServer interface {
	Deploy([]Deployment)
}

type CaddyServer struct{}

func getCaddyStaticRouteBoilerplate(publicUrl string, filePath string) string {
	// this is unfortunately much easier than creating a struct that can handle
	// all of these object types
	return fmt.Sprintf(`{
		"match": [{ "host": ["%s"] }],
		"handle": [
		{
			"handler": "subroute",
			"routes": [
			{
				"handle": [
				{
					"handler": "vars",
					"root": "%s"
				}
				]
			},
			{
				"handle": [
				{
					"handler": "headers",
					"response": { "set": { "Cache-Control": ["max-age=0,no-store"] } }
				}
				],
				"match": [{ "path": ["*/"] }]
			},
			{
				"handle": [
				{
					"handler": "encode",
					"encodings": { "gzip": {}, "zstd": {} },
					"prefer": ["zstd", "gzip"]
				},
				{"browse": {"file_limit": 1000}, "handler": "file_server" }
				]
			}
			]
		}
		],
		"terminal": true
	}`, publicUrl, filePath)
}

func (c CaddyServer) Deploy(deployments []Deployment) error {
	fmt.Printf("deploying %+v\n", deployments)
	// TODO: implement https://caddyserver.com/docs/api#concurrent-config-changes

	var configs []string
	for _, deployment := range deployments {
		uriParts := strings.Split(deployment.ResourceUri, "://")
		if uriParts[0] == "file" {
			configs = append(configs, getCaddyStaticRouteBoilerplate(deployment.Matcher, uriParts[1]))
		}
		// TODO: docker cases
	}

	config := fmt.Sprintf(`{
		"apps": {
			"http": {
				"servers": {
					"srv0": {
						"automatic_https": {"disable": true},
						"listen": [":443", ":80", ":8989"],
						"routes": [%s]
					}
				}
			}
		}
	}`, strings.Join(configs, ",\n"))

	fmt.Print(config + "\n")

	resp, err := http.Post("http://localhost:2019/config/", "application/json", strings.NewReader(config))
	if err != nil {
		panic(err)
	}

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Print(sb)

	return nil
}
