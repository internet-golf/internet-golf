package deppy

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func checkHMAC(ctx huma.Context, next func(huma.Context)) {
	// TODO: check HMAC. i guess the hmac should be in a header and should be
	// based on a hash of the request body's bytes plus a (deployment-specific?) password
	next(ctx)
}

type ContainerDeploymentInput struct {
	Body struct {
		ContainerUrl string `json:"containerUrl"`
	}
}

type ContainerDeploymentOutput struct {
	Body struct {
		Thing string `json:"thing"`
	}
}

type StaticDeploymentInput struct {
	RawBody huma.MultipartFormFiles[struct {
		Id       string        `form:"id"`
		Matcher  string        `form:"matcher"`
		Contents huma.FormFile `form:"contents" contentType:"application/gzip"`
		// ...settings
	}]
}

type StaticDeploymentOutput struct {
	Body struct {
		Id string `json:"id"`
	}
}

type AdminApi struct {
	Web      DeploymentBus
	Settings struct {
		// will be automatically set to $HOME/.deppy if not set
		DataDirectory string
	}
}

func (a AdminApi) Start() {
	// still not sure if chi is better than the current standard library router
	router := chi.NewMux()
	api := humachi.New(
		router,
		huma.DefaultConfig("Deployment Agent API", "0.5.0"),
	)

	var dataDirectoryError error
	a.Settings.DataDirectory, dataDirectoryError = getDataDirectory(a.Settings.DataDirectory)
	if dataDirectoryError != nil {
		panic("Could not create data directory: " + dataDirectoryError.Error())
	}

	api.UseMiddleware(checkHMAC)

	huma.Put(api, "/deploy/container", func(
		ctx context.Context, input *ContainerDeploymentInput,
	) (*ContainerDeploymentOutput, error) {
		fmt.Println(input.Body.ContainerUrl)
		a.Web.PutDeployment(
			Deployment{
				Id:          "identifier",
				Matcher:     "mitch.website/thing",
				ResourceUri: "file://whatever",
			})
		resp := &ContainerDeploymentOutput{}
		resp.Body.Thing = "hi"
		return resp, nil
	})

	huma.Put(
		api,
		"/deploy/files",
		func(
			ctx context.Context, input *StaticDeploymentInput,
		) (*StaticDeploymentOutput, error) {
			output := StaticDeploymentOutput{}

			formData := input.RawBody.Data()

			outDir := path.Join(a.Settings.DataDirectory, formData.Id)
			if tarGzError := extractTarGz(formData.Contents, outDir); tarGzError != nil {
				return nil, tarGzError
			}

			a.Web.PutDeployment(Deployment{
				Id:          formData.Id,
				Matcher:     formData.Matcher,
				ResourceUri: "file://" + outDir,
			})

			output.Body.Id = formData.Id

			return &output, nil
		})

	fmt.Println("Starting admin API server at http://127.0.0.1:8888")
	http.ListenAndServe("127.0.0.1:8888", router)
}
