
package runtime
import models "vortex-edge/internal/models"

// import "vortex-edge/internal/registery"

type Runtime interface {
	Run(service*models.Service) error 

	Stop(service*models.Service) error
	
	Delete(service* models.Service) error

	Logs(service* models.Service) (string, error)

}