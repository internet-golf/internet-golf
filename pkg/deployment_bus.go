package internetgolf

import (
	"fmt"
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

type Url struct {
	Domain string `json:"domain"`
	Path   string `json:"path,omitempty"`
}

// TODO: since this is used by the Huma API, it probably should have more docs
// and stuff in the struct tags
type DeploymentMetadata struct {
	// this has omitempty so that the name doesn't have to be specified when
	// POSTing an object of this type to the API. this is also the ID field that
	// storm uses when saving deployments
	Name string `json:"name,omitempty" storm:"id"`
	Urls []Url  `json:"urls"`

	// assuming that there won't be multiple external sources...

	// for github repos, this is repoOwner/repoName or repoOwner/repoName#branch-name
	ExternalSource     string             `json:"externalSource,omitempty"`
	ExternalSourceType ExternalSourceType `json:"externalSourceType,omitempty"`

	Tags []string `json:"tags,omitempty"`

	// if this is true and the deployment is at the path "/thing", then the
	// "/thing" in the path will be transparently passed through to the
	// underlying resource instead of being removed (which is the default)
	PreserveExternalPath bool `json:"preserveExternalPath,omitempty"`

	// this is `true` for internal deployments like the one for the admin API
	DontPersist bool `json:"-"`
}

func (d *DeploymentMetadata) Equals(e *DeploymentMetadata) bool {
	return (d.DontPersist == e.DontPersist && d.ExternalSource == e.ExternalSource &&
		d.ExternalSourceType == e.ExternalSourceType && d.Name == e.Name &&
		d.PreserveExternalPath == e.PreserveExternalPath && slices.Equal(d.Urls, e.Urls) &&
		slices.Equal(d.Tags, e.Tags))
}

type DeploymentContent struct {
	// this is false if no actual content has been added to the deployment (yet)
	HasContent bool `json:"hasContent"`
	// for static files, this is the path to a local directory; for a docker
	// container, this is a port number (?); for a redirect, this is a url or url
	// path; for a reverse proxy, this is a host and port (probably "localhost:[port]")
	ServedThing     string          `json:"servedThing"`
	ServedThingType ServedThingType `json:"servedThingType"`
}

type Deployment struct {
	DeploymentMetadata `storm:"inline"`
	DeploymentContent  `storm:"inline"`
}

type DeploymentBus struct {
	deployments []Deployment
	Server      PublicWebServer
	Db          Storage
}

// brings any persisted deployments back to life and initializes the
// DeploymentBus' Server with them
func (bus *DeploymentBus) Init() {
	deployments, err := bus.Db.GetDeployments()
	if err != nil {
		panic(err)
	}
	bus.deployments = deployments
	bus.Server.DeployAll(bus.deployments)
}

func (bus *DeploymentBus) Stop() error {
	return bus.Server.Stop()
}

func (bus *DeploymentBus) persistDeployments() error {
	return bus.Db.SaveDeployments(bus.deployments)
}

// create a deployment or, if a deployment with the same name as the input
// metadata already exists, update its metadata
func (bus *DeploymentBus) SetupDeployment(metadata DeploymentMetadata) error {
	fmt.Printf("adding deployment %+v\n", metadata)
	existingIndex := slices.IndexFunc(bus.deployments, func(d Deployment) bool {
		return d.Name == metadata.Name
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

// TODO: method to rename existing deployment (will need to enforce that names
// are unique, and aren't the empty string, etc)

func (bus *DeploymentBus) getDeploymentIndexByName(name string) int {
	return slices.IndexFunc(bus.deployments, func(d Deployment) bool {
		return d.Name == name
	})
}

func (bus *DeploymentBus) GetDeploymentByName(name string) (Deployment, error) {
	index := bus.getDeploymentIndexByName(name)
	if index == -1 {
		return Deployment{}, fmt.Errorf("Deployment with name %s not found", name)
	}
	return bus.deployments[index], nil
}

func (bus *DeploymentBus) PutDeploymentContentByName(
	name string, content DeploymentContent,
) error {
	fmt.Printf("updating deployment content %+v\n", content)
	existingIndex := bus.getDeploymentIndexByName(name)
	if existingIndex == -1 {
		return fmt.Errorf(
			"Could not find deployment with name \"%v\" to update content", name,
		)
	}
	return bus.updateDeploymentContentByIndex(existingIndex, content)
}

// updates the content of the deployment at the given index, pushes the
// deployments to the public web server, and then saves them
func (bus *DeploymentBus) updateDeploymentContentByIndex(
	index int, content DeploymentContent,
) error {
	bus.deployments[index].DeploymentContent = content
	bus.deployments[index].DeploymentContent.HasContent = true

	deploymentErr := bus.Server.DeployAll(bus.deployments)
	if deploymentErr != nil {
		// rollback to previous persisted state?
		return deploymentErr
	}

	return bus.persistDeployments()
}

// deletes the deployment from the given name, pushes the deployment set
// (without the deleted one) to the public web server, and then saves the new
// deployment set
func (bus *DeploymentBus) DeleteDeployment(name string) error {
	index := bus.getDeploymentIndexByName(name)
	if index == -1 {
		return fmt.Errorf("could not find deployment with name \"%s\" to delete it", name)
	}

	bus.deployments = slices.Delete(bus.deployments, index, index+1)

	deploymentErr := bus.Server.DeployAll(bus.deployments)
	if deploymentErr != nil {
		return deploymentErr
	}

	bus.persistDeployments()

	return nil
}
