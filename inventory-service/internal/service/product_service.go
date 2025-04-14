package service

import (
    "context"
    "inventory-service/internal/domain"
    "github.com/google/uuid"
)

type ProductRepository interface {
    Create(ctx context.Context, product *domain.Product) error
    GetByID(ctx context.Context, id string) (*domain.Product, error
    Update(ctx context.Context, product *domain.Product) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, page, pageSize int32) ([]*domain.Product, error)
}

type ProductService struct {
    repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
    return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(ctx context.Context, name, description, category string, price float64, stock int32) (*domain.Product, error) {
    product := &domain.Product{
        ID:          uuid.New().String(),
        Name:        name,
        Description: description,
        Price:       price,
        Stock:       stock,
        Category:    category,
    }
    if err := s.repo.Create(ctx, product); err != nil {
        return nil, err
    }
    return product, nil
}

func (s *ProductService) GetProductByID(ctx context.Context, id string) (*domain.Product, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *ProductService) UpdateProduct(ctx context.Context, id, name, description, category string, price float64, stock int32) (*domain.Product, error) {
    product := &domain.Product{
        ID:          id,
        Name:        name,
        Description: description,
        Price:       price,
        Stock:       stock,
        Category:    category,
    }
    if err := s.repo.Update(ctx, product); err != nil {
        return nil, err
    }
    return product, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
    return s.repo.Delete(ctx, id)
}

func (s *ProductService) ListProducts(ctx context.Context, page, pageSize int32) ([]*domain.Product, error) {
    return s.repo.List(ctx, page, pageSize)
}