package main

import (
	internetgolf "github.com/toBeOfUse/internet-golf/pkg"
)

func main() {

	// TODO: receive non-default settings as command-line argument and pass in here
	settings := internetgolf.StorageSettings{}
	settings.Init("")
	adminApiPort := "8888"

	// the core architecture of this app consists of these three actors:

	// 1. interface to the web server that actually deploys the deployments
	deploymentServer := internetgolf.CaddyServer{}
	// TODO: from cli arg
	deploymentServer.Settings.LocalOnly = true

	// 2. object that (persistently) stores the active deployments and
	// broadcasts them to the deploymentServer when necessary
	deploymentBus := internetgolf.DeploymentBus{
		// TODO: why does & work here ???
		Server:          &deploymentServer,
		StorageSettings: settings,
	}
	// this initializes the deployment bus and the server that it controls
	deploymentBus.Init()

	// 3. api server that receives admin API requests and updates the active
	// deployments in response to them
	adminApi := internetgolf.AdminApi{
		Web:             deploymentBus,
		StorageSettings: settings,
		Port:            adminApiPort,
	}
	// create a deployment for the admin api (slightly premature)
	deploymentBus.PutDeployment(internetgolf.Deployment{
		Id: "internet-golf-admin-api",
		// this could perhaps be configurable
		Matcher:             "/internet--golf--admin",
		SiteResourceType:    internetgolf.ReverseProxy,
		SiteResourceLocator: "localhost:" + adminApiPort,
		Persist:             false,
	})
	adminApi.Start()
}
