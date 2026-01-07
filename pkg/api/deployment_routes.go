package api

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"slices"

	"github.com/danielgtaylor/huma/v2"
	"github.com/internet-golf/internet-golf/pkg/db"
)

// input types ======================

type DeployFilesBody struct {
	Url                    string        `form:"url" required:"true" doc:"The URL of the deployment that you're updating." example:"mysite.mydomain.com"`
	Contents               huma.FormFile `form:"contents" contentType:"application/gzip,application/octet-stream" doc:"A .tar.gz that contains the files to be deployed."`
	KeepLeadingDirectories bool          `form:"keepLeadingDirectories" doc:"By default, if you upload a .tar.gz whose contents are all in one folder, the contents of that folder will be used instead of the folder itself. For example, if you upload a folder called 'dist' for the deployment 'mysite.com', the URL of your site content will not be at 'mysite.com/dist'. Setting this to true turns off that auto-unpacking." default:"false"`
	PreserveExistingFiles  bool          `form:"preserveExistingFiles" doc:"Leave the existing files for the current deployment in place instead of completely replacing them."`
}
type DeployFilesInput struct {
	RawBody huma.MultipartFormFiles[DeployFilesBody]
}

type DeployAliasBody struct {
	Url string `form:"url" required:"true" doc:"The URL of the deployment that you're updating." example:"mysite.mydomain.com"`
	AliasBase
}
type DeployAliasInput struct {
	Body DeployAliasBody
}

type DeployAdminDashBody struct {
	Url string `json:"url" required:"true" doc:"The URL that you want to deploy the admin dashboard to." example:"dash.mydomain.com"`
}
type DeployAdminDashInput struct {
	Body DeployAdminDashBody
}

// this input type cheats and reuses the output type declared below
type DeploymentCreateInput struct {
	Body struct{ DeploymentBase }
}

// output types ============================

type DeploymentBase struct {
	Url string `json:"url" doc:"URL that this deployment will appear at. The DNS for the domain has to be set up first." example:"mysite.mydomain.com"`

	// assuming that there won't be multiple external sources...
	ExternalSource     string `json:"externalSource,omitempty" required:"false" doc:"Original repository for this deployment's source. Can include a branch name." example:"user/repo or user/repo#branch-name"`
	ExternalSourceType string `json:"externalSourceType,omitempty" required:"false" doc:"Place where the original repository lives."`

	Tags []string `json:"tags" required:"false" doc:"Tags used for metadata."`

	PreserveExternalPath bool `json:"preserveExternalPath" required:"false" doc:"if this is true and the deployment url has a path like \"/thing\", then the \"/thing\" in the path will be transparently passed through to the underlying resource instead of being removed (which is the default)"`
}

// this could go in DeploymentBase if DeploymentCreateInput didn't cheat and use
// it for input
type DeploymentType struct {
	Type string `json:"type" enum:"StaticSite,Alias,Empty" doc:"Type of deployment contents."`
}

type StaticSiteBase struct {
	// these values are pointers so that they will be properly omitted from the
	// JSON response if not set by the API handler (which will happen when
	// creating a DeploymentBody for a non-static-site deployment)
	ServerContentLocation *string `json:"serverContentLocation,omitempty" doc:"The path to this deployment's files on the server."`
	SpaMode               *bool   `json:"spaMode,omitempty" doc:"Whether this deployment is set up to support a Single Page App by using /index.html as a fallback for all requests."`
}

type AliasBase struct {
	// these values are pointers so that they will be properly omitted from the
	// JSON response if not set by the API handler (which will happen when
	// creating a DeploymentBody for a non-alias deployment)
	AliasedTo *string `json:"aliasedTo,omitempty" doc:"The URL that this deployment is an alias for."`
	Redirect  *bool   `json:"redirect,omitempty" doc:"If this is true, visitors to this deployment's URL will be completely redirected to the URL that this alias is for."`
}

// this mostly exists to make absolutely sure that the different deployment base
// types can be distinguished between by e.g. OpenAPI validation
type EmptyBase struct {
	NoContentYet bool `json:"noContentYet" doc:"Set to true to indicate that this deployment has not yet been set up."`
}

// this is the type that is returned for a deployment from the api. it combines
// the properties of all the possible deployment types. to let api consumers use
// this more safely, a custom schema (`deploymentSchema`) is created for it in
// the API code and inserted into the OpenAPI specification; see below
type DeploymentModel struct {
	DeploymentBase
	DeploymentType
	AliasBase
	StaticSiteBase
	EmptyBase
}
type GetDeploymentOutput struct {
	Body DeploymentModel
}

type GetDeploymentsOutput struct {
	Body struct {
		Deployments []DeploymentModel `json:"deployments" required:"true"`
	}
}

// api code =================================

func deploymentToApiModel(deployment db.Deployment) (DeploymentModel, error) {
	var output DeploymentModel
	output.DeploymentBase = DeploymentBase{
		Url:                  deployment.Url.String(),
		ExternalSource:       deployment.ExternalSource,
		ExternalSourceType:   string(deployment.ExternalSourceType),
		Tags:                 deployment.Tags,
		PreserveExternalPath: deployment.PreserveExternalPath,
	}
	if deployment.ServedThingType == db.StaticFiles {
		output.Type = "StaticSite"
		output.StaticSiteBase.ServerContentLocation = &deployment.ServedThing
		output.StaticSiteBase.SpaMode = &deployment.SpaMode
	} else if deployment.ServedThingType == db.Alias {
		output.Type = "Alias"
		aliasedTo := deployment.AliasedTo.String()
		output.AliasBase.AliasedTo = &aliasedTo
		output.AliasBase.Redirect = &deployment.Redirect
	} else if len(deployment.ServedThingType) == 0 {
		output.Type = "Empty"
		output.EmptyBase.NoContentYet = true
	} else {
		return DeploymentModel{}, fmt.Errorf("Deployment type %v not supported by the API", deployment.ServedThingType)
	}
	return output, nil
}

func (a *AdminApi) addDeploymentRoutes(api huma.API) {

	// TODO: abstract out permissions checks, which are currently very repetitive

	huma.Register(api, huma.Operation{
		OperationID: "CreateDeployment",
		Description: "Create a new deployment.",
		Method:      http.MethodPut,
		Path:        "/deploy/new",
	}, func(
		ctx context.Context, input *DeploymentCreateInput,
	) (*SuccessOutput, error) {
		permissions, permissionsOk := ctx.Value("permissions").(Permissions)
		if !permissionsOk {
			return nil, huma.Error500InternalServerError("Auth check failed somehow")
		}

		if !permissions.CanCreateDeployment() {
			return nil, huma.Error401Unauthorized("Not authorized to create deployments")
		}

		tags := input.Body.Tags
		if tags == nil {
			tags = []string{}
		}
		putDeploymentErr := a.web.SetupDeployment(db.DeploymentMetadata{
			Url:                  urlFromString(input.Body.Url),
			ExternalSource:       input.Body.ExternalSource,
			ExternalSourceType:   db.ExternalSourceType(input.Body.ExternalSourceType),
			Tags:                 tags,
			PreserveExternalPath: input.Body.PreserveExternalPath,
			DontPersist:          false,
		})
		if putDeploymentErr != nil {
			return nil, putDeploymentErr
		}
		var output SuccessOutput
		output.Body.Success = true
		output.Body.Message = fmt.Sprintf("Created deployment with url %s", input.Body.Url)
		return &output, nil
	})

	registry := api.OpenAPI().Components.Schemas

	// these types and this schema match the DeploymentModel type except it
	// narrows the structs for the different types of deployments down to one
	// for each instance of the response object, with "type" as the field that
	// indicates which one you're getting. go doesn't really have this feature,
	// but it produces a discriminated union if you use it to generate
	// typescript

	type StaticSiteDeployment struct {
		DeploymentBase
		DeploymentType
		StaticSiteBase
	}
	type AliasDeployment struct {
		DeploymentBase
		DeploymentType
		AliasBase
	}
	type EmptyDeployment struct {
		DeploymentBase
		DeploymentType
		EmptyBase
	}
	deploymentSchema := &huma.Schema{
		OneOf: []*huma.Schema{
			registry.Schema(reflect.TypeFor[StaticSiteDeployment](), true, ""),
			registry.Schema(reflect.TypeFor[AliasDeployment](), true, ""),
			registry.Schema(reflect.TypeFor[EmptyDeployment](), true, ""),
		},
		Discriminator: &huma.Discriminator{
			PropertyName: "type", Mapping: map[string]string{
				"StaticSite": "#/components/schemas/StaticSiteDeployment",
				"Alias":      "#/components/schemas/AliasDeployment",
				"Empty":      "#/components/schemas/EmptyDeployment",
			},
		},
	}
	deploymentsSchema := &huma.Schema{
		Properties: map[string]*huma.Schema{
			"deployments": {
				Type:  "array",
				Items: deploymentSchema,
			},
		},
	}

	huma.Register(api, huma.Operation{
		OperationID: "GetDeployments",
		Description: "Retrieve all active deployments.",
		Method:      http.MethodGet,
		Path:        "/deployments",
		Responses: map[string]*huma.Response{
			"200": {
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: deploymentsSchema,
					},
				},
			},
		},
	}, func(ctx context.Context, input *struct{}) (*GetDeploymentsOutput, error) {
		permissions, permissionsOk := ctx.Value("permissions").(Permissions)
		if !permissionsOk {
			return nil, fmt.Errorf("Auth check failed somehow")
		}

		// make a copy of the deployment bus' deployments
		deployments := append([]db.Deployment{}, a.web.deployments...)
		// delete the deployments that the current user shouldn't be able to see
		deployments = slices.DeleteFunc(deployments, func(d db.Deployment) bool {
			return !permissions.CanViewDeployment(&d)
		})

		var output GetDeploymentsOutput
		output.Body.Deployments = []DeploymentModel{}
		for _, d := range deployments {
			model, error := deploymentToApiModel(d)
			if error != nil {
				fmt.Println(error.Error())
			} else {
				output.Body.Deployments = append(output.Body.Deployments, model)
			}
		}

		return &output, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "GetDeployment",
		Description: "Retrieve an active deployment.",
		Method:      http.MethodGet,
		Path:        "/deployment/{url}",
		Responses: map[string]*huma.Response{
			"200": {
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: deploymentSchema,
					},
				},
			},
		},
	}, func(ctx context.Context, input *struct {
		Url string `path:"url"`
	}) (*GetDeploymentOutput, error) {
		permissions, permissionsOk := ctx.Value("permissions").(Permissions)
		if !permissionsOk {
			return nil, fmt.Errorf("Auth check failed somehow")
		}

		url := urlFromString(input.Url)

		deployment, err := a.web.GetDeploymentByUrl(&url)
		if err != nil || deployment.DontPersist {
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

		model, err := deploymentToApiModel(deployment)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}
		var output GetDeploymentOutput
		output.Body = model
		return &output, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "CreateAlias",
		Method:      http.MethodPut,
		Description: "Create an alias deployment.",
		Path:        "/deploy/alias",
	}, func(ctx context.Context, input *DeployAliasInput) (*SuccessOutput, error) {
		permissions, permissionsOk := ctx.Value("permissions").(Permissions)
		if !permissionsOk {
			return nil, huma.Error500InternalServerError("Auth check failed somehow")
		}

		if !permissions.CanCreateDeployment() {
			return nil, huma.Error401Unauthorized("Not authorized to create deployments")
		}

		err := a.web.PutAliasDeployment(
			urlFromString(input.Body.Url),
			urlFromString(*input.Body.AliasBase.AliasedTo),
			*input.Body.AliasBase.Redirect,
		)

		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		var output SuccessOutput
		output.Body.Success = true
		output.Body.Message = "Created alias deployment"
		return &output, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "DeployFiles",
		Description: "Put files in an existing deployment.",
		Method:      http.MethodPut,
		Path:        "/deploy/files",
	}, func(
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

	huma.Register(api, huma.Operation{
		OperationID: "DeployAdminDash",
		Description: "Deploy the admin dashboard to a specified URL.",
		Method:      http.MethodPut,
		Path:        "/admin-dash",
	}, func(ctx context.Context, input *DeployAdminDashInput) (*SuccessOutput, error) {
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
