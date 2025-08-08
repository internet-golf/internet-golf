package internetgolf

import (
	"fmt"
	"slices"
)

type CacheControlSetting int

const (
	AllButHtml CacheControlSetting = iota
	All
	Nothing
)

func (cc CacheControlSetting) String() string {
	return "CacheControlSetting(" + map[CacheControlSetting]string{
		AllButHtml: "AllButHtml",
		All:        "All",
		Nothing:    "Nothing",
	}[cc] + ")"
}

type DeploymentSettings struct {
	// defaults to AllButHtml since that's the 0 value for the enum
	CacheControlMode CacheControlSetting `json:"cacheControlMode"`
	// 404 page address for static sites?
	// TODO: ensure config exists for SPAs that serves the same file for all paths
	// also: basic auth config
}

type Deployment struct {
	Id          string             `json:"id"`
	Matcher     string             `json:"matcher"`
	ResourceUri string             `json:"resourceUri"`
	Settings    DeploymentSettings `json:"settings"`
}

type DeploymentBus struct {
	deployments        []Deployment
	OnDeploymentChange func([]Deployment)
}

func (bus DeploymentBus) PutDeployment(newDeployment Deployment) error {
	fmt.Printf("adding deployment %+v\n", newDeployment)
	newDeployments := slices.DeleteFunc(bus.deployments, func(d Deployment) bool {
		return d.Id == newDeployment.Id
	})
	bus.deployments = append(newDeployments, newDeployment)
	bus.OnDeploymentChange(bus.deployments)

	return nil
}

func (bus DeploymentBus) DeleteDeployment(id string) error {
	fmt.Printf("removing deployment with id %v\n", id)
	bus.deployments = slices.DeleteFunc(bus.deployments, func(d Deployment) bool {
		return d.Id == id
	})
	bus.OnDeploymentChange(bus.deployments)

	return nil
}
