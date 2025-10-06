package content

import (
	"context"
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

func (*ContainerManager) StartContainer(deploymentName string, name string, registry string, port string) (string, error) {

	cli, err := client.NewClientWithOpts()
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	container, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:           registry + "/" + name,
			AttachStdout:    false,
			AttachStderr:    false,
			StopTimeout:     new(int),
			NetworkDisabled: true,
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				nat.Port(port): []nat.PortBinding{
					nat.PortBinding{HostPort: port},
				},
			},
		},
		nil,
		nil,
		deploymentName,
	)

	if err != nil {
		return "", err
	}

	return container.ID, nil
}
