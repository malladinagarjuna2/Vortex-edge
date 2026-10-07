package registery

import "vortex-edge/internal/models"

// ServiceRegistry tracks every service the control plane has deployed,
// so it can answer "which services are running on node X?" when that
// node fails. MemoryRegistery is the in-memory implementation.
type ServiceRegistry interface {
	Add(service *models.Service) error
	Get(id string) (*models.Service, error)
	Update(service *models.Service) error
	Delete(id string) error
	List() []*models.Service
}

var _ ServiceRegistry = (*MemoryRegistery)(nil)
