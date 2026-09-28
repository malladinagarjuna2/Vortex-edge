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
		if node.CPU < cpuRequired {
			continue
		}

		if node.Memory < memoryRequired {
			continue
		}

		return node, nil
	}

	return nil, ErrNoSuitableNode
}