package  models 
import "time"
type ServiceStatus string 
//service ke baare me bata rahe hai 
const (
	Deploying ServiceStatus = "deploying"
	Running  ServiceStatus = "running"
	Stopped  ServiceStatus = "stopped"
	Failed   ServiceStatus = "failed"
	Deleting ServiceStatus = "deleting"
)

type Service struct {
	ID          string
	Name        string
	Image       string
	ContainerID string
	HostPort    int
    CPU    int
	Memory int64
    NodeID string
	Status ServiceStatus

	CreatedAt time.Time
	UpdatedAt time.Time
}
// Service
//  ├── CPU = 2
//  ├── Memory = 1GB
//  │
//  ↓
// Scheduler
//  │
//  ↓
// Node-02
//  │
//  └── service.NodeID = "node-02"