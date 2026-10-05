
package agent

import (
	"testing"  //which package?
	"context"
	"vortex-edge/internal/runtime"
	"vortex-edge/internal/models"
"vortex-edge/proto"
)

///
func TestNewAgent(t *testing.T){
	 node := &models.Node{ 
		ID: "node-01",
		Name: "worker-01",
	Address: "192.168.1.10:9000",
		Status:  models.NodeReady,
		CPU:     8,
		Memory:  16 * 1024 * 1024 * 1024,
	 }

	 dockerRuntime, err:= runtime.NewDockerRuntime()
	 if err!= nil{
		t.Fatalf("failed to create Docker runtime:%v", err)
	 }
	 agent:= NewAgent(node, dockerRuntime)
    
	 if agent == nil {
		t.Fatal("expected agent, got nil")
	}

	if agent.Node().ID != "node-01" {
		t.Fatalf("expected node-01, got %s", agent.Node().ID)
	}
	
}
//we need to add an agent deployment test  
func TestRunService(t*testing.T){
	 node:= &models.Node{
		ID:      "node-01",
		Name:    "worker-01",
		Status:  models.NodeReady,
		CPU:     4,
		Memory:  8 * 1024 * 1024 * 1024,
	 }
	 dockerRuntime, err:= runtime.NewDockerRuntime()
	 	if err != nil {
		t.Fatalf("failed to create docker runtime: %v", err)
	}
	agent := NewAgent(node, dockerRuntime)

	response, err := agent.RunService(
		context.Background(),
		&proto.RunServiceRequest{
			ServiceId: "service-01",
			Name:      "vortex-test",
			Image:     "hello-world",
			Cpu:       1,
			Memory:    256 * 1024 * 1024,
		},
	)

	if err != nil {
		t.Fatalf("failed to run service: %v", err)
	}

	if !response.Started {
		t.Fatal("expected service to start")
	}

	if response.ContainerId == "" {
		t.Fatal("expected container ID")
	}

	t.Logf("Container started: %s", response.ContainerId)
	
}

//so why we are using [protoc files here the protoc files is uesd here for serailising and deserialising the data 
// serializtion:converting the data into the  given format for storage and transformation
//deserialisation:making the data from the given format
// these proto buffers make the data in such a way that it is ready for the transmission through the 
// rfc buffere channnels. thus it is really important to understand about it 
//While protoc is powerful, Protocol Buffers have some limitations:

// They are not ideal for extremely large datasets that cannot fit into memory.

// Messages are not inherently compressed, though external compression can be applied.

// They are less efficient for scientific data involving large, multi-dimensional arrays.

// In summary, protoc is a critical tool for working with Protocol Buffers, enabling developers to define, serialize, and deserialize structured data efficiently across multiple languages and platforms.