package recovery

import (
	"errors"
	"testing"
	"time"

	"vortex-edge/internal/cluster"
	"vortex-edge/internal/models"
	"vortex-edge/internal/registery"
)

type fakeRescheduler struct {
	calls []string
	err   error
}

func (f *fakeRescheduler) RescheduleNodeServices(nodeID string) error {
	f.calls = append(f.calls, nodeID)
	return f.err
}

func TestRunOnceReschedulesFailedNodesOnce(t *testing.T) {
	registry := registery.NewNodeRegistry()

	dead := &models.Node{
		ID:            "node-01",
		Status:        models.NodeReady,
		LastHeartbeat: time.Now().Add(-time.Minute),
	}
	alive := &models.Node{
		ID:            "node-02",
		Status:        models.NodeReady,
		LastHeartbeat: time.Now(),
	}
	_ = registry.Add(dead)
	_ = registry.Add(alive)

	rescheduler := &fakeRescheduler{}
	manager := NewManager(
		cluster.NewHealthMonitor(registry, 5*time.Second),
		rescheduler,
	)

	failed := manager.RunOnce()
	if len(failed) != 1 || failed[0].ID != "node-01" {
		t.Fatalf("expected node-01 to fail, got %v", failed)
	}
	if len(rescheduler.calls) != 1 || rescheduler.calls[0] != "node-01" {
		t.Fatalf("expected reschedule of node-01, got %v", rescheduler.calls)
	}

	// node-01 is already Offline, so a second pass must not reschedule again.
	manager.RunOnce()
	if len(rescheduler.calls) != 1 {
		t.Fatalf("expected no further reschedules, got %v", rescheduler.calls)
	}
}

func TestRunOnceContinuesAfterRescheduleError(t *testing.T) {
	registry := registery.NewNodeRegistry()
	_ = registry.Add(&models.Node{ID: "node-01", Status: models.NodeReady})
	_ = registry.Add(&models.Node{ID: "node-02", Status: models.NodeReady})

	rescheduler := &fakeRescheduler{err: errors.New("no capacity")}
	manager := NewManager(
		cluster.NewHealthMonitor(registry, 5*time.Second),
		rescheduler,
	)

	manager.RunOnce()

	if len(rescheduler.calls) != 2 {
		t.Fatalf("expected both nodes handled despite errors, got %v", rescheduler.calls)
	}
}
