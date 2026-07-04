package runtime

import (
	"vortex-edge/internal/models"
)

type DockerRuntime struct{}

func NewDockerRuntime() *DockerRuntime {
	return &DockerRuntime{}
}

func (d *DockerRuntime) Run(service *models.Service) error {
 if err := d.ensureImage(service); err != nil {
        return err
    }

    containerID, err := d.createContainer(service)
    if err != nil {
        return err
    }

    if err := d.startContainer(containerID); err != nil {
        return err
    }

    return d.inspectContainer(service, containerID)
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