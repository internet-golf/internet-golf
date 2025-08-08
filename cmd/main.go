package main

import (
	deppy "github.com/toBeOfUse/deployment-agent/pkg"
)

func main() {
	// the core architecture of this app consists of these three actors:

	// interface to the web server that actually deploys the deployments
	deploymentServer := deppy.CaddyServer{}

	// object that stores (and theoretically persists) the active deployments
	// and broadcasts them to the deploymentServer when necessary
	deploymentBus := deppy.DeploymentBus{
		OnDeploymentChange: func(d []deppy.Deployment) {
			deploymentServer.Deploy(d)
		},
	}

	// api server that receives admin API requests and updates the active
	// deployments in response to them
	adminApi := deppy.AdminApi{Web: deploymentBus}
	adminApi.Start()
}
