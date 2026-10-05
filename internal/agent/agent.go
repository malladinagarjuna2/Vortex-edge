package agent 
import (
	"context"
	"vortex-edge/internal/models"
	"vortex-edge/internal/runtime"
	"vortex-edge/proto"
	"time"
)

type Agent struct {
	proto.UnimplementedNodeAgentServer
	 node *models.Node
	 runtime runtime.Runtime
}

func NewAgent(node *models.Node, runtime runtime.Runtime)*Agent{
	  return &Agent{
		node: node, 
		runtime: runtime,
	  }
}
func(a*Agent)Node() *models.Node{
	 return a.node
}

func (a*Agent)Heartbeat(
	ctx context.Context,	
	req *proto.HeartbeatRequest,

)(*proto.HeartbeatResponse,error){
	return &proto.HeartbeatResponse{
		Acknowledged: true,
	}, nil
}

func (a *Agent) DeleteService(
    ctx context.Context,
    req *proto.DeleteServiceRequest,
) (*proto.DeleteServiceResponse, error) {
// 	Orchestrator
//      ↓
// DeleteService request
//      ↓
// Node Agent
//      ↓
// runtime.Delete()
//      ↓
// DockerRuntime
//      ↓
// ContainerRemove()
//      ↓
// Docker container deleted

    service := &models.Service{
        ID:          req.ServiceId,
        ContainerID: req.ContainerId,
        Status:      models.Stopped,
    }

    err := a.runtime.Delete(service)
    if err != nil {
        return &proto.DeleteServiceResponse{
            Deleted: false,
        }, err
    }

    return &proto.DeleteServiceResponse{
        Deleted: true,
    }, nil
}
func (a *Agent) StartHeartbeat(client *Client, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C

		acknowledged, err := client.Heartbeat(a.node.ID)
		if err != nil {
			a.node.Status = models.NodeUnhealthy
			continue
		}

		if acknowledged {
			a.node.Status = models.NodeReady
			a.node.LastHeartbeat = time.Now()
		}
	}
}

func(a *Agent)RunService(ctx context.Context, req*proto.RunServiceRequest,)(*proto.RunServiceResponse, error){
		service := &models.Service{
		ID:        req.ServiceId,
		Name:      req.Name,
		Image:     req.Image,
		CPU:       int(req.Cpu),
		Memory:    req.Memory,
		HostPort:  int(req.HostPort),
		Status:    models.Deploying,
	}

	err := a.runtime.Run(service)
	if err != nil {
		return &proto.RunServiceResponse{
			Started: false,
		}, err
	}

	return &proto.RunServiceResponse{
		Started: true,
		   ContainerId: service.ContainerID,
	}, nil
}


// gRPC request
//       ↓
// models.Service
//       ↓
// runtime.Run(service)