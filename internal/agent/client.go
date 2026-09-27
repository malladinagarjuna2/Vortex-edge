package agent

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"vortex-edge/internal/models"
	"vortex-edge/proto"
)
//client architecture grpc server ke thru baat cheet krega
type Client struct {
	conn   *grpc.ClientConn
	client proto.NodeAgentClient
}
// naya client banane ke liye grpc server ek naya client register
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

func (c *Client) Close() error {
	return c.conn.Close()
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
