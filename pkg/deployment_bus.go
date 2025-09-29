package internetgolf

import (
	"fmt"
	"slices"
	"strings"
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
	Github ExternalSourceType = "Github"
)

type Url struct {
	Domain string `json:"domain"`
	Path   string `json:"path,omitempty"`
}

func (u *Url) Equals(v *Url) bool {
	return u.Domain == v.Domain && u.Path == v.Path
}

func (u Url) String() string {
	return u.Domain + u.Path
}

func urlFromString(url string) Url {
	// TODO: return an error if the domain is invalid or the path has, like, # or ?
	firstSlash := strings.Index(url, "/")
	if firstSlash == -1 {
		return Url{Domain: url}
	} else {
		path := url[firstSlash:]
		return Url{Domain: url[0:firstSlash], Path: path}
	}
}

// TODO: since this is used by the Huma API, it probably should have more docs
// and stuff in the struct tags

type DeploymentMetadata struct {
	Url Url `json:"" storm:"id"`

	// assuming that there won't be multiple external sources...
	// TODO: probably move this to the auth section?
	// for github repos, the ExternalSource has the format "repoOwner/repoName"
	// or "repoOwner/repoName#branch-name"
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

// TODO: create struct just for API that does not have "hasContent" (since it's
// an internal value)
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
	Db          Db
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
	// TODO: make sure its URL does not overlap with any existing deployments
	// (except the one it is replacing), and that at least the domain is present
	// and a valid domain name?

	existingIndex := slices.IndexFunc(bus.deployments, func(d Deployment) bool {
		return d.Url.Equals(&metadata.Url)
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

// TODO: method to change the URL of a deployment?

func (bus *DeploymentBus) getDeploymentIndexByUrl(url *Url) int {
	return slices.IndexFunc(bus.deployments, func(d Deployment) bool {
		return d.Url.Equals(url)
	})
}

func (bus *DeploymentBus) GetDeploymentByUrl(url *Url) (Deployment, error) {
	index := bus.getDeploymentIndexByUrl(url)
	if index == -1 {
		return Deployment{}, fmt.Errorf("Deployment with URL \"%s\" not found", url)
	}
	return bus.deployments[index], nil
}

func (bus *DeploymentBus) PutDeploymentContentByUrl(
	url Url, content DeploymentContent,
) error {
	existingIndex := bus.getDeploymentIndexByUrl(&url)
	if existingIndex == -1 {
		return fmt.Errorf(
			"Could not find deployment with url \"%s\" to update content", url,
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
func (bus *DeploymentBus) DeleteDeployment(url Url) error {
	index := bus.getDeploymentIndexByUrl(&url)
	if index == -1 {
		return fmt.Errorf("could not find deployment with URL \"%s\" to delete it", url)
	}

	bus.deployments = slices.Delete(bus.deployments, index, index+1)

	deploymentErr := bus.Server.DeployAll(bus.deployments)
	if deploymentErr != nil {
		return deploymentErr
	}

	bus.persistDeployments()

	return nil
}
