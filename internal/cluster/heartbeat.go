package cluster

import (
	"time"

	"vortex-edge/internal/models"
	"vortex-edge/internal/registery"
)

type HeartbeatManager struct {
	registry *registery.NodeRegistry
}

func NewHeartbeatManager(registry *registery.NodeRegistry) *HeartbeatManager {
	return &HeartbeatManager{
		registry: registry,
	}
}

func (h *HeartbeatManager) Handle(nodeID string) error {
	node, err := h.registry.Get(nodeID)
	if err != nil {
		return err
	}

	if node == nil {
		return nil
	}

	node.LastHeartbeat = time.Now()
	node.Status = models.NodeReady
	node.UpdatedAt = time.Now()

	return h.registry.Update(node)
}