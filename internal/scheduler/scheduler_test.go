// Service Request
//       ↓
// Scheduler
//       ↓
// Node-02 selected
//       ↓
// gRPC
//       ↓
// Node Agent
//       ↓
// Docker Runtime
//       ↓
// Container
package scheduler

import (
	"testing"

	"vortex-edge/internal/cluster"
	"vortex-edge/internal/models"
	"vortex-edge/internal/registery"
)

func TestScheduler(t *testing.T) {
	registry := registery.NewNodeRegistry()

	node1 := &models.Node{
		ID:     "node-01",
		Name:   "worker-01",
		Status: models.NodeReady,
		CPU:    4,
		Memory: 8 * 1024 * 1024 * 1024,
	}

	node2 := &models.Node{
		ID:     "node-02",
		Name:   "worker-02",
		Status: models.NodeReady,
		CPU:    8,
		Memory: 16 * 1024 * 1024 * 1024,
	}

	err := registry.Add(node1)
	if err != nil {
		t.Fatalf("failed to add node1: %v", err)
	}

	err = registry.Add(node2)
	if err != nil {
		t.Fatalf("failed to add node2: %v", err)
	}

	clusterState := cluster.NewClusterState(registry)

	scheduler := NewScheduler(clusterState)

	selectedNode, err := scheduler.Schedule(
		8,
		12*1024*1024*1024,
	)
	if err != nil {
		t.Fatalf("failed to schedule workload: %v", err)
	}

	if selectedNode.ID != "node-02" {
		t.Fatalf(
			"expected node-02, got %s",
			selectedNode.ID,
		)
	}
}