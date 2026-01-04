package api

import (
	_ "embed"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/internet-golf/internet-golf/pkg/db"
	"github.com/internet-golf/internet-golf/pkg/public"
	"github.com/internet-golf/internet-golf/pkg/resources"
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

// the DeploymentBus handles data and config received by the admin API and
// persists them and turns them into websites.
type DeploymentBus struct {
	deployments []db.Deployment
	server      public.PublicWebServer
	db          db.Db
	files       *resources.FileManager
}

func NewDeploymentBus(server public.PublicWebServer, db db.Db, files *resources.FileManager) (*DeploymentBus, error) {
	deployments, err := db.GetDeployments()
	if err != nil {
		return nil, err
	}

	if err := server.DeployAll(deployments); err != nil {
		return nil, err
	}

	return &DeploymentBus{
		deployments: deployments,
		server:      server,
		db:          db,
		files:       files,
	}, nil

}

func (bus *DeploymentBus) Stop() error {
	return bus.server.Stop()
}

func (bus *DeploymentBus) persistDeployments() error {
	return bus.db.SaveDeployments(bus.deployments)
}

// create a deployment or, if a deployment with the same name as the input
// metadata already exists, update its metadata
func (bus *DeploymentBus) SetupDeployment(metadata db.DeploymentMetadata) error {
	// TODO: make sure its URL does not overlap with any existing deployments
	// (except the one it is replacing), and that at least the domain is present
	// and a valid domain name? also validate externalSourceType if that's a thing

	existingIndex := slices.IndexFunc(bus.deployments, func(d db.Deployment) bool {
		return d.Url.Equals(&metadata.Url)
	})

	if existingIndex == -1 {
		bus.deployments = append(bus.deployments, db.Deployment{DeploymentMetadata: metadata})
	} else {
		bus.deployments[existingIndex].DeploymentMetadata = metadata
	}

	deploymentErr := bus.server.DeployAll(bus.deployments)
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

func (bus *DeploymentBus) PutStaticFilesForDeployment(
	deployment db.Deployment, gzippedDir io.ReadSeeker, keepLeadingDirectories bool,
) error {

	outDir, extractionErr := bus.files.TarGzToDeploymentFiles(
		gzippedDir, deployment.Url.String(),
		keepLeadingDirectories,
	)

	if extractionErr != nil {
		return extractionErr
	}

	bus.PutDeploymentContentByUrl(deployment.Url, db.DeploymentContent{
		HasContent:      true,
		ServedThingType: db.StaticFiles,
		ServedThing:     outDir,
	})

	// TODO: delete the old directory after deployContent is
	// finished? presumably that'll be safe (INT-42)

	return nil
}

func (bus *DeploymentBus) PutAdminDash(url db.Url) error {
	if err := bus.SetupDeployment(db.DeploymentMetadata{
		Url:  url,
		Tags: []string{"system"},
	}); err != nil {
		return err
	}
	return bus.PutDeploymentContentByUrl(url, db.DeploymentContent{
		HasContent:      true,
		ServedThingType: db.StaticFiles,
		ServedThing:     bus.files.DashSpaPath,
		SpaMode:         true,
	})
}

// updates the content of the deployment at the given index, pushes the
// deployments to the public web server, and then saves them
func (bus *DeploymentBus) updateDeploymentContentByIndex(
	index int, content db.DeploymentContent,
) error {
	bus.deployments[index].DeploymentContent = content
	bus.deployments[index].DeploymentContent.HasContent = true

	deploymentErr := bus.server.DeployAll(bus.deployments)
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

	deploymentErr := bus.server.DeployAll(bus.deployments)
	if deploymentErr != nil {
		return deploymentErr
	}

	bus.persistDeployments()

	return nil
}
