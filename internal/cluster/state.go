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
	"fmt"
	"sync"
	"time"

	"vortex-edge/internal/models"
	"vortex-edge/internal/registery"
)
type ClusterState struct{
	registry *registery.NodeRegistry

	// mu makes check-and-allocate atomic so two deploys can't both
	// claim the same remaining capacity.
	mu sync.Mutex
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
// AllocateResources reserves cpu/memory on a node, failing if the node
// doesn't have enough available capacity left.
func (c *ClusterState) AllocateResources(
	nodeID string,
	cpu int,
	memory int64,
) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, err := c.registry.Get(nodeID)
	if err != nil {
		return err
	}

	if node == nil {
		return fmt.Errorf("node %s not found", nodeID)
	}

	availableCPU := node.CPU - node.AllocatedCPU
	availableMemory := node.Memory - node.AllocatedMemory

	if availableCPU < cpu {
		return fmt.Errorf("insufficient CPU on node %s", nodeID)
	}

	if availableMemory < memory {
		return fmt.Errorf("insufficient memory on node %s", nodeID)
	}

	node.AllocatedCPU += cpu
	node.AllocatedMemory += memory

	node.UpdatedAt = time.Now()

	return c.registry.Update(node)
}

// ReleaseResources returns previously allocated cpu/memory to a node.
func (c *ClusterState) ReleaseResources(
	nodeID string,
	cpu int,
	memory int64,
) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, err := c.registry.Get(nodeID)
	if err != nil {
		return err
	}

	if node == nil {
		return fmt.Errorf("node %s not found", nodeID)
	}

	node.AllocatedCPU -= cpu
	node.AllocatedMemory -= memory

	// Defensive: never go negative.
	if node.AllocatedCPU < 0 {
		node.AllocatedCPU = 0
	}

	if node.AllocatedMemory < 0 {
		node.AllocatedMemory = 0
	}

	node.UpdatedAt = time.Now()

	return c.registry.Update(node)
}
