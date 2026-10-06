package scheduler

import (
	"errors"
	"testing"

	"vortex-edge/internal/cluster"
	"vortex-edge/internal/models"
	"vortex-edge/internal/registery"
)

const gb = int64(1024 * 1024 * 1024)

// TestResourceAccounting walks one node through allocate → over-allocate →
// release, checking AllocatedCPU/AllocatedMemory at every step.
func TestResourceAccounting(t *testing.T) {
	registry := registery.NewNodeRegistry()

	node := &models.Node{
		ID:     "node-01",
		Name:   "worker-01",
		Status: models.NodeReady,
		CPU:    8,
		Memory: 16 * gb,
	}

	err := registry.Add(node)
	if err != nil {
		t.Fatalf("failed to add node: %v", err)
	}

	clusterState := cluster.NewClusterState(registry)
	s := NewScheduler(clusterState)

	expectAllocated := func(step string, cpu int, memory int64) {
		t.Helper()

		if node.AllocatedCPU != cpu || node.AllocatedMemory != memory {
			t.Fatalf(
				"%s: expected allocated %d CPU / %d GB, got %d CPU / %d GB",
				step,
				cpu, memory/gb,
				node.AllocatedCPU, node.AllocatedMemory/gb,
			)
		}
	}

	// schedule finds a node for the request and reserves it, like Deploy does.
	schedule := func(step string, cpu int, memory int64) {
		t.Helper()

		selected, err := s.Schedule(cpu, memory)
		if err != nil {
			t.Fatalf("%s: schedule failed: %v", step, err)
		}

		err = clusterState.AllocateResources(selected.ID, cpu, memory)
		if err != nil {
			t.Fatalf("%s: allocate failed: %v", step, err)
		}
	}

	// Test 1 — allocation
	schedule("first deploy", 2, 4*gb)
	expectAllocated("first deploy", 2, 4*gb)

	// Test 2 — second deployment
	schedule("second deploy", 4, 8*gb)
	expectAllocated("second deploy", 6, 12*gb)

	// Test 3 — over-allocation: only 2 CPU left, 3 requested
	_, err = s.Schedule(3, 2*gb)
	if !errors.Is(err, ErrNoSuitableNode) {
		t.Fatalf("over-allocation: expected ErrNoSuitableNode, got %v", err)
	}

	err = clusterState.AllocateResources(node.ID, 3, 2*gb)
	if err == nil {
		t.Fatal("over-allocation: expected AllocateResources to fail")
	}
	expectAllocated("over-allocation", 6, 12*gb)

	// Test 4 — stop the second service
	err = clusterState.ReleaseResources(node.ID, 4, 8*gb)
	if err != nil {
		t.Fatalf("stop: release failed: %v", err)
	}
	expectAllocated("stop", 2, 4*gb)

	// Test 5 — delete the remaining service
	err = clusterState.ReleaseResources(node.ID, 2, 4*gb)
	if err != nil {
		t.Fatalf("delete: release failed: %v", err)
	}
	expectAllocated("delete", 0, 0)
}
