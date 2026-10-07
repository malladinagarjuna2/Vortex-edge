package orchestrator

import (
	"errors"
	"fmt"

	"vortex-edge/internal/agent"
	"vortex-edge/internal/cluster"
	"vortex-edge/internal/models"
	"vortex-edge/internal/registery"
	"vortex-edge/internal/scheduler"
	"vortex-edge/proto"
)

type Orchestrator struct {
	scheduler    *scheduler.Scheduler
	clusterState *cluster.ClusterState

	// registry remembers every deployed service. Without it the control
	// plane can't answer "which services were running on the failed node?"
	registry registery.ServiceRegistry
}

func NewOrchestrator(
	scheduler *scheduler.Scheduler,
	clusterState *cluster.ClusterState,
	registry registery.ServiceRegistry,
) *Orchestrator {
	return &Orchestrator{
		scheduler:    scheduler,
		clusterState: clusterState,
		registry:     registry,
	}
}

func (o *Orchestrator) Deploy(service *models.Service) error {
	// Register first so a duplicate service ID is rejected before we
	// reserve resources or start a container.
	service.Status = models.Deploying
	err := o.registry.Add(service)
	if err != nil {
		return err
	}

	_, err = o.placeService(service, "")
	if err != nil {
		_ = o.registry.Delete(service.ID)
		service.Status = models.Failed
		service.LastError = err.Error()
		return err
	}

	return o.registry.Update(service)
}

// placeService schedules the service onto a Ready node (never excludeNodeID),
// reserves resources there and starts it through the node's agent.
// On success the service's NodeID/ContainerID/Status are updated.
// On failure the reserved resources are released and the service is untouched.
func (o *Orchestrator) placeService(
	service *models.Service,
	excludeNodeID string,
) (*models.Node, error) {
	node, err := o.scheduler.Schedule(
		service.CPU,
		service.Memory,
	)
	if err != nil {
		return nil, err
	}

	if node.ID == excludeNodeID {
		return nil, fmt.Errorf(
			"scheduler selected failed node %s",
			excludeNodeID,
		)
	}

	// Reserve resources before starting the container.
	err = o.clusterState.AllocateResources(
		node.ID,
		service.CPU,
		service.Memory,
	)
	if err != nil {
		return nil, err
	}

	release := func() {
		_ = o.clusterState.ReleaseResources(
			node.ID,
			service.CPU,
			service.Memory,
		)
	}

	client, err := agent.NewClient(node.Address)
	if err != nil {
		release()
		return nil, err
	}
	defer func() { _ = client.Close() }()

	response, err := client.RunService(
		&proto.RunServiceRequest{
			ServiceId: service.ID,
			Name:      service.Name,
			Image:     service.Image,
			Cpu:       int32(service.CPU),
			Memory:    service.Memory,
			HostPort:  int32(service.HostPort),
		},
	)
	if err != nil {
		release()
		return nil, err
	}

	if !response.Started {
		release()
		return nil, fmt.Errorf(
			"service %s failed to start on node %s",
			service.ID,
			node.ID,
		)
	}

	service.NodeID = node.ID
	service.ContainerID = response.ContainerId
	service.Status = models.Running
	service.ResourcesAllocated = true
	service.LastError = ""

	return node, nil
}

func (o *Orchestrator) Delete(service *models.Service) error {
	node, err := o.clusterState.NodeByID(service.NodeID)
	if err != nil {
		return err
	}
	if node == nil {
		return fmt.Errorf("node %s not found", service.NodeID)
	}

	client, err := agent.NewClient(node.Address)
	if err != nil {
		return err
	}
	defer client.Close()

	response, err := client.DeleteService(
		&proto.DeleteServiceRequest{
			ServiceId:   service.ID,
			ContainerId: service.ContainerID,
		},
	)
	if err != nil {
		return err
	}

	if !response.Deleted {
		return fmt.Errorf("service %s failed to delete", service.ID)
	}

	err = o.releaseResources(service)
	if err != nil {
		return err
	}

	service.ContainerID = ""
	service.Status = models.Stopped

	// The service no longer exists, so stop tracking it. This also lets
	// the same service ID be deployed again later.
	return o.registry.Delete(service.ID)
}

func (o *Orchestrator) Stop(service *models.Service) error {
	node, err := o.clusterState.NodeByID(service.NodeID)
	if err != nil {
		return err
	}
	if node == nil {
		return fmt.Errorf("node %s not found", service.NodeID)
	}

	client, err := agent.NewClient(node.Address)
	if err != nil {
		return err
	}
	defer client.Close()

	response, err := client.StopService(
		&proto.StopServiceRequest{
			ServiceId:   service.ID,
			ContainerId: service.ContainerID,
		},
	)
	if err != nil {
		return err
	}

	if !response.Stopped {
		return fmt.Errorf("service %s failed to stop", service.ID)
	}

	err = o.releaseResources(service)
	if err != nil {
		return err
	}

	service.Status = models.Stopped

	return o.registry.Update(service)
}

// releaseResources gives a service's CPU/memory back to its node.
// It only releases once, so calling Stop and then Delete doesn't
// double-count the release.
func (o *Orchestrator) releaseResources(service *models.Service) error {
	if !service.ResourcesAllocated {
		return nil
	}

	err := o.clusterState.ReleaseResources(
		service.NodeID,
		service.CPU,
		service.Memory,
	)
	if err != nil {
		return err
	}

	service.ResourcesAllocated = false

	return nil
}

func (o *Orchestrator) Logs(service *models.Service) (string, error) {
	node, err := o.clusterState.NodeByID(service.NodeID)
	if err != nil {
		return "", err
	}

	if node == nil {
		return "", fmt.Errorf(
			"node %s not found",
			service.NodeID,
		)
	}

	client, err := agent.NewClient(node.Address)
	if err != nil {
		return "", err
	}
	defer client.Close()

	response, err := client.Logs(
		&proto.LogsRequest{
			ServiceId:   service.ID,
			ContainerId: service.ContainerID,
		},
	)
	if err != nil {
		return "", err
	}

	return response.Logs, nil
}

// ServicesOnNode returns the running services placed on a node.
//
// Node 1 failed
//     ↓
// ServicesOnNode("node-01")
//     ↓
// [service-A, service-B, service-C]
func (o *Orchestrator) ServicesOnNode(nodeID string) []*models.Service {
	services := o.registry.List()

	affected := make([]*models.Service, 0)

	for _, service := range services {
		if service.NodeID == nodeID &&
			service.Status == models.Running {
			affected = append(affected, service)
		}
	}

	return affected
}

// RescheduleService moves a service off its (failed) node onto a healthy one.
//
// failed Node
//     ↓
// Scheduler → healthy Node
//     ↓
// allocate resources → gRPC RunService
//     ↓
// update NodeID + ContainerID → Running
//
// If no node can take it, the service is marked Failed with LastError set.
func (o *Orchestrator) RescheduleService(service *models.Service) error {
	oldNodeID := service.NodeID
	hadAllocation := service.ResourcesAllocated

	_, err := o.placeService(service, oldNodeID)
	if err != nil {
		service.Status = models.Failed
		service.LastError = err.Error()
		_ = o.registry.Update(service)

		return fmt.Errorf("reschedule service %s: %w", service.ID, err)
	}

	// The service now lives on the new node, so stop counting it against
	// the old one. The old node is Offline so the scheduler ignores it
	// anyway, but if it comes back it shouldn't carry phantom allocations.
	if hadAllocation {
		_ = o.clusterState.ReleaseResources(
			oldNodeID,
			service.CPU,
			service.Memory,
		)
	}

	return o.registry.Update(service)
}

// RescheduleNodeServices reschedules every running service on a failed node.
// It keeps going when one service fails and returns all errors joined.
func (o *Orchestrator) RescheduleNodeServices(nodeID string) error {
	services := o.ServicesOnNode(nodeID)

	var errs []error

	for _, service := range services {
		err := o.RescheduleService(service)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
