package cluster

import (
	"time"

	"vortex-edge/internal/models"
	"vortex-edge/internal/registery"
)

type HealthMonitor struct {
	registry *registery.NodeRegistry
	timeout  time.Duration
}

func NewHealthMonitor(
	registry *registery.NodeRegistry,
	timeout time.Duration,
) *HealthMonitor {
	return &HealthMonitor{
		registry: registry,
		timeout:  timeout,
	}
}

// Check marks nodes whose heartbeat has timed out as Offline and returns
// only the nodes that went Offline during this check, so the caller
// reschedules each failed node's services exactly once.
func (h *HealthMonitor) Check() []*models.Node {
	nodes := h.registry.List()

	now := time.Now()
	failed := make([]*models.Node, 0)

	for _, node := range nodes {
		if node.Status == models.NodeOffline {
			continue
		}

		if now.Sub(node.LastHeartbeat) > h.timeout {
			node.Status = models.NodeOffline
			node.UpdatedAt = now

			_ = h.registry.Update(node)
			failed = append(failed, node)
		}
	}

	return failed
}
