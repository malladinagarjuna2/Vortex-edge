
package registery 

import "vortex-edge/internal/models"

type Registery interface {
	AddService(service *models.Service) error
	GetService(id string) (*models.Service, error)
	UpdateService(service *models.Service) error
	DeleteService(id string) error
	List() []*models.Service
}