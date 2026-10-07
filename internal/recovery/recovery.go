// Node failure recovery.
//
// Node 1
//   ↓
// Heartbeat timeout
//   ↓
// HealthMonitor.Check()  → [node-01]
//   ↓
// Rescheduler.RescheduleNodeServices("node-01")
//   ↓
// Scheduler → Node 2 → gRPC RunService → Docker
//
// This lives in its own package (not in cluster) because it needs both
// cluster and orchestrator, and orchestrator already imports cluster.
// Putting it in cluster would create the cycle cluster → orchestrator → cluster.
package recovery

import (
	"context"
	"log"
	"time"

	"vortex-edge/internal/models"
)

// FailureDetector reports nodes that have just gone offline.
// Implemented by *cluster.HealthMonitor.
type FailureDetector interface {
	Check() []*models.Node
}

// Rescheduler moves every service off a failed node.
// Implemented by *orchestrator.Orchestrator.
type Rescheduler interface {
	RescheduleNodeServices(nodeID string) error
}

type Manager struct {
	detector    FailureDetector
	rescheduler Rescheduler
}

func NewManager(detector FailureDetector, rescheduler Rescheduler) *Manager {
	return &Manager{
		detector:    detector,
		rescheduler: rescheduler,
	}
}

// RunOnce detects newly failed nodes and reschedules their services.
// It returns the nodes that were handled.
func (m *Manager) RunOnce() []*models.Node {
	failedNodes := m.detector.Check()

	for _, node := range failedNodes {
		log.Printf("node %s is offline, rescheduling its services", node.ID)

		err := m.rescheduler.RescheduleNodeServices(node.ID)
		if err != nil {
			log.Printf(
				"failed to reschedule services from node %s: %v",
				node.ID,
				err,
			)
		}
	}

	return failedNodes
}

// Run calls RunOnce every interval until ctx is cancelled.
func (m *Manager) Run(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.RunOnce()
		}
	}
}
