package orchestrator
import(
	"fmt"	
"vortex-edge/internal/cluster"
	"vortex-edge/internal/agent"
	"vortex-edge/internal/models"
	"vortex-edge/internal/scheduler"
	"vortex-edge/proto"
)

type Orchestrator struct{
	scheduler *scheduler.Scheduler
	    clusterState *cluster.ClusterState
}
func NewOrchestrator(
	scheduler *scheduler.Scheduler,
	clusterState *cluster.ClusterState,
) *Orchestrator {
	return &Orchestrator{
		scheduler:    scheduler,
		clusterState: clusterState,
	}
}

func(o*Orchestrator)Deploy(service *models.Service,)error{
	 node, err:= o.scheduler.Schedule(
		service.CPU,
		service.Memory,

	 )
	 if err!=nil {
		 return err
	 }

	// Reserve resources before starting the container.
	err = o.clusterState.AllocateResources(
		node.ID,
		service.CPU,
		service.Memory,
	)
	if err != nil {
		return err
	}

	service.NodeID = node.ID
	service.ResourcesAllocated = true

	  client, err:= agent.NewClient(node.Address)
	  if err!=nil{
		o.releaseResources(service)
	    return err
	}
	 defer client.Close()
	 response, err := client.RunService(
		&proto.RunServiceRequest{
			ServiceId: service.ID,
			Name:      service.Name,
			Image:     service.Image,
			Cpu:       int32(service.CPU),
			Memory:    service.Memory,
			HostPort:  int32(service.HostPort),
		},
	)
	if err != nil {
		o.releaseResources(service)
		return err
	}

	if !response.Started {
		o.releaseResources(service)
		return fmt.Errorf(
			"service %s failed to start",
			service.ID,
		)
	}

	service.ContainerID = response.ContainerId
	service.Status = models.Running

	return nil
}

func( o*Orchestrator)Delete(service*models.Service)error{
node, err := o.clusterState.NodeByID(service.NodeID)
	 if err!=nil{
		return err
	 }
	 if node== nil {
		 return fmt.Errorf("node %s not found", service.NodeID)
	 }

	    client, err := agent.NewClient(node.Address)
    if err != nil {
        return err
    }
    defer client.Close()
	 response, err:=client.DeleteService(
		&proto.DeleteServiceRequest{
			 ServiceId: service.ID,
			 ContainerId: service.ContainerID,
		},
		
	 )
	 if err != nil {
        return err
    }

    if !response.Deleted {
        return fmt.Errorf("service %s failed to delete", service.ID)
    }

    err = o.releaseResources(service)
    if err != nil {
        return err
    }

    service.ContainerID = ""
    service.Status = models.Stopped

    return nil

}

func (o *Orchestrator) Stop(service *models.Service) error {
	node, err := o.clusterState.NodeByID(service.NodeID)
	if err != nil {
		return err
	}
	if node == nil {
		return fmt.Errorf("node %s not found", service.NodeID)
	}

	client, err := agent.NewClient(node.Address)
	if err != nil {
		return err
	}
	defer client.Close()

	response, err := client.StopService(
		&proto.StopServiceRequest{
			ServiceId:   service.ID,
			ContainerId: service.ContainerID,
		},
	)
	if err != nil {
		return err
	}

	if !response.Stopped {
		return fmt.Errorf("service %s failed to stop", service.ID)
	}

	err = o.releaseResources(service)
	if err != nil {
		return err
	}

	service.Status = models.Stopped

	return nil
}

// releaseResources gives a service's CPU/memory back to its node.
// It only releases once, so calling Stop and then Delete doesn't
// double-count the release.
func (o *Orchestrator) releaseResources(service *models.Service) error {
	if !service.ResourcesAllocated {
		return nil
	}

	err := o.clusterState.ReleaseResources(
		service.NodeID,
		service.CPU,
		service.Memory,
	)
	if err != nil {
		return err
	}

	service.ResourcesAllocated = false

	return nil
}

func (o *Orchestrator) Logs(service *models.Service) (string, error) {
    node, err := o.clusterState.NodeByID(service.NodeID)
    if err != nil {
        return "", err
    }

    if node == nil {
        return "", fmt.Errorf(
            "node %s not found",
            service.NodeID,
        )
    }

    client, err := agent.NewClient(node.Address)
    if err != nil {
        return "", err
    }
    defer client.Close()

    response, err := client.Logs(
        &proto.LogsRequest{
            ServiceId:   service.ID,
            ContainerId: service.ContainerID,
        },
    )
    if err != nil {
        return "", err
    }

    return response.Logs, nil
}