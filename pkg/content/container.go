package content

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/internet-golf/internet-golf/pkg/db"
)

const DOCKER_HUB_REGISTRY = "docker.io"

type ContainerManager struct{ Settings db.StorageSettings }

// probably want to add future support for a full name like "docker.io/hello-world" being passed
func (*ContainerManager) PullContainer(name string, registry string, authToken string) error {

	cli, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}

	ctx := context.Background()
	reader, err := cli.ImagePull(
		ctx,
		registry+"/"+name,
		image.PullOptions{
			RegistryAuth: authToken,
		},
	)

	if err != nil {
		return err
	}

	defer reader.Close()

	if _, err := io.Copy(os.Stdout, reader); err != nil {
		return err
	}

	return nil
}

func (*ContainerManager) StartContainer(name string) (string, error) {

	cli, err := client.NewClientWithOpts()
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	inspection, err := cli.ImageInspect(ctx, name)
	if err != nil {
		return "", errors.New("failed to inspect image for ports")
	}

	portBindings := nat.PortMap{}
	for port := range inspection.Config.ExposedPorts {
		portBindings[nat.Port(port)] = []nat.PortBinding{nat.PortBinding{}}
	}

	cont, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:        name,
			AttachStdout: false,
			AttachStderr: false,
		},
		&container.HostConfig{PortBindings: portBindings},
		nil,
		nil,
		"",
	)

	if err != nil {
		return "", errors.New("failed to create container from image")
	}

	containerId := cont.ID
	if err := cli.ContainerStart(ctx, containerId, container.StartOptions{}); err != nil {
		return "", errors.New("failed start container " + containerId)
	}

	return containerId, nil
}
