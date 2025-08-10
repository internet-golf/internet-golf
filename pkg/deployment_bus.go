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

type LocalResourceType int

const (
	Files LocalResourceType = iota
	DockerContainer
)

func (r LocalResourceType) String() string {
	return "LocalResourceType(" + map[LocalResourceType]string{
		Files:           "Files",
		DockerContainer: "DockerContainer",
	}[r] + ")"
}

type DeploymentSettings struct {
	// defaults to AllButHtml since that's the 0 value for the enum
	CacheControlMode CacheControlSetting `json:"cacheControlMode"`
	// 404 page address for static sites?
	// TODO: ensure config exists for SPAs that serves the same file for all paths
	// also: basic auth config
}

// why did i make this json serializable ? it should probably be kept
// internal
type Deployment struct {
	// TODO: do i need or want an ID? or does the Matcher fill that role effectively?
	Id                   string             `json:"id"`
	Matcher              string             `json:"matcher"`
	LocalResourceLocator string             `json:"localResourceLocator"`
	LocalResourceType    LocalResourceType  `json:"localResourceType"`
	Settings             DeploymentSettings `json:"settings"`
}

type DeploymentBus struct {
	deployments []Deployment
	DeployAll   func([]Deployment)
}

func (bus DeploymentBus) PutDeployment(newDeployment Deployment) error {
	fmt.Printf("adding deployment %+v\n", newDeployment)
	newDeployments := slices.DeleteFunc(bus.deployments, func(d Deployment) bool {
		return d.Id == newDeployment.Id
	})
	bus.deployments = append(newDeployments, newDeployment)
	bus.DeployAll(bus.deployments)

	return nil
}

func (bus DeploymentBus) DeleteDeployment(id string) error {
	fmt.Printf("removing deployment with id %v\n", id)
	bus.deployments = slices.DeleteFunc(bus.deployments, func(d Deployment) bool {
		return d.Id == id
	})
	bus.DeployAll(bus.deployments)

	return nil
}
