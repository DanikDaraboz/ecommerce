package service

import (
	"context"
	"log"

	"github.com/danikdaraboz/ecommerce/api-gateway/proto"
	"google.golang.org/grpc"
)

type OrderClient struct {
	client proto.OrderServiceClient
}

func NewOrderClient(address string) *OrderClient {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to order service: %v", err)
	}
	client := proto.NewOrderServiceClient(conn)
	return &OrderClient{client: client}
}

func (c *OrderClient) CreateOrder(ctx context.Context, userID string, items []*proto.OrderItem) (*proto.Order, error) {
	req := &proto.CreateOrderRequest{UserId: userID, Items: items}
	return c.client.CreateOrder(ctx, req)
}
