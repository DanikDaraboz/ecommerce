package domain

import (
	"context"
	"github.com/danikdaraboz/ecommerce/order-service/internal/domain/model"
)

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrInvalidOrder  = errors.New("invalid order")
)

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) (*model.Order, error)
	FindByID(ctx context.Context, id string) (*model.Order, error)
	FindByUserID(ctx context.Context, userID string) ([]*model.Order, error)
	UpdateStatus(ctx context.Context, id string, status model.OrderStatus) error
	Delete(ctx context.Context, id string) error
}