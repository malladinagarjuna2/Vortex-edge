// Next → Cluster State

// Now we need to maintain a current view of the entire cluster.

// For example:

// Cluster
// │
// ├── node-01
// │   ├── Ready
// │   ├── 8 CPU
// │   └── 16 GB RAM
// │
// ├── node-02
// │   ├── Ready
// │   ├── 16 CPU
// │   └── 32 GB RAM
// │
// └── node-03
//     ├── Offline
//     ├── 8 CPU
//     └── 16 GB RAM

// Why do we need this?

// Because the Scheduler is coming next.

// The scheduler will eventually ask:

// "I have a workload requiring 8 CPU and 12 GB RAM. Which available node can run it?"

// It needs reliable cluster information to answer that.

// So we'll create:
// internal/cluster/state.go


package cluster
import(
	"vortex-edge/internal/models"
	"vortex-edge/internal/registery"
)
type ClusterState struct{
	registry *registery.NodeRegistry
}

func NewClusterState(registry *registery.NodeRegistry)*ClusterState{
	 return &ClusterState{
		registry: registry,
	 }
}

func (c *ClusterState) Nodes() []*models.Node {
	return c.registry.List()
}

func (c *ClusterState) NodeByID(id string) (*models.Node, error) {
    return c.registry.Get(id)
}
func (c *ClusterState) ReadyNodes() []*models.Node {
	nodes := c.registry.List()

	readyNodes := make([]*models.Node, 0)

	for _, node := range nodes {
		if node.Status == models.NodeReady {
			readyNodes = append(readyNodes, node)
		}
	}

	return readyNodes
}