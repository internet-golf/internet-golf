package internetgolf

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"slices"
)

// TODO: make into string so it can be json-serialized more durably
type CacheControlSetting int

// TODO: prefix names???
const (
	AllButHtml CacheControlSetting = iota
	Default
	Nothing
)

func (cc CacheControlSetting) String() string {
	return "CacheControlSetting(" + map[CacheControlSetting]string{
		AllButHtml: "AllButHtml",
		Default:    "Default",
		Nothing:    "Nothing",
	}[cc] + ")"
}

// TODO: make into string so it can be json-serialized more durably
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
	// and config for ignoring .html extensions?
	// also: basic auth config
	// CORS
}

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
		defer infile.Close()
		decoder := json.NewDecoder(infile)
		decoderError := decoder.Decode(&bus.deployments)
		if decoderError != nil {
			fmt.Printf("error decoding existing deployments: %v", decoderError)
		}
	}
	bus.Server.DeployAll(bus.deployments)
}

func (bus *DeploymentBus) persistDeployments() error {
	// TODO: this doesn't seem that atomic since an error after this would leave
	// an empty bus.deploymentsFile with no way to rollback. realistically
	// probably would want something like sqlite for better durability
	outfile, outfileErr := os.Create(bus.deploymentsFile)
	if outfileErr != nil {
		return outfileErr
	}
	defer outfile.Close()

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

	return bus.persistDeployments()
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

	bus.persistDeployments()

	return nil
}
