package agent

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"vortex-edge/internal/models"
	"vortex-edge/internal/runtime"
	"vortex-edge/proto"
)

func TestGRPC(t *testing.T) {
	// Create node
	node := &models.Node{
		ID:      "node-01",
		Name:    "worker-01",
		Status:  models.NodeReady,
		CPU:     8,
		Memory:  16 * 1024 * 1024 * 1024,
	}

	// Create Docker runtime
	dockerRuntime, err := runtime.NewDockerRuntime()
	if err != nil {
		t.Fatalf("failed to create Docker runtime: %v", err)
	}

	// Create Agent
	agent := NewAgent(node, dockerRuntime)

	// Create TCP listener
	listener, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()

	// Create gRPC server
	server := grpc.NewServer()

	// Register our Agent
	proto.RegisterNodeAgentServer(server, agent)

	// Start server
	go func() {
		if err := server.Serve(listener); err != nil {
			t.Logf("server stopped: %v", err)
		}
	}()

	// Connect to server
	conn, err := grpc.NewClient(
		"localhost:9000",
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// Create generated gRPC client
	client := proto.NewNodeAgentClient(conn)

	// Call Heartbeat
	response, err := client.Heartbeat(
		context.Background(),
		&proto.HeartbeatRequest{
			NodeId: "node-01",
		},
	)
	if err != nil {
		t.Fatalf("heartbeat failed: %v", err)
	}

	// Verify response
	if !response.Acknowledged {
		t.Fatal("expected heartbeat to be acknowledged")
	}

	t.Log("gRPC Heartbeat successful")
}