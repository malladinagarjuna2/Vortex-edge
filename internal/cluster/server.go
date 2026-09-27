
// Node Agent
//     │
//     │ Register(node information)
//     ▼
// Control Plane gRPC Server
//     │
//     ▼
// NodeRegistry
//     │
//     ▼
// "node-01 registered"



// User
//  │
//  │ Deploy workload
//  ▼
// Control Plane
//  │
//  ▼
// Scheduler
//  │
//  ├── Node 1 ❌ insufficient
//  ├── Node 2 ✅
//  └── Node 3 ✅
//  │
//  ▼
// Node 2
//  │
//  ▼
// Agent
//  │
//  ▼
// Docker
//  │
//  ▼
// Container


package cluster

import (
	"context"
	"time"

	"vortex-edge/internal/models"
	"vortex-edge/internal/registery"
	"vortex-edge/proto"
)
//server ek schema hai 
type Server struct {
	proto.UnimplementedNodeAgentServer

	registry *registery.NodeRegistry
}
// ek naya server register kiya hai humne 
func NewServer(registry *registery.NodeRegistry) *Server {
	return &Server{
		registry: registry,
	}
}
// control plane me server register kiya hai humne waha pe model  register kiya hai 
func (s *Server) Register(
	ctx context.Context,
	req *proto.RegisterRequest,
) (*proto.RegisterResponse, error) {

	node := &models.Node{
		ID:      req.NodeId,
		Name:    req.Name,
		Address: req.Address,
		Status:  models.NodeReady,
		CPU:     int(req.Cpu),
		Memory:  req.Memory,
	}

	err := s.registry.Add(node)
	if err != nil {
		return nil, err
	}

	return &proto.RegisterResponse{
		Accepted: true,
	}, nil
}
// control plane kaa heartbeat implement kro 
func (s *Server) Heartbeat(
	ctx context.Context,
	req *proto.HeartbeatRequest,
) (*proto.HeartbeatResponse, error) {

	node, err := s.registry.Get(req.NodeId)
	if err != nil {
		return nil, err
	}

	if node == nil {
		return &proto.HeartbeatResponse{
			Acknowledged: false,
		}, nil
	}

	node.LastHeartbeat = time.Now()
	node.Status = models.NodeReady

	err = s.registry.Update(node)
	if err != nil {
		return nil, err
	}

	return &proto.HeartbeatResponse{
		Acknowledged: true,
	}, nil
}