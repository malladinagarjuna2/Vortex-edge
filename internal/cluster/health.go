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

func (h *HealthMonitor) Check() {
	nodes := h.registry.List()

	now := time.Now()

	for _, node := range nodes {
		if now.Sub(node.LastHeartbeat) > h.timeout {
			node.Status = models.NodeOffline
			node.UpdatedAt = now

			_ = h.registry.Update(node)
		}
	}
}