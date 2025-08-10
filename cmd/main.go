package main

import (
	internetgolf "github.com/toBeOfUse/internet-golf/pkg"
)

func main() {

	// TODO: receive non-default data directory as command-line argument and
	// pass in here
	settings := internetgolf.StorageSettings{}
	settings.Init("")

	// the core architecture of this app consists of these three actors:

	// interface to the web server that actually deploys the deployments
	deploymentServer := internetgolf.CaddyServer{}
	// TODO: from cli arg
	deploymentServer.Settings.LocalOnly = true

	// object that (persistently) stores the active deployments and broadcasts
	// them to the deploymentServer when necessary
	deploymentBus := internetgolf.DeploymentBus{
		// TODO: why does & work here ???
		Server:          &deploymentServer,
		StorageSettings: settings,
	}
	// this initializes the deployment bus and the server that it controls
	deploymentBus.Init()

	// api server that receives admin API requests and updates the active
	// deployments in response to them
	adminApi := internetgolf.AdminApi{
		Web:             deploymentBus,
		StorageSettings: settings,
	}
	adminApi.Start()
}
