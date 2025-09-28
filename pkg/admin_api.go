package internetgolf

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/gosimple/slug"
)

func readAuth(api huma.API, authManager AuthManager) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {

		remoteAddr := ctx.RemoteAddr()
		authHeader := ctx.Header("Authorization")

		permissions, _ := authManager.getPermissionsForRequest(remoteAddr, authHeader)
		ctx = huma.WithValue(ctx, "permissions", permissions)

		next(ctx)
	}
}

type DeploymentCreateInput struct {
	Body struct {
		DeploymentMetadata
	}
}

type DeployFilesBody struct {
	// Name is only required when the Authorization header does not imply a
	// specific, single deployment with a token that is scoped to an
	// externalSource and externalSourceType
	Name                   string        `form:"name" required:"false"`
	Contents               huma.FormFile `form:"contents" contentType:"application/gzip,application/octet-stream"`
	KeepLeadingDirectories bool          `form:"keepLeadingDirectories"`
	PreserveExistingFiles  bool          `form:"preserveExistingFiles"`
}
type DeployFilesInput struct {
	RawBody huma.MultipartFormFiles[DeployFilesBody]
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
		Deployment
	}
}

type AdminApi struct {
	Web  DeploymentBus
	Auth AuthManager
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
		if len(input.Body.Name) == 0 {
			if len(input.Body.Urls) == 1 {
				input.Body.Name = (input.Body.Urls[0].Domain + input.Body.Urls[0].Path)
			} else {
				return nil, huma.Error400BadRequest(
					"Specifying a name is required unless there is exactly one URL to use as a name",
				)
			}
		}
		// TODO: validate externalSourceType and i guess Domain and Path
		putDeploymentErr := a.Web.SetupDeployment(input.Body.DeploymentMetadata)
		if putDeploymentErr != nil {
			return nil, putDeploymentErr
		}
		var output SuccessOutput
		output.Body.Success = true
		output.Body.Message = "Created deployment with name " + input.Body.Name
		return &output, nil
	})

	huma.Get(api, "/deployment/{name}", func(ctx context.Context, input *struct {
		Name string `path:"name"`
	}) (*GetDeploymentOutput, error) {
		permissions, permissionsOk := ctx.Value("permissions").(Permissions)
		if !permissionsOk {
			return nil, fmt.Errorf("Auth check failed somehow")
		}

		deployment, err := a.Web.GetDeploymentByName(input.Name)
		if err != nil {
			return nil, huma.Error404NotFound(
				fmt.Sprintf("Could not find deployment called \"%s\"", input.Name),
			)
		}

		if !permissions.canViewDeployment(&deployment) {
			return nil, huma.Error403Forbidden(
				fmt.Sprintf(
					"You are not authorized to view the deployment \"%s\"",
					input.Name,
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

			permissions, permissionsOk := ctx.Value("permissions").(Permissions)
			if !permissionsOk {
				return nil, fmt.Errorf("Auth check failed somehow")
			}

			deployment, findDeploymentError := a.Web.GetDeploymentByName(formData.Name)
			if findDeploymentError != nil {
				return nil, huma.Error404NotFound(
					fmt.Sprintf("could not find deployment with name \"%s\"", formData.Name),
				)
			}

			if !permissions.canModifyDeployment(&deployment) {
				return nil, huma.Error403Forbidden(
					fmt.Sprintf("insufficient permissions to modify deployment \"%s\"", formData.Name),
				)
			}

			// 2. actually create the deployment content locally

			hash, hashErr := hashStream(formData.Contents)
			if hashErr != nil {
				return nil, fmt.Errorf("could not hash files for %v", formData.Name)
			}
			outDir := path.Join(
				// this is kind of a hack and file storage should probably be
				// encapsulated in its own interface instead of being handled by
				// random utility functions here in the admin api route handler
				a.Auth.Db.GetStorageDirectory(),
				slug.Make(formData.Name),
				hash,
			)
			// weirdly, formData.Contents is a seekable stream, which i'm pretty
			// sure means its entire contents must be being kept in memory so that
			// they can be sought back to (unless it falls back to saving them
			// to disk for large files?) this seems like an annoying limitation
			if tarGzError := extractTarGz(
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

			err := a.Web.PutDeploymentContentByName(
				formData.Name,
				DeploymentContent{
					ServedThingType: StaticFiles,
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
			output.Body.Message = "Updated content for the deployment called \"" + formData.Name + "\""
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
