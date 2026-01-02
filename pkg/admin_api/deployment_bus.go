package admin_api

import (
	_ "embed"
	"fmt"
	"slices"
	"strings"

	"github.com/internet-golf/internet-golf/pkg/db"
	"github.com/internet-golf/internet-golf/pkg/public_web_server"
)

func urlFromString(url string) db.Url {
	// TODO: return an error if the domain is invalid or the path has, like, # or ?
	firstSlash := strings.Index(url, "/")
	if firstSlash == -1 {
		return db.Url{Domain: url}
	} else {
		path := url[firstSlash:]
		return db.Url{Domain: url[0:firstSlash], Path: path}
	}
}

// this struct provides access to the active set of deployments and also, more
// importantly, sends those deployments to the PublicWebServer and the database
// when necessary.
type DeploymentBus struct {
	deployments []db.Deployment
	Server      public_web_server.PublicWebServer
	Db          db.Db
}

func NewDeploymentBus(server public_web_server.PublicWebServer, db db.Db) (*DeploymentBus, error) {
	deployments, err := db.GetDeployments()
	if err != nil {
		return nil, err
	}

	if err := server.DeployAll(deployments); err != nil {
		return nil, err
	}

	return &DeploymentBus{
		deployments: deployments,
		Server:      server,
		Db:          db,
	}, nil

}

func (bus *DeploymentBus) Stop() error {
	return bus.Server.Stop()
}

func (bus *DeploymentBus) persistDeployments() error {
	return bus.Db.SaveDeployments(bus.deployments)
}

// create a deployment or, if a deployment with the same name as the input
// metadata already exists, update its metadata
func (bus *DeploymentBus) SetupDeployment(metadata db.DeploymentMetadata) error {
	// TODO: make sure its URL does not overlap with any existing deployments
	// (except the one it is replacing), and that at least the domain is present
	// and a valid domain name?

	existingIndex := slices.IndexFunc(bus.deployments, func(d db.Deployment) bool {
		return d.Url.Equals(&metadata.Url)
	})

	if existingIndex == -1 {
		bus.deployments = append(bus.deployments, db.Deployment{DeploymentMetadata: metadata})
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

func (bus *DeploymentBus) getDeploymentIndexByUrl(url *db.Url) int {
	return slices.IndexFunc(bus.deployments, func(d db.Deployment) bool {
		return d.Url.Equals(url)
	})
}

func (bus *DeploymentBus) GetDeploymentByUrl(url *db.Url) (db.Deployment, error) {
	index := bus.getDeploymentIndexByUrl(url)
	if index == -1 {
		return db.Deployment{}, fmt.Errorf("Deployment with URL \"%s\" not found", url)
	}
	return bus.deployments[index], nil
}

func (bus *DeploymentBus) PutDeploymentContentByUrl(
	url db.Url, content db.DeploymentContent,
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
	index int, content db.DeploymentContent,
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
func (bus *DeploymentBus) DeleteDeployment(url db.Url) error {
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
