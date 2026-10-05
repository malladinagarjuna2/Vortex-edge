package agent

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"vortex-edge/internal/models"
	"vortex-edge/proto"
)

// Client communicates with a Node Agent through gRPC.
type Client struct {
	conn   *grpc.ClientConn
	client proto.NodeAgentClient
}

// NewClient creates a new gRPC client.
func NewClient(address string) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:   conn,
		client: proto.NewNodeAgentClient(conn),
	}, nil
}

func (c *Client) Heartbeat(nodeID string) (bool, error) {
	ctx := context.Background()

	response, err := c.client.Heartbeat(
		ctx,
		&proto.HeartbeatRequest{
			NodeId: nodeID,
		},
	)
	if err != nil {
		return false, err
	}

	return response.Acknowledged, nil
}

func (c *Client) Register(node *models.Node) (bool, error) {
	ctx := context.Background()

	response, err := c.client.Register(
		ctx,
		&proto.RegisterRequest{
			NodeId:  node.ID,
			Name:    node.Name,
			Address: node.Address,
			Cpu:     int32(node.CPU),
			Memory:  node.Memory,
		},
	)
	if err != nil {
		return false, err
	}

	return response.Accepted, nil
}

func (c *Client) RunService(
	req *proto.RunServiceRequest,
) (*proto.RunServiceResponse, error) {

	ctx := context.Background()

	return c.client.RunService(ctx, req)
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c*Client)DeleteService( req*proto.DeleteServiceRequest,)(*proto.DeleteServiceResponse, error){
	 ctx:= context.Background()
	 return c.client.DeleteService(ctx,req)
}

func (c *Client) Logs(
    req *proto.LogsRequest,
) (*proto.LogsResponse, error) {
    ctx := context.Background()

    return c.client.Logs(ctx, req)
}