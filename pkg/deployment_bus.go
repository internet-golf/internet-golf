package internetgolf

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"slices"
)

type ServedThingType string

const (
	StaticFiles     ServedThingType = "StaticFiles"
	DockerContainer ServedThingType = "DockerContainer"
	Redirect        ServedThingType = "Redirect"
	// low-level deployment type; currently just used to expose the admin api
	ReverseProxy ServedThingType = "ReverseProxy"
)

type ExternalSourceType string

const (
	GithubRepo ExternalSourceType = "GithubRepo"
)

type DeploymentSettings struct {
}

// TODO: break into internal vs. configurable
type DeploymentMetadata struct {
	// public URL of the deployment. also serves as the deployment's unique ID.
	// contains a hostname with an optional url path. examples: "mitch.website"
	// "mitch.website/path" "mitch.website/path*"
	Url string `json:"url"`

	// for github repos, this is repoOwner/repoName or repoOwner/repoName#branch-name
	ExternalSource     string             `json:"externalSource,omitempty"`
	ExternalSourceType ExternalSourceType `json:"externalSourceType,omitempty"`

	Tags     []string           `json:"tags,omitempty"`
	Settings DeploymentSettings `json:"settings,omitempty"`

	// used for internal deployments like the one for the admin API
	DontPersist bool `json:"-"`
}

type DeploymentContent struct {
	HasContent bool `json:"hasContent"`
	// for static files, this is the path to a local directory; for a docker
	// container, this is a port number (?); for a redirect, this is a url or url
	// path; for a reverse proxy, this is a host and port (probably "localhost:[port]")
	ServedThing     string          `json:"servedThing"`
	ServedThingType ServedThingType `json:"servedThingType"`
}

type Deployment struct {
	DeploymentMetadata
	DeploymentContent
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
	// source of truth for where deployments are persisted to
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
	persistable := slices.Collect(func(yield func(Deployment) bool) {
		for _, d := range bus.deployments {
			if !d.DontPersist {
				if !yield(d) {
					return
				}
			}
		}
	})
	jsonErr := encoder.Encode(persistable)
	if jsonErr != nil {
		return jsonErr
	}

	return nil
}

// create a deployment or update its metadata
func (bus *DeploymentBus) PutDeploymentMetadata(metadata DeploymentMetadata) error {
	fmt.Printf("adding deployment %+v\n", metadata)
	existingIndex := slices.IndexFunc(bus.deployments, func(d Deployment) bool {
		return d.Url == metadata.Url
	})

	if existingIndex == -1 {
		bus.deployments = append(bus.deployments, Deployment{DeploymentMetadata: metadata})
	} else {
		bus.deployments[existingIndex].DeploymentMetadata = metadata
	}

	deploymentErr := bus.Server.DeployAll(bus.deployments)
	if deploymentErr != nil {
		// rollback to persisted state?
		return deploymentErr
	}

	return bus.persistDeployments()
}

func (bus *DeploymentBus) PutDeploymentContent(url string, content DeploymentContent) error {
	fmt.Printf("updating deployment content %+v\n", content)
	existingIndex := slices.IndexFunc(bus.deployments, func(d Deployment) bool {
		return d.Url == url
	})

	if existingIndex == -1 {
		return fmt.Errorf("Could not find deployment with URL \"%v\" to update content", url)
	} else {
		bus.deployments[existingIndex].DeploymentContent = content
		bus.deployments[existingIndex].DeploymentContent.HasContent = true
	}

	deploymentErr := bus.Server.DeployAll(bus.deployments)
	if deploymentErr != nil {
		// rollback to persisted state?
		return deploymentErr
	}

	return bus.persistDeployments()
}

func (bus *DeploymentBus) DeleteDeployment(url string) error {
	fmt.Printf("removing deployment with url %v\n", url)
	bus.deployments = slices.DeleteFunc(bus.deployments, func(d Deployment) bool {
		return d.Url == url
	})

	deploymentErr := bus.Server.DeployAll(bus.deployments)
	if deploymentErr != nil {
		return deploymentErr
	}

	bus.persistDeployments()

	return nil
}
