package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"vortex-edge/internal/cluster"
	"vortex-edge/internal/registery"
	"vortex-edge/proto"
)

func main() {
	registry := registery.NewNodeRegistry()
	server := cluster.NewServer(registry)

	listener, err := net.Listen("tcp", "localhost:9001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterNodeAgentServer(grpcServer, server)

	log.Println("Control plane listening on localhost:9001")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("gRPC server stopped: %v", err)
	}
}
