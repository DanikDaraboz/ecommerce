package server

import (
    "context"
    "order-service/internal/service"
    pb "order-service/proto"
    "net"
    "google.golang.org/grpc"
)

type OrderServer struct {
    pb.UnimplementedOrderServiceServer
    svc *service.OrderService
}

func NewOrderServer(svc *service.OrderService) *OrderServer {
    return &OrderServer{svc: svc}
}

func (s *OrderServer) Run(addr string) error {
    lis, err := net.Listen("tcp", addr)
    if err != nil {
        return err
    }
    grpcServer := grpc.NewServer()
    pb.RegisterOrderServiceServer(grpcServer, s)
    return grpcServer.Serve(lis)
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
    items := make([]service.OrderItem, len(req.Items))
    for i, item := range req.Items {
        items[i] = service.OrderItem{
            ProductID: item.ProductId,
            Quantity:  item.Quantity,
        }
    }

    order, err := s.svc.CreateOrder(ctx, req.UserId, items)
    if err != nil {
        return nil, err
    }

    resp := &pb.OrderResponse{
        Id:     order.ID,
        UserId: order.UserID,
        Status: order.Status,
        Total:  order.Total,
    }
    for _, item := range order.Items {
        resp.Items = append(resp.Items, &pb.OrderItem{
            ProductId: item.ProductID,
            Quantity:  item.Quantity,
        })
    }
    return resp, nil
}

func (s *OrderServer) GetOrderByID(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {
    order, err := s.svc.GetOrderByID(ctx, req.Id)
    if err != nil {
        return nil, err
    }

    resp := &pb.OrderResponse{
        Id:     order.ID,
        UserId: order.UserID,
        Status: order.Status,
        Total:  order.Total,
    }
    for _, item := range order.Items {
        resp.Items = append(resp.Items, &pb.OrderItem{
            ProductId: item.ProductID,
            Quantity:  item.Quantity,
        })
    }
    return resp, nil
}

func (s *OrderServer) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.OrderResponse, error) {
    order, err := s.svc.UpdateOrderStatus(ctx, req.Id, req.Status)
    if err != nil {
        return nil, err
    }

    resp := &pb.OrderResponse{
        Id:     order.ID,
        UserId: order.UserID,
        Status: order.Status,
        Total:  order.Total,
    }
    for _, item := range order.Items {
        resp.Items = append(resp.Items, &pb.OrderItem{
            ProductId: item.ProductID,
            Quantity:  item.Quantity,
        })
    }
    return resp, nil
}

func (s *OrderServer) ListUserOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
    orders, err := s.svc.ListUserOrders(ctx, req.UserId, req.Page, req.PageSize)
    if err != nil {
        return nil, err
    }

    resp := &pb.ListOrdersResponse{}
    for _, order := range orders {
        orderResp := &pb.OrderResponse{
            Id:     order.ID,
            UserId: order.UserID,
            Status: order.Status,
            Total:  order.Total,
        }
        for _, item := range order.Items {
            orderResp.Items = append(orderResp.Items, &pb.OrderItem{
                ProductId: item.ProductID,
                Quantity:  item.Quantity,
            })
        }
        resp.Orders = append(resp.Orders, orderResp)
    }
    return resp, nil
}