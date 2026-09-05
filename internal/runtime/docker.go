package runtime

import (
	"context"

	"vortex-edge/internal/models"

	   "github.com/moby/moby/client"
    // "github.com/moby/moby/api/types"
)

type DockerRuntime struct {
	client *client.Client
}

// NewDockerRuntime creates a Docker SDK client and returns the runtime.
func NewDockerRuntime() (*DockerRuntime, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}

	return &DockerRuntime{
		client: cli,
	}, nil
}

func (d *DockerRuntime) ServerVersion() (string, error) {
    ctx := context.Background()

    version, err := d.client.ServerVersion(ctx,
	    client.ServerVersionOptions{},)
    if err != nil {
        return "", err
    }

    return version.Version, nil
}

func (d *DockerRuntime) Run(service *models.Service) error {
	// TODO: Pull image
	// TODO: Create container
	// TODO: Start container
	// TODO: Inspect container
	return nil
}

func (d *DockerRuntime) Stop(service *models.Service) error {
	// TODO: Stop Docker container
	return nil
}

func (d *DockerRuntime) Delete(service *models.Service) error {
	// TODO: Remove Docker container
	return nil
}

func (d *DockerRuntime) Logs(service *models.Service) (string, error) {
	// TODO: Fetch Docker logs
	return "", nil
}

// Ping verifies that the Docker daemon is reachable.
func (d *DockerRuntime) Ping() error {
	ctx := context.Background()

	_, err := d.client.Ping(ctx,
	client.PingOptions{},)
	return err
}