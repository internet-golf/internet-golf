package main

import (
	internetgolf "github.com/toBeOfUse/internet-golf/pkg"
)

func main() {
	// the core architecture of this app consists of these three actors:

	// interface to the web server that actually deploys the deployments
	deploymentServer := internetgolf.CaddyServer{}

	// object that stores (and theoretically persists) the active deployments
	// and broadcasts them to the deploymentServer when necessary
	deploymentBus := internetgolf.DeploymentBus{
		OnDeploymentChange: func(d []internetgolf.Deployment) {
			deploymentServer.Deploy(d)
		},
	}

	// api server that receives admin API requests and updates the active
	// deployments in response to them
	adminApi := internetgolf.AdminApi{Web: deploymentBus}
	adminApi.Start()
}
