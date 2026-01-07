package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/internet-golf/internet-golf/pkg/db"
	"github.com/internet-golf/internet-golf/pkg/utils"
)

// this returns a huma middleware function that figures out the permissions
// assigned to the entity making the request and stores them in the request
// context
func readAuth(api huma.API, authManager *AuthManager) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {

		// special case: don't try to determine user permissions if they're just
		// getting a health check
		if ctx.URL().Path == "/alive" {
			next(ctx)
			return
		}

		// this header is set by the internal caddy reverse-proxy when it is
		// forwarding a request - we don't want to mistake those for "true"
		// localhost requests
		remoteAddr := ctx.Header("X-Forwarded-For")
		if len(remoteAddr) == 0 {
			remoteAddr = ctx.RemoteAddr()
		}
		authHeader := ctx.Header("Authorization")

		// TODO: recover from any panics in getPermissionForRequest?
		permissions, error := authManager.GetPermissionsForRequest(remoteAddr, authHeader)
		if error != nil {
			fmt.Fprintf(os.Stderr, "Error getting permissions for request: %s\n", error.Error())
		} else {
			ctx = huma.WithValue(ctx, "permissions", permissions)
		}

		next(ctx)
	}
}

type AddExternalUserBody struct {
	ExternalUserHandle string                `json:"externalUserHandle,omitempty" docs:"A username, like \"internet-golf\" for Github user @internet-golf. Will be ignored if externalUserId is specified."`
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

type AdminApi struct {
	web    *DeploymentBus
	auth   *AuthManager
	config *utils.Config
}

func NewAdminApi(bus *DeploymentBus, db db.Db, config *utils.Config) *AdminApi {
	return &AdminApi{
		web:    bus,
		auth:   NewAuthManager(db),
		config: config,
	}
}

var humaConfig = huma.DefaultConfig("Internet Golf API", "0.5.0")

// this function sets up the endpoints for the server's admin API. note that the
// route handlers are meant to be fairly thin; their job is to parse incoming
// data, verify the permissions of the entity making the request, and respond to
// the client afterward. it's easier to handle business logic in dedicated
// entities like DeploymentBus and AuthManager.
func (a *AdminApi) addRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "HealthCheck",
		Description: "Find out if the server is up or not.",
		Method:      http.MethodGet,
		Path:        "/alive",
	},
		func(ctx context.Context, i *struct{}) (*HealthCheckOutput, error) {
			resp := &HealthCheckOutput{}
			resp.Body.Ok = true
			return resp, nil
		})

	a.addDeploymentRoutes(api)

	// TODO: separate out user/deployment routes, just like deployment routes
	// have their own file and method

	huma.Put(api, "/user/register", func(ctx context.Context, input *AddExternalUserInput) (*SuccessOutput, error) {
		permissions, permissionsOk := ctx.Value("permissions").(Permissions)
		if !permissionsOk {
			return nil, fmt.Errorf("Auth check failed somehow")
		}

		if !permissions.CanCreateCredentials() {
			return nil, huma.Error401Unauthorized("You are not authorized to add a user")
		}

		if len(input.Body.ExternalUserHandle) == 0 && len(input.Body.ExternalUserId) == 0 {
			return nil, huma.Error400BadRequest("Either ID or handle must be specified.")
		}

		// TODO: move this logic into AuthManager somehow

		if len(input.Body.ExternalUserId) == 0 {
			if input.Body.ExternalUserSource == db.Github {
				// example api url: https://api.github.com/users/internet-golf
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

		a.auth.RegisterExternalUser(db.ExternalUser{
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
		token, err := a.auth.CreateBearerToken(input.Body.FullPermissions)
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
	b, _ := api.OpenAPI().DowngradeYAML()
	openApiErr := os.WriteFile(outputPath, b, 0644)
	if openApiErr != nil {
		fmt.Println("Could not write OpenAPI spec")
	}
}

func (a *AdminApi) CreateServer() *http.Server {
	if len(a.config.AdminApiPort) == 0 {
		panic("Admin API port not set")
	}

	router := http.NewServeMux()
	api := humago.New(router, humaConfig)

	api.UseMiddleware(readAuth(api, a.auth))

	a.addRoutes(api)

	fmt.Println("Starting admin API server at http://127.0.0.1:" + a.config.AdminApiPort)
	address := "0.0.0.0"
	if a.config.LocalOnly {
		address = "127.0.0.1"
	}
	server := http.Server{Addr: address + ":" + a.config.AdminApiPort, Handler: router}
	return &server
}
