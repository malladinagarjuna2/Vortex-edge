package cluster

import (
	"testing"
	"time"

	"vortex-edge/internal/models"
	"vortex-edge/internal/registery"
)

func TestHealthMonitor(t *testing.T) {
	registry := registery.NewNodeRegistry()

	node := &models.Node{
		ID:            "node-01",
		Name:          "worker-01",
		Status:        models.NodeReady,
		LastHeartbeat: time.Now().Add(-10 * time.Second),
	}

	err := registry.Add(node)
	if err != nil {
		t.Fatalf("failed to add node: %v", err)
	}

	// Node is considered offline if heartbeat is older than 5 seconds.
	monitor := NewHealthMonitor(
		registry,
		5*time.Second,
	)

	monitor.Check()

	updatedNode, err := registry.Get("node-01")
	if err != nil {
		t.Fatalf("failed to get node: %v", err)
	}

	if updatedNode.Status != models.NodeOffline {
		t.Fatalf(
			"expected node to be Offline, got %s",
			updatedNode.Status,
		)
	}
}