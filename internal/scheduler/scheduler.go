package scheduler
import(
	"errors"
"vortex-edge/internal/models"
	"vortex-edge/internal/cluster"
)

var ErrNoSuitableNode= errors.New("no suitable node found")

type Scheduler struct{
	 	clusterState *cluster.ClusterState
}

func NewScheduler(clusterState *cluster.ClusterState) *Scheduler {
	return &Scheduler{
		clusterState: clusterState,
	}
}


func (s *Scheduler) Schedule(
	cpuRequired int,
	memoryRequired int64,
) (*models.Node, error) {

	nodes := s.clusterState.ReadyNodes()

	for _, node := range nodes {
		availableCPU := node.CPU - node.AllocatedCPU
		availableMemory := node.Memory - node.AllocatedMemory

		if availableCPU < cpuRequired {
			continue
		}

		if availableMemory < memoryRequired {
			continue
		}

		return node, nil
	}

	return nil, ErrNoSuitableNode
}