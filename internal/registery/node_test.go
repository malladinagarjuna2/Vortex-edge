
package registery 

import(
	"testing"
	"vortex-edge/internal/models"
)

func TestNodeRegistry( t*testing.T){
	 registry := NewNodeRegistry()
	
	 node := &models.Node{
		 ID: "node-01",
		 Name:"worker-01",
		 Address:"192.168.1.10:9000", 
		 Status: models.NodeReady, 
		 CPU: 8, 
		 Memory: 16*1024*1024*1024, 


	 }
	 err:= registry.Add(node)
	if err!= nil{
		 t.Fatalf("failed to add node: %v", err)
	}
	 
	got,err:= registry.Get("node-01")
	if err!=nil {
		 t.Fatalf("failed to get node:%v", err)
	}
	 
   	if got == nil {
		t.Fatal("expected node, got nil")
	}

	if got.Name != "worker-01" {
		t.Fatalf("expected worker-01, got %s", got.Name)
	}

	nodes := registry.List()

	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(nodes))
	}
}