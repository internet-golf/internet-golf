package deppy

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/gosimple/slug"
)

func checkHMAC(ctx huma.Context, next func(huma.Context)) {
	// TODO: check HMAC. i guess the hmac should be in a header and should be
	// based on a hash of the request body's bytes plus a (deployment-specific?) password
	next(ctx)
}

type ContainerDeploymentInput struct {
	Body struct {
		// should this be json even though the static deployment input is form data??
		ContainerUrl    string `json:"containerUrl"`
		InternalAppPort int    `json:"internalAppPort"`
		PublicUrl       string `json:"publicUrl"`
	}
}

type ContainerDeploymentOutput struct {
	Body struct {
		// no real idea what to return from the api call for any of these
		Thing string `json:"thing"`
	}
}

type StaticDeploymentInput struct {
	RawBody huma.MultipartFormFiles[struct {
		PublicUrl              string        `form:"publicUrl"`
		Contents               huma.FormFile `form:"contents" contentType:"application/gzip,application/octet-stream"`
		KeepLeadingDirectories bool          `form:"keepLeadingDirectories"`
		// TODO: PreserveExistingFiles
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
				ResourceUri: "docker://thing:" + strconv.Itoa((input.Body.InternalAppPort)),
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

			if len(formData.PublicUrl) < 1 {
				panic("public URL for deployment is required")
			}

			fmt.Printf("received form data: %+v\n", formData)

			outDir := path.Join(
				a.Settings.DataDirectory, slug.Make(formData.PublicUrl),
				// TODO: hash of formData.Contents
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
			// not the same directory (i.e. the hashes are unequal))

			a.Web.PutDeployment(Deployment{
				Id:          formData.PublicUrl,
				Matcher:     formData.PublicUrl,
				ResourceUri: "file://" + outDir,
			})

			// delete the old directory after PutDeployment is finished?
			// presumably that'll be safe

			output.Body.Id = "whatever"

			return &output, nil
		})

	fmt.Println("Starting admin API server at http://127.0.0.1:8888")
	http.ListenAndServe("127.0.0.1:8888", router)
}
