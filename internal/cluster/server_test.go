package cluster

import (
	"context"
	"testing"

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