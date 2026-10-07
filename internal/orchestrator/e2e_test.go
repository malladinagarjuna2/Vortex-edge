package orchestrator

import (
	"testing"
	"time"

	"vortex-edge/internal/cluster"
	"vortex-edge/internal/models"
	"vortex-edge/internal/recovery"
	"vortex-edge/internal/registery"
	"vortex-edge/internal/runtime"
	"vortex-edge/internal/scheduler"
)

// TestMultiNodeFailover is the Phase 2 end-to-end test:
//
//	Control Plane ── Orchestrator
//	       ┌────────┴────────┐
//	    Node 1             Node 2
//	    Agent              Agent
//	    Docker             Docker
//
// Deploy → Node 1, Node 1 stops heartbeating, recovery moves the
// service to Node 2 and a real container is started there.
//
// Both agents share the local Docker daemon, so when Node 1 "dies" the
// test removes its container itself. Otherwise the restarted container
// would clash with the old one's name.
func TestMultiNodeFailover(t *testing.T) {
	dockerRuntime, err := runtime.NewDockerRuntime()
	if err != nil {
		t.Skipf("docker unavailable: %v", err)
	}
	if _, err := dockerRuntime.ServerVersion(); err != nil {
		t.Skipf("docker daemon not reachable: %v", err)
	}

	// 1. Registries
	nodes := registery.NewNodeRegistry()
	services := registery.NewMemoryRegistery()

	// 2-3. Nodes
	node1 := &models.Node{
		ID:            "node-01",
		Name:          "worker-01",
		Status:        models.NodeReady,
		CPU:           4,
		Memory:        8 * 1024 * 1024 * 1024,
		LastHeartbeat: time.Now(),
	}
	node2 := &models.Node{
		ID:            "node-02",
		Name:          "worker-02",
		Status:        models.NodeReady,
		CPU:           4,
		Memory:        8 * 1024 * 1024 * 1024,
		LastHeartbeat: time.Now(),
	}

	// 4-5. Agents, each on its own port
	startAgent(t, node1, dockerRuntime)
	startAgent(t, node2, dockerRuntime)

	// 6. Register both nodes
	_ = nodes.Add(node1)
	_ = nodes.Add(node2)

	clusterState := cluster.NewClusterState(nodes)
	orch := NewOrchestrator(
		scheduler.NewScheduler(clusterState),
		clusterState,
		services,
	)

	// 7. Deploy service. node-02 is held back so placement is deterministic.
	service := &models.Service{
		ID:     "service-e2e",
		Name:   "vortex-e2e-failover",
		Image:  "hello-world",
		CPU:    1,
		Memory: 256 * 1024 * 1024,
	}

	node2.Status = models.NodeOffline
	err = orch.Deploy(service)
	node2.Status = models.NodeReady
	if err != nil {
		t.Fatalf("deploy failed: %v", err)
	}

	// 8. Verify service is on Node 1
	if service.NodeID != node1.ID {
		t.Fatalf("expected service on %s, got %s", node1.ID, service.NodeID)
	}
	oldContainerID := service.ContainerID

	// 9. Node 1 stops heartbeating; simulate the machine (and its container) dying.
	node1.LastHeartbeat = time.Now().Add(-time.Minute)
	node2.LastHeartbeat = time.Now()
	_ = dockerRuntime.Delete(&models.Service{ContainerID: oldContainerID})

	// 10. Recovery loop detects the failure and reschedules
	manager := recovery.NewManager(
		cluster.NewHealthMonitor(nodes, 5*time.Second),
		orch,
	)
	failed := manager.RunOnce()
	if len(failed) != 1 || failed[0].ID != node1.ID {
		t.Fatalf("expected only node-01 to fail, got %v", failed)
	}

	t.Cleanup(func() {
		_ = orch.Delete(service)
	})

	// 11. Verify service.NodeID == Node 2
	if service.NodeID != node2.ID {
		t.Fatalf(
			"expected service on %s, got %s (status=%s, lastError=%q)",
			node2.ID, service.NodeID, service.Status, service.LastError,
		)
	}
	if service.Status != models.Running {
		t.Fatalf("expected Running, got %s", service.Status)
	}

	// 12. Verify Node 2 allocation increased (and Node 1's was released)
	if node2.AllocatedCPU != service.CPU || node2.AllocatedMemory != service.Memory {
		t.Fatalf(
			"expected node-02 allocation %d / %d, got %d / %d",
			service.CPU, service.Memory, node2.AllocatedCPU, node2.AllocatedMemory,
		)
	}
	if node1.AllocatedCPU != 0 {
		t.Fatalf("expected node-01 allocation released, got %d CPU", node1.AllocatedCPU)
	}

	// 13. Verify a new container exists on Node 2
	if service.ContainerID == "" || service.ContainerID == oldContainerID {
		t.Fatalf("expected a new container, got %q", service.ContainerID)
	}

	t.Logf(
		"failover ok: %s → %s, container %s",
		node1.ID, node2.ID, service.ContainerID,
	)
}
