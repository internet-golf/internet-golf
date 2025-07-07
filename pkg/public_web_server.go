package deppy

import (
	"fmt"
)

type PublicWebServer interface {
	Deploy([]Deployment)
}

type CaddyServer struct{}

func (c CaddyServer) Deploy(deployments []Deployment) error {
	fmt.Printf("deploying %+v\n", deployments)
	return nil
}
