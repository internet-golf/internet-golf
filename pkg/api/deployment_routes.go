package api

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/internet-golf/internet-golf/pkg/db"
)

// TODO: update these structs to not rely on the database schema types, and be
// deployment-wise polymorphic (INT-52)

// currently unused
type DeployContainerInput struct {
	Body struct {
		// should this be json even though the static deployment input is form data??
		ContainerUrl    string `json:"containerUrl"`
		InternalAppPort int    `json:"internalAppPort"`
		Id              string `json:"id" required:"false"`
	}
}

type DeploymentCreateBody struct {
	db.DeploymentMetadata
	// the DeploymentMetadata type already has a Url field, but that uses
	// the internal Url type. this just receives a string so that the
	// internal Url type is hidden from the outside world
	Url string `json:"url"`
}
type DeploymentCreateInput struct {
	Body struct{ DeploymentCreateBody }
}

type DeployFilesBody struct {
	Url                    string        `form:"url" required:"true"`
	Contents               huma.FormFile `form:"contents" contentType:"application/gzip,application/octet-stream"`
	KeepLeadingDirectories bool          `form:"keepLeadingDirectories"`
	PreserveExistingFiles  bool          `form:"preserveExistingFiles"`
}
type DeployFilesInput struct {
	RawBody huma.MultipartFormFiles[DeployFilesBody]
}

type DeployAdminDashBody struct {
	Url string `json:"url" required:"true"`
}
type DeployAdminDashInput struct {
	Body DeployAdminDashBody
}

func (a *AdminApi) addDeploymentRoutes(api huma.API) {

	// TODO: abstract out permissions checks, which are currently very repetitive

	huma.Put(api, "/deploy/new", func(
		ctx context.Context, input *DeploymentCreateInput,
	) (*SuccessOutput, error) {
		permissions, permissionsOk := ctx.Value("permissions").(Permissions)
		if !permissionsOk {
			return nil, huma.Error500InternalServerError("Auth check failed somehow")
		}

		if !permissions.CanCreateDeployment() {
			return nil, huma.Error401Unauthorized("Not authorized to create deployments")
		}

		input.Body.DeploymentMetadata.Url = urlFromString(input.Body.Url)

		putDeploymentErr := a.web.SetupDeployment(input.Body.DeploymentMetadata)
		if putDeploymentErr != nil {
			return nil, putDeploymentErr
		}
		var output SuccessOutput
		output.Body.Success = true
		output.Body.Message = fmt.Sprintf("Created deployment with url %s", input.Body.Url)
		return &output, nil
	})

	type GetDeploymentOutput struct {
		Body struct{ db.Deployment }
	}
	huma.Get(api, "/deployment/{url}", func(ctx context.Context, input *struct {
		Url string `path:"url"`
	}) (*GetDeploymentOutput, error) {
		permissions, permissionsOk := ctx.Value("permissions").(Permissions)
		if !permissionsOk {
			return nil, fmt.Errorf("Auth check failed somehow")
		}

		url := urlFromString(input.Url)

		deployment, err := a.web.GetDeploymentByUrl(&url)
		if err != nil {
			return nil, huma.Error404NotFound(
				fmt.Sprintf("Could not find deployment with URL \"%s\"", url),
			)
		}

		if !permissions.CanViewDeployment(&deployment) {
			return nil, huma.Error403Forbidden(
				fmt.Sprintf(
					"You are not authorized to view the deployment \"%s\"",
					url,
				),
			)
		}

		var output GetDeploymentOutput
		output.Body.Deployment = deployment
		return &output, nil
	})

	huma.Put(
		api,
		"/deploy/files",
		func(
			ctx context.Context, input *DeployFilesInput,
		) (*SuccessOutput, error) {
			formData := input.RawBody.Data()

			permissions, permissionsOk := ctx.Value("permissions").(Permissions)
			if !permissionsOk {
				return nil, fmt.Errorf("Auth check failed somehow")
			}

			url := urlFromString(formData.Url)
			deployment, findDeploymentError := a.web.GetDeploymentByUrl(&url)
			if findDeploymentError != nil {
				return nil, huma.Error404NotFound(
					fmt.Sprintf("could not find deployment with URL \"%s\"", url),
				)
			}

			if !permissions.CanModifyDeployment(&deployment) {
				return nil, huma.Error403Forbidden(
					fmt.Sprintf("insufficient permissions to modify deployment \"%s\"", url),
				)
			}

			filesErr := a.web.PutStaticFilesForDeployment(deployment, formData.Contents, formData.KeepLeadingDirectories)

			if filesErr != nil {
				return nil, huma.Error500InternalServerError(
					"Error occurred while unpacking uploaded files: " + filesErr.Error(),
				)
			}

			output := SuccessOutput{}
			output.Body.Success = true
			output.Body.Message = "Updated content for " + url.String()
			return &output, nil
		})

	huma.Put(api, "/admin-dash", func(ctx context.Context, input *DeployAdminDashInput) (*SuccessOutput, error) {
		permissions, permissionsOk := ctx.Value("permissions").(Permissions)
		if !permissionsOk {
			return nil, huma.Error500InternalServerError("Auth check failed somehow")
		}

		if !permissions.CanCreateDeployment() {
			return nil, huma.Error401Unauthorized("Not authorized to create deployments")
		}

		a.web.PutAdminDash(urlFromString(input.Body.Url))

		var output SuccessOutput
		output.Body.Success = true
		output.Body.Message = "admin dashboard deployed"
		return &output, nil
	})
}
