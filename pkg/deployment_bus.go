package internetgolf

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
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
	deployments     []Deployment
	deploymentsFile string
	Server          PublicWebServer
	StorageSettings StorageSettings
}

// brings any persisted deployments back to life and initializes the
// DeploymentBus' Server with them
func (bus *DeploymentBus) Init() {
	bus.deploymentsFile = path.Join(bus.StorageSettings.DataDirectory, "deployments.json")
	infile, infileErr := os.Open(bus.deploymentsFile)
	// TODO: handle file not existing yet differently from other errors
	if infileErr == nil {
		decoder := json.NewDecoder(infile)
		decoderError := decoder.Decode(&bus.deployments)
		if decoderError != nil {
			fmt.Printf("error decoding existing deployments: %v", decoderError)
		}
	}
	bus.Server.DeployAll(bus.deployments)
}

func (bus *DeploymentBus) PersistDeployments() error {
	// TODO: this doesn't seem that atomic since an error after this would leave
	// an empty bus.deploymentsFile with no way to rollback
	outfile, outfileErr := os.Create(bus.deploymentsFile)
	if outfileErr != nil {
		return outfileErr
	}

	encoder := json.NewEncoder(outfile)
	jsonErr := encoder.Encode(bus.deployments)
	if jsonErr != nil {
		return jsonErr
	}

	return nil
}

func (bus *DeploymentBus) PutDeployment(newDeployment Deployment) error {
	fmt.Printf("adding deployment %+v\n", newDeployment)
	newDeployments := slices.DeleteFunc(bus.deployments, func(d Deployment) bool {
		return d.Id == newDeployment.Id
	})
	bus.deployments = append(newDeployments, newDeployment)

	deploymentErr := bus.Server.DeployAll(bus.deployments)
	if deploymentErr != nil {
		return deploymentErr
	}

	return bus.PersistDeployments()
}

func (bus *DeploymentBus) DeleteDeployment(id string) error {
	fmt.Printf("removing deployment with id %v\n", id)
	bus.deployments = slices.DeleteFunc(bus.deployments, func(d Deployment) bool {
		return d.Id == id
	})

	deploymentErr := bus.Server.DeployAll(bus.deployments)
	if deploymentErr != nil {
		return deploymentErr
	}

	bus.PersistDeployments()

	return nil
}
