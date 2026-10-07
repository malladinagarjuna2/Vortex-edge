// Service Request
//       ↓
// Scheduler
//       ↓
// Node-02 selected
//       ↓
// gRPC
//       ↓
// Node Agent
//       ↓
// Docker Runtime
//       ↓
// Container
package scheduler

import (
	"errors"
	"testing"

	"vortex-edge/internal/cluster"
	"vortex-edge/internal/models"
	"vortex-edge/internal/registery"
)

func TestScheduler(t *testing.T) {
	registry := registery.NewNodeRegistry()

	node1 := &models.Node{
		ID:     "node-01",
		Name:   "worker-01",
		Status: models.NodeReady,
		CPU:    4,
		Memory: 8 * 1024 * 1024 * 1024,
	}

	node2 := &models.Node{
		ID:     "node-02",
		Name:   "worker-02",
		Status: models.NodeReady,
		CPU:    8,
		Memory: 16 * 1024 * 1024 * 1024,
	}

	err := registry.Add(node1)
	if err != nil {
		t.Fatalf("failed to add node1: %v", err)
	}

	err = registry.Add(node2)
	if err != nil {
		t.Fatalf("failed to add node2: %v", err)
	}

	clusterState := cluster.NewClusterState(registry)

	scheduler := NewScheduler(clusterState)

	selectedNode, err := scheduler.Schedule(
		8,
		12*1024*1024*1024,
	)
	if err != nil {
		t.Fatalf("failed to schedule workload: %v", err)
	}

	if selectedNode.ID != "node-02" {
		t.Fatalf(
			"expected node-02, got %s",
			selectedNode.ID,
		)
	}
}
//insufficient cpu test 
func TestScheduleInsufficientCPU(t *testing.T) {
    registry := registery.NewNodeRegistry()

    node := &models.Node{
        ID:            "node-01",
        Name:          "node-01",
        Status:        models.NodeReady,
        CPU:           8,
        Memory:        16 * 1024 * 1024 * 1024,
        AllocatedCPU:  6,
        AllocatedMemory: 4 * 1024 * 1024 * 1024,
    }

    if err := registry.Add(node); err != nil {
        t.Fatal(err)
    }

    clusterState := cluster.NewClusterState(registry)
    scheduler := NewScheduler(clusterState)

    _, err := scheduler.Schedule(
        3,
        2*1024*1024*1024,
    )

    if !errors.Is(err, ErrNoSuitableNode) {
        t.Fatalf(
            "expected ErrNoSuitableNode, got %v",
            err,

        )
    }
}
// insufficient memory test 
func TestScheduleInsufficientMemory(t*testing.T){
	registery:= registery.NewNodeRegistry()
  node := &models.Node{
	    ID:              "node-01",
        Name:            "node-01",
        Status:          models.NodeReady,
        CPU:             8,
        Memory:          16 * 1024 * 1024 * 1024,
        AllocatedCPU:    2,
        AllocatedMemory: 14 * 1024 * 1024 * 1024,
  }
  if err:= registery.Add(node); err!= nil{
	 t.Fatal(err)
  }

  clusterState := cluster.NewClusterState(registery)
  scheduler:= NewScheduler(clusterState)
  
    _, err := scheduler.Schedule(
        2,
        4*1024*1024*1024,
    )

    if !errors.Is(err, ErrNoSuitableNode) {
        t.Fatalf(
            "expected ErrNoSuitableNode, got %v",
            err,
        )
    }
}	

func TestScheduleSkipsOverloadedNode(t *testing.T){
	 registry:= registery.NewNodeRegistry()
	 overloaded := &models.Node{
		    ID:              "node-01",
        Name:            "node-01",
        Status:          models.NodeReady,
        CPU:             8,
        Memory:          16 * 1024 * 1024 * 1024,
        AllocatedCPU:    8,
        AllocatedMemory: 8 * 1024 * 1024 * 1024,
	 }

	  available := &models.Node{
        ID:              "node-02",
        Name:            "node-02",
        Status:          models.NodeReady,
        CPU:             8,
        Memory:          16 * 1024 * 1024 * 1024,
        AllocatedCPU:    2,
        AllocatedMemory: 2 * 1024 * 1024 * 1024,
    }

    if err := registry.Add(overloaded); err != nil {
        t.Fatal(err)
    }

    if err := registry.Add(available); err != nil {
        t.Fatal(err)
    }

    clusterState := cluster.NewClusterState(registry)
    scheduler := NewScheduler(clusterState)

    node, err := scheduler.Schedule(
        4,
        4*1024*1024*1024,
    )

    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if node.ID != "node-02" {
        t.Fatalf(
            "expected node-02, got %s",
            node.ID,
        )
    }
      
}