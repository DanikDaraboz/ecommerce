package grpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/danikdaraboz/ecommerce/order-service/internal/domain/model"
	"github.com/danikdaraboz/ecommerce/order-service/internal/usecase"
	"github.com/danikdaraboz/ecommerce/order-service/pb"
)

type OrderHandler struct {
	orderUC usecase.OrderUseCase
	pb.UnimplementedOrderServiceServer
}

func NewOrderHandler(orderUC usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{orderUC: orderUC}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	if req.GetUserId() == "" || len(req.GetItems()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var items []model.OrderItem
	for _, item := range req.GetItems() {
		items = append(items, model.OrderItem{
			ProductID: item.GetProductId(),
			Quantity:  int(item.GetQuantity()),
			Price:     item.GetPrice(),
		})
	}

	order, err := h.orderUC.CreateOrder(ctx, req.GetUserId(), items)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.OrderResponse{
		Order: h.toProtoOrder(order),
	}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "order ID is required")
	}

	order, err := h.orderUC.GetOrder(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return nil, status.Error(codes.NotFound, "order not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.OrderResponse{
		Order: h.toProtoOrder(order),
	}, nil
}

// Implement other gRPC methods (GetUserOrders, UpdateOrderStatus, CancelOrder)

func (h *OrderHandler) toProtoOrder(order *model.Order) *pb.Order {
	var items []*pb.OrderItem
	for _, item := range order.Items {
		items = append(items, &pb.OrderItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
			Price:     item.Price,
		})
	}

	return &pb.Order{
		Id:        order.ID,
		UserId:    order.UserID,
		Items:     items,
		Total:     order.Total,
		Status:    string(order.Status),
		CreatedAt: order.CreatedAt.Format(time.RFC3339),
	}
}