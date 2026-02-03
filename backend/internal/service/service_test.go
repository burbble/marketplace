package service_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/google/uuid"

	"github.com/burbble/marketplace/internal/domain"
	"github.com/burbble/marketplace/internal/mocks"
	"github.com/burbble/marketplace/internal/service"
)

func TestProductService_GetByID_Success(t *testing.T) {
	id := uuid.New()
	repo := &mocks.ProductRepositoryMock{
		GetByIDFunc: func(_ context.Context, gotID uuid.UUID) (*domain.Product, error) {
			if gotID != id {
				t.Errorf("expected id %s, got %s", id, gotID)
			}
			return &domain.Product{ID: id, Name: "Phone"}, nil
		},
	}

	svc := service.NewProductService(repo)
	p, err := svc.GetByID(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "Phone" {
		t.Errorf("expected name 'Phone', got %q", p.Name)
	}
	if len(repo.GetByIDCalls()) != 1 {
		t.Errorf("expected 1 call to GetByID, got %d", len(repo.GetByIDCalls()))
	}
}

func TestProductService_GetByID_NotFound(t *testing.T) {
	repo := &mocks.ProductRepositoryMock{
		GetByIDFunc: func(_ context.Context, _ uuid.UUID) (*domain.Product, error) {
			return nil, sql.ErrNoRows
		},
	}

	svc := service.NewProductService(repo)
	_, err := svc.GetByID(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestProductService_GetByFilter(t *testing.T) {
	repo := &mocks.ProductRepositoryMock{
		GetByFilterFunc: func(_ context.Context, f domain.ProductFilter) (*domain.ProductList, error) {
			if f.Limit != 10 {
				t.Errorf("expected limit 10, got %d", f.Limit)
			}
			return &domain.ProductList{
				Products: []domain.Product{{Name: "A"}, {Name: "B"}},
				Total:    2,
				Page:     1,
				PageSize: 10,
			}, nil
		},
	}

	svc := service.NewProductService(repo)
	result, err := svc.GetByFilter(context.Background(), domain.ProductFilter{Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}
	if len(repo.GetByFilterCalls()) != 1 {
		t.Errorf("expected 1 call to GetByFilter, got %d", len(repo.GetByFilterCalls()))
	}
}

func TestProductService_GetByFilter_Error(t *testing.T) {
	repo := &mocks.ProductRepositoryMock{
		GetByFilterFunc: func(_ context.Context, _ domain.ProductFilter) (*domain.ProductList, error) {
			return nil, fmt.Errorf("db error")
		},
	}

	svc := service.NewProductService(repo)
	_, err := svc.GetByFilter(context.Background(), domain.ProductFilter{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestProductService_GetBrands(t *testing.T) {
	repo := &mocks.ProductRepositoryMock{
		GetBrandsFunc: func(_ context.Context) ([]string, error) {
			return []string{"Apple", "Samsung"}, nil
		},
	}

	svc := service.NewProductService(repo)
	brands, err := svc.GetBrands(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(brands) != 2 {
		t.Errorf("expected 2 brands, got %d", len(brands))
	}
	if len(repo.GetBrandsCalls()) != 1 {
		t.Errorf("expected 1 call to GetBrands, got %d", len(repo.GetBrandsCalls()))
	}
}

func TestProductService_GetBrands_Error(t *testing.T) {
	repo := &mocks.ProductRepositoryMock{
		GetBrandsFunc: func(_ context.Context) ([]string, error) {
			return nil, fmt.Errorf("db error")
		},
	}

	svc := service.NewProductService(repo)
	_, err := svc.GetBrands(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}


func TestCategoryService_GetAll_Success(t *testing.T) {
	repo := &mocks.CategoryRepositoryMock{
		GetAllFunc: func(_ context.Context) ([]domain.Category, error) {
			return []domain.Category{
				{Name: "Phones"},
				{Name: "Laptops"},
			}, nil
		},
	}

	svc := service.NewCategoryService(repo)
	cats, err := svc.GetAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cats) != 2 {
		t.Errorf("expected 2 categories, got %d", len(cats))
	}
	if len(repo.GetAllCalls()) != 1 {
		t.Errorf("expected 1 call to GetAll, got %d", len(repo.GetAllCalls()))
	}
}

func TestCategoryService_GetAll_Error(t *testing.T) {
	repo := &mocks.CategoryRepositoryMock{
		GetAllFunc: func(_ context.Context) ([]domain.Category, error) {
			return nil, fmt.Errorf("db error")
		},
	}

	svc := service.NewCategoryService(repo)
	_, err := svc.GetAll(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCategoryService_GetByID_Success(t *testing.T) {
	id := uuid.New()
	repo := &mocks.CategoryRepositoryMock{
		GetByIDFunc: func(_ context.Context, gotID uuid.UUID) (*domain.Category, error) {
			return &domain.Category{ID: gotID, Name: "Phones"}, nil
		},
	}

	svc := service.NewCategoryService(repo)
	cat, err := svc.GetByID(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cat.ID != id {
		t.Errorf("expected id %s, got %s", id, cat.ID)
	}
	if len(repo.GetByIDCalls()) != 1 {
		t.Errorf("expected 1 call to GetByID, got %d", len(repo.GetByIDCalls()))
	}
}

func TestCategoryService_GetByID_NotFound(t *testing.T) {
	repo := &mocks.CategoryRepositoryMock{
		GetByIDFunc: func(_ context.Context, _ uuid.UUID) (*domain.Category, error) {
			return nil, sql.ErrNoRows
		},
	}

	svc := service.NewCategoryService(repo)
	_, err := svc.GetByID(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
