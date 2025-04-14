package server

import (
    "context"
    "inventory-service/internal/service"
    pb "inventory-service/proto"
    "net"
    "google.golang.org/grpc"
)

type InventoryServer struct {
    pb.UnimplementedInventoryServiceServer
    svc *service.ProductService
}

func NewInventoryServer(svc *service.ProductService) *InventoryServer {
    return &InventoryServer{svc: svc}
}

func (s *InventoryServer) Run(addr string) error {
    lis, err := net.Listen("tcp", addr)
    if err != nil {
        return err
    }
    grpcServer := grpc.NewServer()
    pb.RegisterInventoryServiceServer(grpcServer, s)
    return grpcServer.Serve(lis)
}

func (s *InventoryServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
    product, err := s.svc.CreateProduct(ctx, req.Name, req.Description, req.Category, req.Price, req.Stock)
    if err != nil {
        return nil, err
    }
    return &pb.ProductResponse{
        Id:          product.ID,
        Name:        product.Name,
        Description: product.Description,
        Price:       product.Price,
        Stock:       product.Stock,
        Category:    product.Category,
    }, nil
}

func (s *InventoryServer) GetProductByID(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
    product, err := s.svc.GetProductByID(ctx, req.Id)
    if err != nil {
        return nil, err
    }
    return &pb.ProductResponse{
        Id:          product.ID,
        Name:        product.Name,
        Description: product.Description,
        Price:       product.Price,
        Stock:       product.Stock,
        Category:    product.Category,
    }, nil
}

func (s *InventoryServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
    product, err := s.svc.UpdateProduct(ctx, req.Id, req.Name, req.Description, req.Category, req.Price, req.Stock)
    if err != nil {
        return nil, err
    }
    return &pb.ProductResponse{
        Id:          product.ID,
        Name:        product.Name,
        Description: product.Description,
        Price:       product.Price,
        Stock:       product.Stock,
        Category:    product.Category,
    }, nil
}

func (s *InventoryServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.Empty, error) {
    if err := s.svc.DeleteProduct(ctx, req.Id); err != nil {
        return nil, err
    }
    return &pb.Empty{}, nil
}

func (s *InventoryServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
    products, err := s.svc.ListProducts(ctx, req.Page, req.PageSize)
    if err != nil {
        return nil, err
    }
    resp := &pb.ListProductsResponse{}
    for _, p := range products {
        resp.Products = append(resp.Products, &pb.ProductResponse{
            Id:          p.ID,
            Name:        p.Name,
            Description: p.Description,
            Price:       p.Price,
            Stock:       p.Stock,
            Category:    p.Category,
        })
    }
    return resp, nil
}