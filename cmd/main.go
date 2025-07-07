package main

import (
	deppy "github.com/toBeOfUse/deployment-agent/pkg"
)

func main() {

	deploymentServer := deppy.CaddyServer{}
	deploymentBus := deppy.DeploymentBus{
		OnDeploymentChange: func(d []deppy.Deployment) {
			deploymentServer.Deploy(d)
		},
	}
	adminApi := deppy.AdminApi{Web: deploymentBus}
	adminApi.Start()
}
