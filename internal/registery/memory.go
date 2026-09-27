package registery
import(
	"fmt"
	"sync"

	"vortex-edge/internal/models"
)
// yaha pe tumhaari cheeze store ho rahi hai 
// jb tumhara node register kiya hai toh woh thode der ke liye mutex lock me store hoga 
type MemoryRegistery struct {
	services map[string]*models.Service
	mu  sync.RWMutex
}                   

func NewMemoryRegistery() *MemoryRegistery {
	return &MemoryRegistery{
		services: make(map[string]*models.Service),
	}
}

func (r*MemoryRegistery) Add(service *models.Service) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.services[service.ID]; exists{
		return fmt.Errorf("service with ID %s already exists", service.ID)
	}

	r.services[service.ID] = service
	return nil
}


func(r *MemoryRegistery) Get(id string) (*models.Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	service, exists := r.services[id]
		if !exists {
		return nil, fmt.Errorf("service with ID %s not found", id)
	}

	return service, nil
}

func (r *MemoryRegistery) List() []*models.Service{
	r.mu.RLock()
	defer r.mu.RUnlock()

	services := make([]*models.Service, 0, len(r.services))
	for _, service := range r.services {
		services = append(services, service)
	}
	return services
}
func (r *MemoryRegistery) Update(service *models.Service) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.services[service.ID]; !exists {
		return fmt.Errorf("service with ID %s not found", service.ID)
	}

	r.services[service.ID] = service
	return nil
}
func (r *MemoryRegistery) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.services[id]; !exists {
		return fmt.Errorf("service with ID %s not found", id)
	}

	delete(r.services, id)
	return nil
}