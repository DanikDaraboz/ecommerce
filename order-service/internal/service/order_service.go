package service

import (
    "context"
    "order-service/internal/domain"
    "github.com/google/uuid"
)

type OrderRepository interface {
    Create(ctx context.Context, order *domain.Order) error
    GetByID(ctx context.Context, id string) (*domain.Order, error)
    UpdateStatus(ctx context.Context, id, status string) error
    ListByUser(ctx context.Context, userID string, page, pageSize int32) ([]*domain.Order, error)
}

type OrderService struct {
    repo OrderRepository
}

func NewOrderService(repo OrderRepository) *OrderService {
    return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID string, items []domain.OrderItem) (*domain.Order, error) {
    // Simplified total calculation (in real app, fetch product prices)
    total := 0.0
    for _, item := range items {
        total += float64(item.Quantity) * 10.0 // Mock price
    }

    order := &domain.Order{
        ID:     uuid.New().String(),
        UserID: userID,
        Items:  items,
        Status: "pending",
        Total:  total,
    }
    if err := s.repo.Create(ctx, order); err != nil {
        return nil, err
    }
    return order, nil
}

func (s *OrderService) GetOrderByID(ctx context.Context, id string) (*domain.Order, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, id, status string) (*domain.Order, error) {
    if err := s.repo.UpdateStatus(ctx, id, status); err != nil {
        return nil, err
    }
    return s.repo.GetByID(ctx, id)
}

func (s *OrderService) ListUserOrders(ctx context.Context, userID string, page, pageSize int32) ([]*domain.Order, error) {
    return s.repo.ListByUser(ctx, userID, page, pageSize)
}