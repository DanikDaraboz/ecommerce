package service

import (
	"context"
	"log"

	"github.com/danikdaraboz/ecommerce/api-gateway/pb"
	"google.golang.org/grpc"
)

type InventoryClient struct {
	client proto.InventoryServiceClient
}

func NewInventoryClient(address string) *InventoryClient {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to inventory service: %v", err)
	}
	client := proto.NewInventoryServiceClient(conn)
	return &InventoryClient{client: client}
}

func (c *InventoryClient) GetProductByID(ctx context.Context, id string) (*proto.Product, error) {
	req := &proto.GetProductRequest{Id: id}
	return c.client.GetProductByID(ctx, req)
}
