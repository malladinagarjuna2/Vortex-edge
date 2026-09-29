            //  Deploy Service
            //        │
            //        ▼
            //  ┌───────────┐
            //  │ Scheduler │
            //  └─────┬─────┘
            //        │
            //  node-02 selected
            //        │
            //        ▼
            //  gRPC Client
            //        │
            //        ▼
            //  Node Agent
            //        │
            //        ▼
            //  DockerRuntime
            //        │
            //        ▼
            //  Docker Container
package runtime

import (
	"io"
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
  // abhi hum yaha pe sirf  yeh bata rahe hai kee docker client kaise yaha pe ek service ko accept krke uskee image pull krlega
   ctx:=context.Background()
   reader,err := d.client.ImagePull(
	ctx,
	service.Image, 
	client.ImagePullOptions{},
   )
    if err!= nil{
		 return err
	}
	defer reader.Close()
	 _, err=io.Copy(io.Discard, reader)
	 if err!= nil {
		 return err
	 }
     // create the container 
	 	response, err := d.client.ContainerCreate(
		ctx,
		client.ContainerCreateOptions{
			Image: service.Image,
			Name:  service.Name,
		},
	)
	if err != nil {
		return err
	}

// start the container 
_,err = d.client.ContainerStart(
	ctx, 
	response.ID,
	client.ContainerStartOptions{}, 
)
if err!= nil {
	 return err
}
	  service.ContainerID= response.ID
	  service.Status= models.Running


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