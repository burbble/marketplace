package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/burbble/marketplace/internal/domain"
	"github.com/burbble/marketplace/internal/repository/postgres"
)

type ProductService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	GetByFilter(ctx context.Context, filter domain.ProductFilter) (*domain.ProductList, error)
	GetBrands(ctx context.Context) ([]string, error)
}

type productService struct {
	repo postgres.ProductRepository
}

func NewProductService(repo postgres.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *productService) GetByFilter(ctx context.Context, filter domain.ProductFilter) (*domain.ProductList, error) {
	return s.repo.GetByFilter(ctx, filter)
}

func (s *productService) GetBrands(ctx context.Context) ([]string, error) {
	return s.repo.GetBrands(ctx)
}
