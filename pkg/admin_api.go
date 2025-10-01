package internetgolf

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/gosimple/slug"
	"github.com/toBeOfUse/internet-golf/pkg/auth"
	"github.com/toBeOfUse/internet-golf/pkg/db"
	"github.com/toBeOfUse/internet-golf/pkg/deployment"
	"github.com/toBeOfUse/internet-golf/pkg/utils"
)

func readAuth(api huma.API, authManager auth.AuthManager) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		remoteAddr := ctx.RemoteAddr()
		authHeader := ctx.Header("Authorization")

		// TODO: recover from any panics in getPermissionForRequest?
		permissions, _ := authManager.GetPermissionsForRequest(remoteAddr, authHeader)
		ctx = huma.WithValue(ctx, "permissions", permissions)

		next(ctx)
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
	Body struct {
		DeploymentCreateBody
	}
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

type AddExternalUserBody struct {
	ExternalUserHandle string                `json:"externalUserHandle,omitempty" docs:"A username, like \"toBeOfUse\" for Github user @toBeOfUse. Will be ignored if externalUserId is specified."`
	ExternalUserId     string                `json:"externalUserId,omitempty" docs:"The ID that the user has in the external system. Not needed if externalUserHandle is specified."`
	ExternalUserSource db.ExternalSourceType `json:"externalUserSource" docs:"The location of the external user. Currently only supports \"Github\"."`
}
type AddExternalUserInput struct {
	Body struct {
		AddExternalUserBody
	}
}

type CreateBearerTokenBody struct {
	FullPermissions bool `json:"fullPermissions"`
	// TODO: more granular permissions
}

type CreateBearerTokenInput struct {
	Body struct {
		CreateBearerTokenBody
	}
}

type CreateBearerTokenOutput struct {
	Body struct {
		Token string `json:"token"`
	}
}

type DeployContainerInput struct {
	Body struct {
		// should this be json even though the static deployment input is form data??
		ContainerUrl    string `json:"containerUrl"`
		InternalAppPort int    `json:"internalAppPort"`
		Id              string `json:"id" required:"false"`
	}
}

type SuccessOutput struct {
	Body struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
}

type HealthCheckOutput struct {
	Body struct {
		Ok bool `json:"ok"`
	}
}

type GetDeploymentOutput struct {
	Body struct {
		db.Deployment
	}
}

type AdminApi struct {
	Web  *deployment.DeploymentBus
	Auth auth.AuthManager
	Port string
}

var humaConfig = huma.DefaultConfig("Internet Golf API", "0.5.0")

func (a *AdminApi) addRoutes(api huma.API) {
	huma.Get(api, "/alive",
		func(ctx context.Context, i *struct{}) (*HealthCheckOutput, error) {
			resp := &HealthCheckOutput{}
			resp.Body.Ok = true
			return resp, nil
		})

	huma.Post(api, "/deploy/new", func(
		ctx context.Context, input *DeploymentCreateInput,
	) (*SuccessOutput, error) {
		permissions, permissionsOk := ctx.Value("permissions").(auth.Permissions)
		if !permissionsOk {
			return nil, huma.Error500InternalServerError("Auth check failed somehow")
		}

		if !permissions.CanCreateDeployment() {
			return nil, huma.Error401Unauthorized("Not authorized to create deployments")
		}

		input.Body.DeploymentMetadata.Url = deployment.UrlFromString(input.Body.Url)

		// TODO: validate externalSourceType and i guess Domain and Path
		putDeploymentErr := a.Web.SetupDeployment(input.Body.DeploymentMetadata)
		if putDeploymentErr != nil {
			return nil, putDeploymentErr
		}
		var output SuccessOutput
		output.Body.Success = true
		output.Body.Message = fmt.Sprintf("Created deployment with url %s", input.Body.Url)
		return &output, nil
	})

	huma.Get(api, "/deployment/{url}", func(ctx context.Context, input *struct {
		Url string `path:"url"`
	}) (*GetDeploymentOutput, error) {
		permissions, permissionsOk := ctx.Value("permissions").(auth.Permissions)
		if !permissionsOk {
			return nil, fmt.Errorf("Auth check failed somehow")
		}

		url := deployment.UrlFromString(input.Url)

		deployment, err := a.Web.GetDeploymentByUrl(&url)
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
			// 0. parse the form data
			formData := input.RawBody.Data()
			fmt.Printf("received form data: %+v\n", formData)

			permissions, permissionsOk := ctx.Value("permissions").(auth.Permissions)
			if !permissionsOk {
				return nil, fmt.Errorf("Auth check failed somehow")
			}

			url := deployment.UrlFromString(formData.Url)
			deployment, findDeploymentError := a.Web.GetDeploymentByUrl(&url)
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

			// 2. actually create the deployment content locally

			hash, hashErr := utils.HashStream(formData.Contents)
			if hashErr != nil {
				return nil, fmt.Errorf("could not hash files for %s", url)
			}
			outDir := path.Join(
				// this is kind of a hack and file storage should probably be
				// encapsulated in its own interface instead of being handled by
				// random utility functions here in the admin api route handler
				a.Auth.Db.GetStorageDirectory(),
				slug.Make(formData.Url),
				hash,
			)
			// weirdly, formData.Contents is a seekable stream, which i'm pretty
			// sure means its entire contents must be being kept in memory so that
			// they can be sought back to (unless it falls back to saving them
			// to disk for large files?) this seems like an annoying limitation
			if tarGzError := utils.ExtractTarGz(
				formData.Contents, outDir, !formData.KeepLeadingDirectories,
			); tarGzError != nil {
				return nil, tarGzError
			}

			// TODO: if PreserveExistingFiles, find the existing deployment for
			// this url and copy its files into the new directory (if they're
			// not the same directory (i.e. if the hashes are unequal) (which
			// will presumably require getting a reference to the existing
			// deployment from the bus - maybe instead of creating deployContent
			// we should be getting a cursor/active record))

			// 3. send the content to the deployment bus using the function that
			// was created in step 1

			err := a.Web.PutDeploymentContentByUrl(
				url,
				db.DeploymentContent{
					ServedThingType: db.StaticFiles,
					ServedThing:     outDir,
				},
			)
			if err != nil {
				return nil, err
			}
			// TODO: delete the old directory after deployContent is
			// finished? presumably that'll be safe

			// 4. return success

			output := SuccessOutput{}
			output.Body.Success = true
			output.Body.Message = "Updated content for " + url.String()
			return &output, nil
		})

	huma.Put(api, "/user/register", func(ctx context.Context, input *AddExternalUserInput) (*SuccessOutput, error) {
		permissions, permissionsOk := ctx.Value("permissions").(auth.Permissions)
		if !permissionsOk {
			return nil, fmt.Errorf("Auth check failed somehow")
		}

		if !permissions.CanCreateCredentials() {
			return nil, huma.Error401Unauthorized("You are not authorized to add a user")
		}

		if len(input.Body.ExternalUserHandle) == 0 && len(input.Body.ExternalUserId) == 0 {
			return nil, huma.Error400BadRequest("Either ID or handle must be specified.")
		}

		if len(input.Body.ExternalUserId) == 0 {
			if input.Body.ExternalUserSource == db.Github {
				// example api url: https://api.github.com/users/toBeOfUse
				resp, err := http.Get(
					"https://api.github.com/users/" + strings.TrimLeft(input.Body.ExternalUserHandle, "@"),
				)
				if err != nil || resp.StatusCode != 200 {
					return nil, huma.Error400BadRequest("Could not find user")
				}
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, huma.Error500InternalServerError("Got unusable response from the Github API")
				}
				var apiObj struct {
					Id int64 `json:"id"`
				}
				err = json.Unmarshal(body, &apiObj)
				if err != nil || apiObj.Id == 0 {
					return nil, huma.Error500InternalServerError("Could not parse JSON from the Github API")
				}
				input.Body.ExternalUserId = strconv.FormatInt(apiObj.Id, 10)
			} else {
				return nil, huma.Error400BadRequest("External source of user not recognized")
			}
		}

		a.Auth.RegisterExternalUser(db.ExternalUser{
			ExternalSource: input.Body.ExternalUserSource,
			ExternalId:     input.Body.ExternalUserId,
			// defaulting to full permissions until more granular permissions are added
			FullPermissions: true,
		})

		var output SuccessOutput
		output.Body.Success = true
		output.Body.Message = fmt.Sprintf(
			"Successfully added %s user with ID %s", input.Body.ExternalUserSource, input.Body.ExternalUserId,
		)
		return &output, nil
	})

	// TODO: get (all?) users endpoint

	huma.Post(api, "/token/generate", func(ctx context.Context, input *CreateBearerTokenInput) (*CreateBearerTokenOutput, error) {
		token, err := a.Auth.CreateBearerToken(input.Body.FullPermissions)
		if err != nil {
			return nil, huma.Error500InternalServerError("Could not generate token: " + err.Error())
		}

		var output CreateBearerTokenOutput
		output.Body.Token = token
		return &output, nil
	})

	// huma.Put(api, "/deploy/container", func(
	// 	ctx context.Context, input *DeployContainerInput,
	// ) (*SuccessOutput, error) {
	// 	// TODO: implement this at all
	// 	panic("not implemented")
	// })
}

func (a *AdminApi) OutputOpenApiSpec(outputPath string) {
	router := http.NewServeMux()
	api := humago.New(router, humaConfig)

	a.addRoutes(api)

	fmt.Printf("Writing OpenAPI spec to %s\n", outputPath)
	b, _ := api.OpenAPI().YAML()
	openApiErr := os.WriteFile(outputPath, b, 0644)
	if openApiErr != nil {
		fmt.Println("Could not write OpenAPI spec")
	}
}

func (a *AdminApi) CreateServer() *http.Server {
	if len(a.Port) == 0 {
		panic("Admin API port not set")
	}

	router := http.NewServeMux()
	api := humago.New(router, humaConfig)

	api.UseMiddleware(readAuth(api, a.Auth))

	a.addRoutes(api)

	fmt.Println("Starting admin API server at http://127.0.0.1:" + a.Port)
	// TODO: bind to more addresses? i guess not bc this is exposed via a caddy
	// reverse proxy
	server := http.Server{Addr: "127.0.0.1:" + a.Port, Handler: router}
	return &server
}
