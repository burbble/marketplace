package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/burbble/marketplace/internal/domain"
	"github.com/burbble/marketplace/internal/repository/postgres"
)

type CategoryService interface {
	GetAll(ctx context.Context) ([]domain.Category, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error)
}

type categoryService struct {
	repo postgres.CategoryRepository
}

func NewCategoryService(repo postgres.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetAll(ctx context.Context) ([]domain.Category, error) {
	return s.repo.GetAll(ctx)
}

func (s *categoryService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	return s.repo.GetByID(ctx, id)
}
