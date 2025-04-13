package domain

import (
	"context"
	"errors"
	"inventoryservice/internal/domain/model"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidID       = errors.New("invalid ID format")
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) (*model.Product, error)
	FindByID(ctx context.Context, id string) (*model.Product, error)
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, skip, limit int64) ([]*model.Product, error)
}