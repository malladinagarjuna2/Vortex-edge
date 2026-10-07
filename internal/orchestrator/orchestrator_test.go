package orchestrator

import (
	"net"
	"testing"

	"google.golang.org/grpc"

	"vortex-edge/internal/agent"
	"vortex-edge/internal/cluster"
	"vortex-edge/internal/models"
	"vortex-edge/internal/registery"
	"vortex-edge/internal/runtime"
	"vortex-edge/internal/scheduler"
	"vortex-edge/proto"
)

func TestDeploy(t *testing.T) {

	// Docker runtime
	dockerRuntime, err := runtime.NewDockerRuntime()
	if err != nil {
		t.Fatalf("failed to create docker runtime: %v", err)
	}

	// Node
	node := &models.Node{
		ID:      "node-01",
		Name:    "worker-01",
		Address: "localhost:9003",
		Status:  models.NodeReady,
		CPU:     4,
		Memory:  8 * 1024 * 1024 * 1024,
	}

	// Agent
	nodeAgent := agent.NewAgent(
		node,
		dockerRuntime,
	)

	// gRPC server
	listener, err := net.Listen(
		"tcp",
		node.Address,
	)
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	proto.RegisterNodeAgentServer(
		grpcServer,
		nodeAgent,
	)

	go func() {
		_ = grpcServer.Serve(listener)
	}()

	defer grpcServer.Stop()
	defer listener.Close()

	// Registry
	registry := registery.NewNodeRegistry()

	err = registry.Add(node)
	if err != nil {
		t.Fatalf("failed to add node: %v", err)
	}

	// Cluster state
	clusterState := cluster.NewClusterState(
		registry,
	)

	// Scheduler
	s := scheduler.NewScheduler(
		clusterState,
	)

	// Orchestrator
	o := NewOrchestrator(
		s,
		clusterState,
		registery.NewMemoryRegistery(),
	)

	// Service
	service := &models.Service{
		ID:     "service-01",
		Name:   "vortex-orchestrator-test",
		Image:  "hello-world",
		CPU:    1,
		Memory: 256 * 1024 * 1024,
	}

	// Deploy
	err = o.Deploy(service)
	if err != nil {
		t.Fatalf("deployment failed: %v", err)
	}
	// Verify selected node
	if service.NodeID != node.ID {
		t.Fatalf(
			"expected node %s, got %s",
			node.ID,
			service.NodeID,
		)
	}

	// Verify container
	if service.ContainerID == "" {
		t.Fatal("expected container ID")
	}

	// Verify status
	if service.Status != models.Running {
		t.Fatalf(
			"expected Running, got %s",
			service.Status,
		)
	}

	// Verify resources were reserved
	if node.AllocatedCPU != service.CPU || node.AllocatedMemory != service.Memory {
		t.Fatalf(
			"expected allocated %d CPU / %d memory, got %d / %d",
			service.CPU, service.Memory,
			node.AllocatedCPU, node.AllocatedMemory,
		)
	}

	t.Logf(
		"Deployment successful: node=%s container=%s",
		service.NodeID,
		service.ContainerID,
	)

	// Delete
	err = o.Delete(service)
	if err != nil {
		t.Fatalf("deletion failed: %v", err)
	}

	if service.ContainerID != "" {
		t.Fatal("expected container ID to be empty after deletion")
	}
	if service.Status != models.Stopped {
    t.Fatalf("expected service status to be Stopped, got %s", service.Status)
}

	// Verify resources were released
	if node.AllocatedCPU != 0 || node.AllocatedMemory != 0 {
		t.Fatalf(
			"expected resources released, got %d CPU / %d memory",
			node.AllocatedCPU, node.AllocatedMemory,
		)
	}
}