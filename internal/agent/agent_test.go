
package agent

import (
	"testing"  //which package?
	
	"vortex-edge/internal/runtime"
	"vortex-edge/internal/models"
)


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