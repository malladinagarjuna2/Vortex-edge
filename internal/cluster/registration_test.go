// Node Agent / Client
//         │
//         │ Register()
//         ▼
//    gRPC network
//         │
//         ▼
// Control Plane Server
//         │
//         ▼
//    Node Registry
//         │
//         ▼
//      node-01

package cluster

import (
	"context"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"vortex-edge/internal/models"
	"vortex-edge/internal/registery"
	"vortex-edge/proto"
)

func TestNodeRegistration(t *testing.T) {
	
	//apni registry banaao 
	registry := registery.NewNodeRegistry()

	//control plane  kaa grpc server hai yeh 
	server := NewServer(registry)
// tcp listerner hai yeh jo regsitry se aayi hui cheezo ko syun sakta hai 
	listener, err := net.Listen("tcp", "localhost:9001")
	if err != nil {
		t.Fatalf("failed to create listener: %v", err)
	}
	defer listener.Close()

	//grpc server banana hai
	grpcServer := grpc.NewServer()

//ek control plane ko register kiya hai 
	proto.RegisterNodeAgentServer(grpcServer, server)

//grpc server ko start kiya
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			t.Logf("gRPC server stopped: %v", err)
		}
	}()

	//server ko start krne ke liye time de rahe hai 
	time.Sleep(100 * time.Millisecond)
//node side se grpc connection generate kr rahe hai 
	conn, err := grpc.NewClient(
		"localhost:9001",
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		t.Fatalf("failed to connect to Control Plane: %v", err)
	}
	defer conn.Close()

	// ek grpc client create hogaya hai 
	client := proto.NewNodeAgentClient(conn)

	// ek node ko register join kiya hai 
	node := &models.Node{
		ID:      "node-01",
		Name:    "worker-01",
		Address: "localhost:9000",
		Status:  models.NodeReady,
		CPU:     8,
		Memory:  16 * 1024 * 1024 * 1024,
	}

	//yaha apna node register kr rahe hai
	response, err := client.Register(
		context.Background(),
		&proto.RegisterRequest{
			NodeId:  node.ID,
			Name:    node.Name,
			Address: node.Address,
			Cpu:     int32(node.CPU),
			Memory:  node.Memory,
		},
	)
	if err != nil {
		t.Fatalf("registration request failed: %v", err)
	}

	if !response.Accepted {
		t.Fatal("expected registration to be accepted")
	}

//yeh verify krna hai kee node plane store kr skta hai kee nahi 
	registeredNode, err := registry.Get("node-01")
	if err != nil {
		t.Fatalf("failed to get registered node: %v", err)
	}

	if registeredNode == nil {
		t.Fatal("node was not added to registry")
	}

	if registeredNode.Name != "worker-01" {
		t.Fatalf(
			"expected worker-01, got %s",
			registeredNode.Name,
		)
	}

	if registeredNode.CPU != 8 {
		t.Fatalf(
			"expected 8 CPU, got %d",
			registeredNode.CPU,
		)
	}

	t.Log("Node successfully registered through gRPC")
}

