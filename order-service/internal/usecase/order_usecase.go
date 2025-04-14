package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/danikdaraboz/ecommerce/order-service/internal/domain"
	"github.com/danikdaraboz/ecommerce/order-service/internal/domain/model"
)

// Интерфейс для использования в handler и других слоях
type OrderUseCase interface {
	CreateOrder(ctx context.Context, userID string, items []model.OrderItem) (*model.Order, error)
	GetOrderByID(ctx context.Context, id string) (*model.Order, error)
	GetOrdersByUserID(ctx context.Context, userID string) ([]*model.Order, error)
	UpdateOrderStatus(ctx context.Context, id string, status model.OrderStatus) error
	DeleteOrder(ctx context.Context, id string) error
}

// Структура, реализующая интерфейс
type OrderUsecase struct {
	repo domain.OrderRepository
}

func NewOrderUseCase(repo domain.OrderRepository) *OrderUsecase {
	return &OrderUsecase{repo: repo} // создаем объект типа OrderUsecase
}

// Обновлённый метод CreateOrder с userID и items
func (uc *OrderUsecase) CreateOrder(ctx context.Context, userID string, items []model.OrderItem) (*model.Order, error) {
	if userID == "" || len(items) == 0 {
		return nil, domain.ErrInvalidOrder
	}

	var total float64
	for _, item := range items {
		total += float64(item.Quantity) * item.Price
	}

	order := &model.Order{
		UserID:    userID,
		Items:     items,
		Total:     total,
		Status:    model.StatusCreated,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

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

// Проверка на реализацию интерфейса
var _ OrderUseCase = (*OrderUsecase)(nil)
