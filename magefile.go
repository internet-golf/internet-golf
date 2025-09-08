//go:build mage
// +build mage

package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/magefile/mage/mg"
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

func BuildServer() error {
	mg.Deps(SetExtension, InstallDeps)
	fmt.Println("Building server...")
	cmd := exec.Command("go", "build", "-o", "golf-server"+extension, "./cmd")
	return cmd.Run()
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

	// https://docs.docker.com/reference/api/engine/sdk/examples/#run-a-container

	ctx := context.Background()
	cont, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:        "openapitools/openapi-generator-cli:latest",
			AttachStdout: false,
			AttachStderr: false,
			Cmd: []string{
				"generate", "-i", "/local/golf-openapi.yaml", "-g", "go", "-o", "/local/client-sdk",
				"--additional-properties=packageName=golfsdk,withGoMod=false",
			},
			// TODO: what is new(int) doing here
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
	openapiCmd := exec.Command("go", "run", "./cmd", "openapi")
	if err := openapiCmd.Run(); err != nil {
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
	cmd := exec.Command("go", "build", "-o", "golf"+extension, "./client-cmd")
	return cmd.Run()
}

func InstallDeps() error {
	fmt.Println("Installing dependencies...")
	cmd := exec.Command("go", "get", ".")
	return cmd.Run()
}

func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("golf.exe")
	os.RemoveAll("golf-server.exe")
	os.RemoveAll("golf")
	os.RemoveAll("golf-server")
}
