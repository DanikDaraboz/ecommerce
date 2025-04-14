package client

import (
    "context"
    "google.golang.org/grpc"
    "log"
    pborder "github.com/DanikDaraboz/ecommerce/proto/order"
)

type OrderClient struct {
    client pb.OrderServiceClient
}

func NewOrderClient(addr string) (*OrderClient, error) {
    conn, err := grpc.Dial(addr, grpc.WithInsecure())
    if err != nil {
        return nil, err
    }
    return &OrderClient{client: pb.NewOrderServiceClient(conn)}, nil
}

func (c *OrderClient) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
    return c.client.CreateOrder(ctx, req)
}

func (c *OrderClient) GetOrderByID(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {
    return c.client.GetOrderByID(ctx, req)
}

func (c *OrderClient) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.OrderResponse, error) {
    return c.client.UpdateOrderStatus(ctx, req)
}

func (c *OrderClient) ListUserOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
    return c.client.ListUserOrders(ctx, req)
}