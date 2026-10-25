package main

import (
	"flag"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	"vortex-edge/internal/agent"
	"vortex-edge/internal/models"
	"vortex-edge/internal/runtime"
	"vortex-edge/proto"
)

func main() {
	nodeID := flag.String("node-id", "node-01", "stable node identity (not its address)")
	name := flag.String("name", "worker-01", "human-readable node name")
	listen := flag.String("listen", "0.0.0.0:9100", "address the agent's gRPC server binds to")
	advertise := flag.String("advertise", "localhost:9100", "address the control plane uses to reach this agent (e.g. its Tailscale IP)")
	controlPlane := flag.String("control-plane", "localhost:9001", "control plane gRPC address")
	cpu := flag.Int("cpu", 8, "CPU cores offered to the cluster")
	memoryGB := flag.Int64("memory-gb", 16, "memory (GB) offered to the cluster")
	flag.Parse()

	node := &models.Node{
		ID:      *nodeID,
		Name:    *name,
		Address: *advertise,
		Status:  models.NodeReady,
		CPU:     *cpu,
		Memory:  *memoryGB * 1024 * 1024 * 1024,
	}

	dockerRuntime, err := runtime.NewDockerRuntime()
	if err != nil {
		log.Fatal(err)
	}

	nodeAgent := agent.NewAgent(node, dockerRuntime)

	// Serve the NodeAgent API so the control plane can run workloads here.
	listener, err := net.Listen("tcp", *listen)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", *listen, err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterNodeAgentServer(grpcServer, nodeAgent)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("agent gRPC server stopped: %v", err)
		}
	}()

	log.Printf("Agent %s serving on %s (advertised as %s)", node.ID, *listen, *advertise)

	client, err := agent.NewClient(*controlPlane)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	accepted, err := client.Register(node)
	if err != nil {
		log.Fatal(err)
	}

	if !accepted {
		log.Fatal("node registration rejected")
	}

	log.Println("Node registered successfully")

	go nodeAgent.StartHeartbeat(client, 5*time.Second)

	select {}
}
