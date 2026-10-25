package cluster

import (
	"context"
	"testing"

	"vortex-edge/internal/models"
	"vortex-edge/internal/registery"
	"vortex-edge/proto"
)

func TestRegisterNode(t *testing.T) {
	registry := registery.NewNodeRegistry()

	server := NewServer(registry)

	response, err := server.Register(
		context.Background(),
		&proto.RegisterRequest{
			NodeId:  "node-01",
			Name:    "worker-01",
			Address: "192.168.1.10:9000",
			Cpu:     8,
			Memory: 16 * 1024 * 1024 * 1024,
		},
	)

	if err != nil {
		t.Fatalf("registration failed: %v", err)
	}

	if !response.Accepted {
		t.Fatal("expected registration to be accepted")
	}

	node, err := registry.Get("node-01")
	if err != nil {
		t.Fatalf("failed to get registered node: %v", err)
	}

	if node == nil {
		t.Fatal("expected registered node, got nil")
	}

	if node.Name != "worker-01" {
		t.Fatalf("expected worker-01, got %s", node.Name)
	}

	t.Log("Node registered successfully")
}
func TestReRegisterKeepsAllocationsAndUpdatesAddress(t *testing.T) {
	registry := registery.NewNodeRegistry()
	server := NewServer(registry)

	request := &proto.RegisterRequest{
		NodeId:  "node-01",
		Name:    "worker-01",
		Address: "100.64.0.1:9100",
		Cpu:     8,
		Memory:  16 * 1024 * 1024 * 1024,
	}

	if _, err := server.Register(context.Background(), request); err != nil {
		t.Fatal(err)
	}

	node, _ := registry.Get("node-01")
	if node.LastHeartbeat.IsZero() {
		t.Fatal("expected registration to set LastHeartbeat")
	}

	// A service is running on the node.
	node.AllocatedCPU = 2
	node.Status = models.NodeOffline

	// The worker reconnects from a new address.
	request.Address = "100.64.0.9:9100"
	if _, err := server.Register(context.Background(), request); err != nil {
		t.Fatal(err)
	}

	node, _ = registry.Get("node-01")
	if node.AllocatedCPU != 2 {
		t.Fatalf("expected allocations kept, got %d CPU", node.AllocatedCPU)
	}
	if node.Address != "100.64.0.9:9100" {
		t.Fatalf("expected new address, got %s", node.Address)
	}
	if node.Status != models.NodeReady {
		t.Fatalf("expected Ready, got %s", node.Status)
	}
}
