package grpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"inventoryservice/internal/usecase"
	"inventoryservice/internal/domain/model"
	"inventoryservice/pb"
)

type InventoryHandler struct {
	productUC usecase.ProductUseCase
	pb.UnimplementedInventoryServiceServer
}

func NewInventoryHandler(productUC usecase.ProductUseCase) *InventoryHandler {
	return &InventoryHandler{productUC: productUC}
}

func (h *InventoryHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	product := &model.Product{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		Stock:       int(req.GetStock()),
		Category:    req.GetCategory(),
	}

	created, err := h.productUC.CreateProduct(ctx, product)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ProductResponse{
		Product: &pb.Product{
			Id:          created.ID,
			Name:        created.Name,
			Description: created.Description,
			Price:       created.Price,
			Stock:       int32(created.Stock),
			Category:    created.Category,
		},
	}, nil
}

func (h *InventoryHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	product, err := h.productUC.GetProduct(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, usecase.ErrProductNotFound) {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ProductResponse{
		Product: &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       int32(product.Stock),
			Category:    product.Category,
		},
	}, nil
}

func (h *InventoryHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
	product := &model.Product{
		ID:          req.GetId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		Stock:       int(req.GetStock()),
		Category:    req.GetCategory(),
	}

	err := h.productUC.UpdateProduct(ctx, product)
	if err != nil {
		if errors.Is(err, usecase.ErrProductNotFound) {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	updated, err := h.productUC.GetProduct(ctx, product.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ProductResponse{
		Product: &pb.Product{
			Id:          updated.ID,
			Name:        updated.Name,
			Description: updated.Description,
			Price:       updated.Price,
			Stock:       int32(updated.Stock),
			Category:    updated.Category,
		},
	}, nil
}

func (h *InventoryHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.Empty, error) {
	err := h.productUC.DeleteProduct(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, usecase.ErrProductNotFound) {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.Empty{}, nil
}

func (h *InventoryHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, err := h.productUC.ListProducts(ctx, int64(req.GetSkip()), int64(req.GetLimit()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var pbProducts []*pb.Product
	for _, p := range products {
		pbProducts = append(pbProducts, &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       int32(p.Stock),
			Category:    p.Category,
		})
	}

	return &pb.ListProductsResponse{
		Products: pbProducts,
		Total:    int32(len(pbProducts)),
	}, nil
}