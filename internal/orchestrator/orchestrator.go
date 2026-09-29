package orchestrator
import(
	"fmt"	
	"vortex-edge/internal/agent"
	"vortex-edge/internal/models"
	"vortex-edge/internal/scheduler"
	"vortex-edge/proto"
)

type Orchestrator struct{
	scheduler *scheduler.Scheduler
}
func NewOrchestrator(scheduler *scheduler.Scheduler,)*Orchestrator{
		return &Orchestrator{
		scheduler: scheduler,
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
	  service.NodeID= node.ID
	  client, err:= agent.NewClient(node.Address)
	  if err!=nil{
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
		return err
	}

	if !response.Started {
		return fmt.Errorf(
			"service %s failed to start",
			service.ID,
		)
	}

	service.ContainerID = response.ContainerId
	service.Status = models.Running

	return nil
}