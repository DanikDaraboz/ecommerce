package client

import (
    "context"
    "google.golang.org/grpc"
    "log"
    pb "ecommerce-platform/proto/inventory"
)

type InventoryClient struct {
    client pb.InventoryServiceClient
}

func NewInventoryClient(addr string) (*InventoryClient, error) {
    conn, err := grpc.Dial(addr, grpc.WithInsecure())
    if err != nil {
        return nil, err
    }
    return &InventoryClient{client: pb.NewInventoryServiceClient(conn)}, nil
}

func (c *InventoryClient) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
    return c.client.CreateProduct(ctx, req)
}

func (c *InventoryClient) GetProductByID(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
    return c.client.GetProductByID(ctx, req)
}

func (c *InventoryClient) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
    return c.client.UpdateProduct(ctx, req)
}

func (c *InventoryClient) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.Empty, error) {
    return c.client.DeleteProduct(ctx, req)
}

func (c *InventoryClient) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
    return c.client.ListProducts(ctx, req)
}