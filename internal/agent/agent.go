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