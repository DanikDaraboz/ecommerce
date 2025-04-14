package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/danikdaraboz/ecommerce/order-service/internal/domain"
	"github.com/danikdaraboz/ecommerce/order-service/internal/domain/model"
)

type OrderUsecase struct {
	repo domain.OrderRepository
}

func NewOrderUsecase(repo domain.OrderRepository) *OrderUsecase {
	return &OrderUsecase{repo: repo}
}

func (uc *OrderUsecase) CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	if order == nil || len(order.Items) == 0 {
		return nil, domain.ErrInvalidOrder
	}

	var total float64
	for _, item := range order.Items {
		total += float64(item.Quantity) * item.Price
	}

	order.Total = total
	order.Status = model.StatusCreated
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	return uc.repo.Create(ctx, order)
}

func (uc *OrderUsecase) GetOrderByID(ctx context.Context, id string) (*model.Order, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *OrderUsecase) GetOrdersByUserID(ctx context.Context, userID string) ([]*model.Order, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	return uc.repo.FindByUserID(ctx, userID)
}

func (uc *OrderUsecase) UpdateOrderStatus(ctx context.Context, id string, status model.OrderStatus) error {
	return uc.repo.UpdateStatus(ctx, id, status)
}

func (uc *OrderUsecase) DeleteOrder(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
