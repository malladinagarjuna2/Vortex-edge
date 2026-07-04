
package runtime

import "vortex-edge/internal/registery"

type Runtime struct {
	Run(service*models.Service) error 

	Stop(service* modelsService) error
	
	Delete(service* modelsService) error

	Logs(service* modelsService) (string, error)

}