//go:build mage
// +build mage

package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/moby/term"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Build

var extension string = ""

func SetExtension() {
	if runtime.GOOS == "windows" {
		extension = ".exe"
	}
}

func Build() {
	mg.Deps(BuildClient, BuildServer)
}

func BuildDash() error {
	fmt.Println("Building dashboard...")
	return sh.Run("docker", "build", "-o", "./pkg/resources/dash-dist/", "./admin-dash/")
}

func BuildServer() error {
	mg.Deps(SetExtension, InstallDeps, BuildDash)
	fmt.Println("Building server...")
	return sh.Run("go", "build", "-o", "golf-server"+extension, "./cmd")
}

func runOpenapiGenerator() error {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		fmt.Println("Unable to create docker client")
		return err
	}

	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return cwdErr
	}

	// https://stackoverflow.com/a/48579861/3962267
	imageName := "openapitools/openapi-generator-cli:latest"
	ctx := context.Background()
	reader, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
	termFd, isTerm := term.GetFdInfo(os.Stderr)
	pullErr := jsonmessage.DisplayJSONMessagesStream(reader, os.Stderr, termFd, isTerm, nil)
	reader.Close()
	if pullErr != nil {
		panic(pullErr)
	}

	// https://docs.docker.com/reference/api/engine/sdk/examples/#run-a-container
	// TODO: not all errors produced by the container are actually propagated
	// back to this script; i think the main process' exit status needs to be
	// checked, or something like that
	cont, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:        imageName,
			AttachStdout: false,
			AttachStderr: false,
			Cmd: []string{
				"generate", "-i", "/local/golf-openapi.yaml", "-g", "go", "-o", "/local/client-sdk",
				"--additional-properties=packageName=golfsdk,withGoMod=false",
				// the generator currently has some weird complaint about the
				// multi-part form data endpoints but all the tests still pass
				// so :P
				"--skip-validate-spec",
			},
			StopTimeout:     new(int),
			NetworkDisabled: true,
		},
		&container.HostConfig{
			Binds: []string{
				cwd + ":/local",
			},
		},
		nil,
		nil,
		"",
	)
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, cont.ID, container.StartOptions{}); err != nil {
		return err
	}

	statusCh, errCh := cli.ContainerWait(ctx, cont.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, cont.ID, container.LogsOptions{ShowStdout: true})
	if err != nil {
		return err
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	if err := cli.ContainerRemove(ctx, cont.ID, container.RemoveOptions{Force: true}); err != nil {
		return err
	}

	return nil
}

func GenerateClientSdk() error {
	mg.Deps(InstallDeps)
	fmt.Println("Regenerating client code...")
	if err := sh.Run("go", "run", "./cmd", "openapi"); err != nil {
		return err
	}
	if openapiErr := runOpenapiGenerator(); openapiErr != nil {
		return openapiErr
	}

	// delete annoying extra files that openapi-generator creates

	if err := os.RemoveAll("client-sdk/test"); err != nil {
		return err
	}
	if err := os.Remove("client-sdk/.travis.yml"); err != nil {
		return err
	}
	if err := os.Remove("client-sdk/git_push.sh"); err != nil {
		return err
	}

	return nil
}

func BuildClient() error {
	mg.Deps(SetExtension, InstallDeps, GenerateClientSdk)
	fmt.Println("Building client...")
	return sh.Run("go", "build", "-o", "golf"+extension, "./client-cmd")
}

func InstallDeps() error {
	fmt.Println("Installing dependencies...")
	return sh.Run("go", "get", ".")
}

func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("golf.exe")
	os.RemoveAll("golf-server.exe")
	os.RemoveAll("golf")
	os.RemoveAll("golf-server")
}
