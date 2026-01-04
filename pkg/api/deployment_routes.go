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

type DeploymentBase struct {
	Url string `json:"url" doc:"URL that this deployment will appear at. The DNS for the domain has to be set up first." example:"mysite.mydomain.com"`

	// assuming that there won't be multiple external sources...
	ExternalSource     string `json:"externalSource,omitempty" required:"false" doc:"Original repository for this deployment's source. Can include a branch name." example:"user/repo or user/repo#branch-name"`
	ExternalSourceType string `json:"externalSourceType,omitempty" required:"false" doc:"Place where the original repository lives."`

	Tags []string `json:"tags" required:"false" doc:"Tags used for metadata."`

	PreserveExternalPath bool `json:"preserveExternalPath" required:"false" doc:"if this is true and the deployment url has a path like \"/thing\", then the \"/thing\" in the path will be transparently passed through to the underlying resource instead of being removed (which is the default)"`

	Type string `json:"type" required:"false" doc:"Type of deployment contents; can be StaticSite, Alias, or Empty."`
}
type DeploymentCreateInput struct {
	Body struct{ DeploymentBase }
}

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

type StaticSiteResponseBase struct {
	// these values are pointers so that they will be properly omitted from the
	// JSON response if not set by the API handler (which will happen when
	// creating a DeploymentBody for a non-static-site deployment)
	ServerContentLocation *string `json:"serverContentLocation,omitempty"`
	SpaMode               *bool   `json:"spaMode,omitempty"`
}
type StaticSiteDeployment struct {
	DeploymentBase
	Type string `json:"type" enum:"StaticSite,Alias"`
	StaticSiteResponseBase
}

type AliasBase struct {
	// these values are pointers so that they will be properly omitted from the
	// JSON response if not set by the API handler (which will happen when
	// creating a DeploymentBody for a non-alias deployment)
	AliasedTo *string `json:"aliasedTo,omitempty" doc:"The URL that this deployment is an alias for."`
	Redirect  *bool   `json:"redirect,omitempty" doc:"If this is true, visitors to this deployment's URL will be completely redirected to the URL that this alias is for."`
}
type AliasDeployment struct {
	DeploymentBase
	AliasBase
}

// this is the type that is returned for a deployment from the api. it combines
// the properties of all the possible deployment types. to let api consumers use
// this more safely, a custom schema (`deploymentSchema`) is created for it and
// inserted into the OpenAPI specification. the custom schema narrows which
// fields are available on each deployment returned by the api based on the
// "type" field (which is in DeploymentBase.)
type DeploymentBody struct {
	DeploymentBase
	StaticSiteResponseBase
	AliasBase
}
type GetDeploymentOutput struct {
	Body DeploymentBody
}

type GetDeploymentsOutput struct {
	Body struct {
		Deployments []DeploymentBody `json:"deployments" required:"true"`
	}
}

func deploymentToApiModel(deployment db.Deployment) (DeploymentBody, error) {
	var output DeploymentBody
	output.DeploymentBase = DeploymentBase{
		Url:                  deployment.Url.String(),
		ExternalSource:       deployment.ExternalSource,
		ExternalSourceType:   string(deployment.ExternalSourceType),
		Tags:                 deployment.Tags,
		PreserveExternalPath: deployment.PreserveExternalPath,
	}
	if deployment.ServedThingType == db.StaticFiles {
		output.Type = "StaticSite"
		output.StaticSiteResponseBase.ServerContentLocation = &deployment.ServedThing
		output.StaticSiteResponseBase.SpaMode = &deployment.SpaMode
	} else if deployment.ServedThingType == db.Alias {
		output.Type = "Alias"
		aliasedTo := deployment.AliasedTo.String()
		output.AliasBase.AliasedTo = &aliasedTo
		output.AliasBase.Redirect = &deployment.Redirect
	} else if len(deployment.ServedThingType) == 0 {
		output.Type = "Empty"
	} else {
		return DeploymentBody{}, fmt.Errorf("Deployment type " + string(deployment.ServedThingType) + " not supported by the API")
	}
	return output, nil
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
	deploymentSchema := &huma.Schema{
		OneOf: []*huma.Schema{
			registry.Schema(reflect.TypeOf(StaticSiteDeployment{}), true, ""),
			registry.Schema(reflect.TypeOf(AliasDeployment{}), true, ""),
			registry.Schema(reflect.TypeOf(DeploymentBase{}), true, ""),
		},
		Discriminator: &huma.Discriminator{
			PropertyName: "type", Mapping: map[string]string{
				"StaticSite": "#/components/schemas/StaticSiteDeployment",
				"Alias":      "#/components/schemas/AliasDeployment",
				"Empty":      "#/components/schemas/DeploymentBase",
			},
		},
	}
	deploymentsSchema := &huma.Schema{
		Properties: map[string]*huma.Schema{
			"deployments": &huma.Schema{
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
		output.Body.Deployments = []DeploymentBody{}
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
		OperationID: "PutAlias",
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

		err := a.web.PutDeploymentContentByUrl(urlFromString(input.Body.Url), db.DeploymentContent{
			HasContent:      true,
			ServedThingType: db.Alias,
			AliasedTo:       urlFromString(*input.Body.AliasedTo),
			Redirect:        *input.Body.Redirect,
		})

		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		var output SuccessOutput
		output.Body.Success = true
		output.Body.Message = "Created alias deployment"
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
