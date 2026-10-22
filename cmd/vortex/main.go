package main

import (
	"context"
	"flag"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	"vortex-edge/internal/cluster"
	"vortex-edge/internal/orchestrator"
	"vortex-edge/internal/recovery"
	"vortex-edge/internal/registery"
	"vortex-edge/internal/scheduler"
	"vortex-edge/proto"
)

const (
	// A node is Offline after missing heartbeats for this long
	// (agents send one every 5s).
	heartbeatTimeout = 15 * time.Second

	// How often the control plane checks for failed nodes.
	recoveryInterval = 5 * time.Second
)

func main() {
	listen := flag.String("listen", "localhost:9001", "control plane gRPC address (use the Tailscale IP in a real cluster)")
	flag.Parse()

	registry := registery.NewNodeRegistry()
	services := registery.NewMemoryRegistery()

	clusterState := cluster.NewClusterState(registry)
	orch := orchestrator.NewOrchestrator(
		scheduler.NewScheduler(clusterState),
		clusterState,
		services,
	)

	// Failure recovery: Offline node → reschedule its services.
	recoveryManager := recovery.NewManager(
		cluster.NewHealthMonitor(registry, heartbeatTimeout),
		orch,
	)
	go recoveryManager.Run(context.Background(), recoveryInterval)

	server := cluster.NewServer(registry)

	listener, err := net.Listen("tcp", *listen)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterControlPlaneServer(grpcServer, server)

	log.Printf("Control plane listening on %s", *listen)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("gRPC server stopped: %v", err)
	}
}
