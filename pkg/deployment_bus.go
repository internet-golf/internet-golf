package internetgolf

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
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
	// POSTing an object of this type to the API
	Name string `json:"name,omitempty"`
	Urls []Url  `json:"urls"`

	// assuming that there won't be multiple external sources...

	// for github repos, this is repoOwner/repoName or repoOwner/repoName#branch-name
	ExternalSource     string             `json:"externalSource,omitempty"`
	ExternalSourceType ExternalSourceType `json:"externalSourceType,omitempty"`

	Tags []string `json:"tags,omitempty"`

	PreserveExternalPath bool `json:"preserveExternalPath,omitempty"`

	// this is `true` for internal deployments like the one for the admin API
	DontPersist bool `json:"-"`
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

// create a deployment or, if a deployment with the same ID as the input
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

// TODO: this function is a prime target for branch test coverage

// gets the index of a stored Deployment with the provided external source,
// external source type, and name. `name` is optional (i.e. it can be an empty
// string, and this function will then just look for a deployment with a
// matching external source and external source type), unless there are multiple
// deployments with this externalSource and externalSourceType, in which case a
// name must be specified to disambiguate between them, or an error will be
// returned. an error will also be returned if no matching instance is found.
//
// if `externalSourceType == GithubRepo`, and externalSource follow the pattern
// `repoOwner/repoName#branch-name`, if no deployment can be found matching the
// full externalSource (and name, if it's present), one will be found that just
// matches the "repoOwner/repoName" part.
//
// the logic in this function is very carefully structured so that if there is
// an error, the most relevant and informative one is always returned, so that
// errors can be passed back to the end user.
func (bus *DeploymentBus) getDeploymentIndexByExternalSource(
	externalSource string, externalSourceType ExternalSourceType, name string,
) (int, error) {

	// get the branchless externalSource if it's a github repo
	lessSpecificExternalSource := externalSource
	if hashIndex := strings.Index(externalSource, "#"); externalSourceType == GithubRepo && hashIndex != -1 {
		lessSpecificExternalSource = externalSource[:hashIndex]
	}

	// the simple case is if `name` is specified
	if len(name) > 0 {
		nameIndex := bus.getDeploymentIndexByName(name)
		if nameIndex == -1 {
			return -1, fmt.Errorf("Could not find deployment with name %s", name)
		} else {
			// a deployment with `name` was found; make sure that the
			// ExternalSource and ExternalSourceType match at some level
			d := bus.deployments[nameIndex]
			if d.ExternalSourceType != externalSourceType || (d.ExternalSource != externalSource && d.ExternalSource != lessSpecificExternalSource) {
				if lessSpecificExternalSource != externalSource {
					return -1, fmt.Errorf(
						"Deployment with name %s does not match %s or %s (%s)",
						name, externalSource, lessSpecificExternalSource, externalSourceType)
				} else {
					return -1, fmt.Errorf(
						"Deployment with name %s does not match %s (%s)",
						name, externalSource, externalSourceType,
					)
				}
			} else {
				// `name` was found and `externalSource` and
				// `externalSourceType` match 👍
				return nameIndex, nil
			}
		}
	}

	// if `name` was not specified, just match on the external source variables

	// get all the indexes that match so that we can report an error if there
	// are multiple (if we didn't care about error reporting, we'd just return
	// the first result)
	var exactMatches []int
	var acceptableMatches []int
	for i, d := range bus.deployments {
		if d.ExternalSourceType == externalSourceType && d.ExternalSource == externalSource {
			exactMatches = append(exactMatches, i)
		} else if d.ExternalSourceType == externalSourceType && d.ExternalSource == lessSpecificExternalSource {
			acceptableMatches = append(acceptableMatches, i)
		}
	}

	if len(exactMatches) > 1 {
		return -1, fmt.Errorf(
			"multiple deployments found for external source \"%s (%s)\"; "+
				"please specify the name of the deployment you're looking for",
			externalSource, externalSourceType,
		)
	} else if len(exactMatches) == 1 {
		return exactMatches[0], nil
	} else {
		// there are no exact matches; fall back to the acceptable ones
		if len(acceptableMatches) > 1 {
			return -1, fmt.Errorf(
				"multiple deployments found for external source \"%s (%s)\"; "+
					"please specify the name of the deployment you're looking for",
				lessSpecificExternalSource, externalSourceType,
			)

		} else if len(acceptableMatches) == 0 {
			return -1, fmt.Errorf(
				"could not find deployment for external source \"%s (%s)\"",
				lessSpecificExternalSource, externalSourceType,
			)
		} else {
			return acceptableMatches[0], nil
		}
	}
}

// updates the struct that points to content for a given deployment, and then
// updates the public web server correspondingly, and then persists the new set
// of deployments. `name` is optional; if it's an empty string, then only
// externalSource and externalSourceType will be used to find the deployment
// whose content is being updated. however, if more than one deployment exists
// with this externalSource and externalSourceType, `name` must be non-empty, or
// an error will be returned.
func (bus *DeploymentBus) PutDeploymentContentByExternalSource(
	externalSource string, externalSourceType ExternalSourceType, name string,
	content DeploymentContent,
) error {
	fmt.Printf("updating deployment content for %s - %s\n", externalSource, externalSourceType)
	deploymentIndex, err := bus.getDeploymentIndexByExternalSource(
		externalSource, externalSourceType, name,
	)
	if err != nil {
		return err
	}
	return bus.updateDeploymentContentByIndex(deploymentIndex, content)
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
