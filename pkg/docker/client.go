package docker

import (
	"github.com/burmilla/os/config"

	dockerClient "github.com/docker/engine-api/client"
	"golang.org/x/net/context"
)

func NewSystemClient() (dockerClient.APIClient, error) {
	return NewClient(config.SystemContainerdHost)
}

func NewDefaultClient() (dockerClient.APIClient, error) {
	return NewClient(config.DockerHost)
}

func NewClient(endpoint string) (dockerClient.APIClient, error) {
	client, err := dockerClient.NewClient(endpoint, "", nil, nil)
	if err != nil {
		return nil, err
	}

	err = ClientOK(endpoint, func() bool {
		_, err := client.Info(context.Background())
		return err == nil
	})

	return client, err
}
