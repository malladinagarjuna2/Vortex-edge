package  models 
import "time"
type ServiceStatus string 

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

	Status ServiceStatus

	CreatedAt time.Time
	UpdatedAt time.Time
}