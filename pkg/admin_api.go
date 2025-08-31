package internetgolf

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/gosimple/slug"
)

type AuthResult struct {
	localRequest       bool
	externalSourceType ExternalSourceType
	externalSource     string
}

func readAuth(api huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		authHeader := strings.Split(ctx.Header("Authorization"), " ")

		var authResult AuthResult

		if ctx.RemoteAddr() == "127.0.0.1" || strings.HasPrefix(ctx.RemoteAddr(), "127.0.0.1:") {
			authResult.localRequest = true
		} else {
			authResult.localRequest = false
		}
		if len(authHeader) == 2 && authHeader[0] == "Github-OIDC" {
			jwtData, err := ParseGithubOidcToken(authHeader[1])
			if err != nil {
				huma.WriteErr(api, ctx, http.StatusInternalServerError, "Failed to parse Github OIDC token")
				return
			}
			authResult.externalSourceType = GithubRepo
			repo := jwtData.Repository
			if jwtData.RefType == "branch" {
				branchNameLocation := strings.LastIndex(jwtData.Ref, "/")
				if branchNameLocation != -1 {
					repo += jwtData.Ref[branchNameLocation+1:]
				}
			}
			authResult.externalSource = repo
		}
		ctx = huma.WithValue(ctx, "authResult", authResult)
		next(ctx)
	}
}

type DeploymentCreateInput struct {
	Body struct {
		DeploymentMetadata
	}
}

type DeployFilesInput struct {
	RawBody huma.MultipartFormFiles[struct {
		// this can be specified manually or it can come from the token in the
		// auth header (which is scoped to specific deployments)
		Url                    string        `form:"url" required:"false"`
		Contents               huma.FormFile `form:"contents" contentType:"application/gzip,application/octet-stream"`
		KeepLeadingDirectories bool          `form:"keepLeadingDirectories"`
		PreserveExistingFiles  bool          `form:"preserveExistingFiles"`
	}]
}

type DeployContainerInput struct {
	Body struct {
		// should this be json even though the static deployment input is form data??
		ContainerUrl    string `json:"containerUrl"`
		InternalAppPort int    `json:"internalAppPort"`
		Url             string `json:"url" required:"false"`
	}
}

type SuccessOutput struct {
	Body struct {
		Success bool `json:"success"`
	}
}

type HealthCheckOutput struct {
	Body struct {
		Ok bool `json:"ok"`
	}
}

type AdminApi struct {
	Web             DeploymentBus
	StorageSettings StorageSettings
	Port            string
}

func (a *AdminApi) Start() {
	if len(a.Port) == 0 {
		panic("Admin API port not set")
	}

	router := http.NewServeMux()
	api := humago.New(
		router,
		huma.DefaultConfig("Internet Golf API", "0.5.0"),
	)

	api.UseMiddleware(readAuth(api))

	huma.Get(api, "/alive",
		func(ctx context.Context, i *struct{}) (*HealthCheckOutput, error) {
			resp := &HealthCheckOutput{}
			resp.Body.Ok = true
			return resp, nil
		})

	huma.Post(api, "/deploy/new", func(
		ctx context.Context, input *DeploymentCreateInput,
	) (*SuccessOutput, error) {
		putDeploymentErr := a.Web.SetupDeployment(input.Body.DeploymentMetadata)
		if putDeploymentErr != nil {
			return nil, putDeploymentErr
		}
		var output SuccessOutput
		output.Body.Success = true
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

			// 1. based on the authorization information set by the readAuth
			// middleware, figure out how the content should be sent to the
			// deployment bus once it's created (or return an error if it
			// shouldn't be sent anywhere.) TODO: extract this logic so it can
			// be used in other /deploy api routes somehow

			authResult, authResultOk := ctx.Value("authResult").(AuthResult)
			if !authResultOk {
				return nil, fmt.Errorf("Auth check failed somehow")
			}
			var deployContent func(DeploymentContent) error
			if len(authResult.externalSourceType) > 0 {
				if !a.Web.DeploymentWithExternalSourceExists(
					authResult.externalSource, authResult.externalSourceType,
				) {
					return nil, fmt.Errorf(
						"Could not find deployment with external source %s (%s)",
						authResult.externalSource, authResult.externalSourceType,
					)
				}
				deployContent = func(content DeploymentContent) error {
					return a.Web.PutDeploymentContentByExternalSource(
						authResult.externalSource, authResult.externalSourceType, content,
					)
				}
			} else if len(formData.Url) > 0 {
				// TODO: other checks using keys or whatever
				if !authResult.localRequest {
					return nil, fmt.Errorf("not authorized to deploy to %s", formData.Url)
				}
				deployContent = func(content DeploymentContent) error {
					return a.Web.PutDeploymentContentByUrl(
						formData.Url,
						content,
					)
				}
			} else {
				return nil, fmt.Errorf("no url or external source; unclear where to deploy to")
			}

			// 2. actually create the deployment content locally

			hash, hashErr := hashStream(formData.Contents)
			if hashErr != nil {
				return nil, fmt.Errorf("could not hash files for %v", formData.Url)
			}
			outDir := path.Join(
				a.StorageSettings.DataDirectory,
				slug.Make(formData.Url),
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

			err := deployContent(DeploymentContent{
				ServedThingType: StaticFiles,
				ServedThing:     outDir,
			})
			if err != nil {
				return nil, err
			}

			// TODO: delete the old directory after deployContent is
			// finished? presumably that'll be safe

			// 4. return success

			output := SuccessOutput{}
			output.Body.Success = true
			return &output, nil
		})

	// huma.Put(api, "/deploy/container", func(
	// 	ctx context.Context, input *DeployContainerInput,
	// ) (*SuccessOutput, error) {
	// 	// TODO: implement this at all
	// 	panic("not implemented")
	// })

	fmt.Println("Writing OpenAPI spec to golf-openapi.yaml")
	b, _ := api.OpenAPI().DowngradeYAML()
	openApiErr := os.WriteFile("golf-openapi.yaml", b, 0644)
	if openApiErr != nil {
		fmt.Println("Could not write OpenAPI spec")
	}

	fmt.Println("Starting admin API server at http://127.0.0.1:" + a.Port)
	// TODO: bind to more addresses? i guess not bc this is exposed via a caddy
	// reverse proxy
	http.ListenAndServe("127.0.0.1:"+a.Port, router)
}
