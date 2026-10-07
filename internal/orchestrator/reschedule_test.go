package orchestrator

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"google.golang.org/grpc"

	"vortex-edge/internal/agent"
	"vortex-edge/internal/cluster"
	"vortex-edge/internal/models"
	"vortex-edge/internal/registery"
	"vortex-edge/internal/runtime"
	"vortex-edge/internal/scheduler"
	"vortex-edge/proto"
)

// fakeRuntime stands in for Docker so rescheduling can be tested without
// a Docker daemon. Requests still travel over real gRPC to a real Agent.
type fakeRuntime struct {
	nodeID  string
	running map[string]bool
}

func newFakeRuntime(nodeID string) *fakeRuntime {
	return &fakeRuntime{
		nodeID:  nodeID,
		running: make(map[string]bool),
	}
}

func (f *fakeRuntime) Run(service *models.Service) error {
	service.ContainerID = fmt.Sprintf("%s-%s", f.nodeID, service.ID)
	service.Status = models.Running
	f.running[service.ContainerID] = true
	return nil
}

func (f *fakeRuntime) Stop(service *models.Service) error {
	delete(f.running, service.ContainerID)
	return nil
}

func (f *fakeRuntime) Delete(service *models.Service) error {
	delete(f.running, service.ContainerID)
	return nil
}

func (f *fakeRuntime) Logs(service *models.Service) (string, error) {
	return "", nil
}

// startAgent serves a node agent on a free local port and sets node.Address.
func startAgent(t *testing.T, node *models.Node, rt runtime.Runtime) {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	node.Address = listener.Addr().String()

	grpcServer := grpc.NewServer()
	proto.RegisterNodeAgentServer(grpcServer, agent.NewAgent(node, rt))

	go func() {
		_ = grpcServer.Serve(listener)
	}()

	t.Cleanup(grpcServer.Stop)
}

type testCluster struct {
	orch     *Orchestrator
	nodes    *registery.NodeRegistry
	services *registery.MemoryRegistery
	node1    *models.Node
	node2    *models.Node
	rt1      *fakeRuntime
	rt2      *fakeRuntime
}

// newTestCluster builds two Ready nodes with running agents. node-01 has
// more free capacity so the scheduler picks it first.
func newTestCluster(t *testing.T, node2CPU int) *testCluster {
	t.Helper()

	node1 := &models.Node{
		ID:     "node-01",
		Name:   "worker-01",
		Status: models.NodeReady,
		CPU:    8,
		Memory: 16 * 1024 * 1024 * 1024,
	}
	node2 := &models.Node{
		ID:     "node-02",
		Name:   "worker-02",
		Status: models.NodeReady,
		CPU:    node2CPU,
		Memory: 16 * 1024 * 1024 * 1024,
	}

	rt1 := newFakeRuntime(node1.ID)
	rt2 := newFakeRuntime(node2.ID)
	startAgent(t, node1, rt1)
	startAgent(t, node2, rt2)

	nodes := registery.NewNodeRegistry()
	_ = nodes.Add(node1)
	_ = nodes.Add(node2)

	clusterState := cluster.NewClusterState(nodes)
	services := registery.NewMemoryRegistery()

	return &testCluster{
		orch: NewOrchestrator(
			scheduler.NewScheduler(clusterState),
			clusterState,
			services,
		),
		nodes:    nodes,
		services: services,
		node1:    node1,
		node2:    node2,
		rt1:      rt1,
		rt2:      rt2,
	}
}

// deployOnNode1 deploys while node-02 is temporarily not Ready, so the
// placement is deterministic (the scheduler iterates a map).
func (c *testCluster) deployOnNode1(t *testing.T, service *models.Service) {
	t.Helper()

	c.node2.Status = models.NodeOffline
	err := c.orch.Deploy(service)
	c.node2.Status = models.NodeReady

	if err != nil {
		t.Fatalf("deploy failed: %v", err)
	}
	if service.NodeID != c.node1.ID {
		t.Fatalf("expected service on %s, got %s", c.node1.ID, service.NodeID)
	}
}

func TestRescheduleService(t *testing.T) {
	c := newTestCluster(t, 8)

	service := &models.Service{
		ID:     "service-01",
		Name:   "web",
		Image:  "nginx",
		CPU:    2,
		Memory: 512 * 1024 * 1024,
	}
	c.deployOnNode1(t, service)

	// node-01 fails.
	c.node1.Status = models.NodeOffline

	err := c.orch.RescheduleService(service)
	if err != nil {
		t.Fatalf("reschedule failed: %v", err)
	}

	if service.NodeID != c.node2.ID {
		t.Fatalf("expected service on %s, got %s", c.node2.ID, service.NodeID)
	}
	if service.Status != models.Running {
		t.Fatalf("expected Running, got %s", service.Status)
	}
	if service.ContainerID != "node-02-service-01" {
		t.Fatalf("expected container on node-02, got %q", service.ContainerID)
	}
	if !c.rt2.running[service.ContainerID] {
		t.Fatal("expected container to be running in node-02's runtime")
	}

	// Resources moved from node-01 to node-02.
	if c.node2.AllocatedCPU != 2 || c.node2.AllocatedMemory != service.Memory {
		t.Fatalf(
			"expected node-02 allocation 2 CPU / %d memory, got %d / %d",
			service.Memory, c.node2.AllocatedCPU, c.node2.AllocatedMemory,
		)
	}
	if c.node1.AllocatedCPU != 0 || c.node1.AllocatedMemory != 0 {
		t.Fatalf(
			"expected node-01 allocation released, got %d / %d",
			c.node1.AllocatedCPU, c.node1.AllocatedMemory,
		)
	}

	// Registry reflects the new placement.
	stored, err := c.services.Get(service.ID)
	if err != nil {
		t.Fatal(err)
	}
	if stored.NodeID != c.node2.ID {
		t.Fatalf("registry still has service on %s", stored.NodeID)
	}
}

func TestRescheduleNodeServicesMarksFailedWhenNoCapacity(t *testing.T) {
	// node-02 only has 1 CPU, so the 2-CPU service can't move there.
	c := newTestCluster(t, 1)

	service := &models.Service{
		ID:     "service-01",
		Name:   "web",
		Image:  "nginx",
		CPU:    2,
		Memory: 512 * 1024 * 1024,
	}
	c.deployOnNode1(t, service)

	c.node1.Status = models.NodeOffline

	err := c.orch.RescheduleNodeServices(c.node1.ID)
	if !errors.Is(err, scheduler.ErrNoSuitableNode) {
		t.Fatalf("expected ErrNoSuitableNode, got %v", err)
	}

	if service.Status != models.Failed {
		t.Fatalf("expected Failed, got %s", service.Status)
	}
	if service.LastError == "" {
		t.Fatal("expected LastError to explain the failure")
	}
	if c.node2.AllocatedCPU != 0 {
		t.Fatalf("expected nothing allocated on node-02, got %d CPU", c.node2.AllocatedCPU)
	}
}

func TestServicesOnNode(t *testing.T) {
	c := newTestCluster(t, 8)

	running := &models.Service{ID: "a", NodeID: "node-01", Status: models.Running}
	stopped := &models.Service{ID: "b", NodeID: "node-01", Status: models.Stopped}
	elsewhere := &models.Service{ID: "c", NodeID: "node-02", Status: models.Running}

	for _, s := range []*models.Service{running, stopped, elsewhere} {
		if err := c.services.Add(s); err != nil {
			t.Fatal(err)
		}
	}

	got := c.orch.ServicesOnNode("node-01")
	if len(got) != 1 || got[0].ID != "a" {
		t.Fatalf("expected only service a, got %v", got)
	}
}
