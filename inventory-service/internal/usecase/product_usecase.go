package usecase

import (
	"context"
	"errors"
	"time"

	"inventoryservice/internal/domain"
	"inventoryservice/internal/domain/model"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidPrice    = errors.New("price must be positive")
	ErrInvalidStock    = errors.New("stock cannot be negative")
)

type ProductUseCase interface {
	CreateProduct(ctx context.Context, product *model.Product) (*model.Product, error)
	GetProduct(ctx context.Context, id string) (*model.Product, error)
	UpdateProduct(ctx context.Context, product *model.Product) error
	DeleteProduct(ctx context.Context, id string) error
	ListProducts(ctx context.Context, skip, limit int64) ([]*model.Product, error)
}

type productUseCase struct {
	repo domain.ProductRepository
}

func NewProductUseCase(repo domain.ProductRepository) ProductUseCase {
	return &productUseCase{repo: repo}
}

func (uc *productUseCase) CreateProduct(ctx context.Context, product *model.Product) (*model.Product, error) {
	if product.Price <= 0 {
		return nil, ErrInvalidPrice
	}
	if product.Stock < 0 {
		return nil, ErrInvalidStock
	}

	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	return uc.repo.Create(ctx, product)
}

func (uc *productUseCase) GetProduct(ctx context.Context, id string) (*model.Product, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *productUseCase) UpdateProduct(ctx context.Context, product *model.Product) error {
	existing, err := uc.repo.FindByID(ctx, product.ID)
	if err != nil {
		return err
	}

	// Preserve created_at
	product.CreatedAt = existing.CreatedAt
	product.UpdatedAt = time.Now()

	return uc.repo.Update(ctx, product)
}

func (uc *productUseCase) DeleteProduct(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *productUseCase) ListProducts(ctx context.Context, skip, limit int64) ([]*model.Product, error) {
	return uc.repo.List(ctx, skip, limit)
}