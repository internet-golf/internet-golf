package internetgolf

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"slices"
)

type SiteResourceType string

const (
	StaticFiles     SiteResourceType = "StaticFiles"
	DockerContainer SiteResourceType = "DockerContainer"
	Redirect        SiteResourceType = "Redirect"
)

type DeploymentSettings struct {
}

type Deployment struct {
	// TODO: do i need or want an ID? or does the Matcher fill that role effectively?
	Id                  string             `json:"id"`
	Matcher             string             `json:"matcher"`
	SiteResourceType    SiteResourceType   `json:"siteResourceType"`
	SiteResourceLocator string             `json:"siteResourceLocator"`
	Settings            DeploymentSettings `json:"settings"`
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
	if infile, infileErr := os.Open(bus.deploymentsFile); infileErr == nil {
		defer infile.Close()
		decoder := json.NewDecoder(infile)
		decoderError := decoder.Decode(&bus.deployments)
		if decoderError != nil {
			fmt.Printf("error decoding existing deployments: %v", decoderError)
		}
	} else if !errors.Is(infileErr, os.ErrNotExist) {
		panic(
			fmt.Sprintf(
				"could not initialize deployment bus; \"%v\" not openable",
				bus.deploymentsFile,
			),
		)
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
