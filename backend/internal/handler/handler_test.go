package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/burbble/marketplace/internal/domain"
	"github.com/burbble/marketplace/internal/mocks"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestProductHandler_List_Success(t *testing.T) {
	svc := &mocks.ProductServiceMock{
		GetByFilterFunc: func(_ context.Context, _ domain.ProductFilter) (*domain.ProductList, error) {
			return &domain.ProductList{
				Products: []domain.Product{{Name: "Phone"}},
				Total:    1,
				Page:     1,
				PageSize: 24,
			}, nil
		},
	}

	h := NewProductHandler(svc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/products?page=1&page_size=24", nil)

	h.List(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var result domain.ProductList
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if result.Total != 1 {
		t.Errorf("expected total 1, got %d", result.Total)
	}
	if len(svc.GetByFilterCalls()) != 1 {
		t.Errorf("expected 1 call to GetByFilter, got %d", len(svc.GetByFilterCalls()))
	}
}

func TestProductHandler_List_InvalidSortField(t *testing.T) {
	svc := &mocks.ProductServiceMock{}
	h := NewProductHandler(svc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/products?sort_fields=invalid_field:asc", nil)

	h.List(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestProductHandler_List_InvalidCategoryID(t *testing.T) {
	svc := &mocks.ProductServiceMock{}
	h := NewProductHandler(svc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/products?category_id=not-a-uuid", nil)

	h.List(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestProductHandler_List_ServiceError(t *testing.T) {
	svc := &mocks.ProductServiceMock{
		GetByFilterFunc: func(_ context.Context, _ domain.ProductFilter) (*domain.ProductList, error) {
			return nil, fmt.Errorf("db error")
		},
	}

	h := NewProductHandler(svc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/products", nil)

	h.List(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestProductHandler_List_WithFilters(t *testing.T) {
	var capturedFilter domain.ProductFilter
	svc := &mocks.ProductServiceMock{
		GetByFilterFunc: func(_ context.Context, f domain.ProductFilter) (*domain.ProductList, error) {
			capturedFilter = f
			return &domain.ProductList{Products: []domain.Product{}, Total: 0, Page: 1, PageSize: 24}, nil
		},
	}

	h := NewProductHandler(svc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/products?brand=Apple&search=iphone&sort_fields=price:asc", nil)

	h.List(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if capturedFilter.Brand == nil || *capturedFilter.Brand != "Apple" {
		t.Errorf("expected brand filter 'Apple'")
	}
	if capturedFilter.Search == nil || *capturedFilter.Search != "iphone" {
		t.Errorf("expected search filter 'iphone'")
	}
}

func TestProductHandler_GetByID_Success(t *testing.T) {
	id := uuid.New()
	svc := &mocks.ProductServiceMock{
		GetByIDFunc: func(_ context.Context, gotID uuid.UUID) (*domain.Product, error) {
			if gotID != id {
				t.Errorf("expected id %s, got %s", id, gotID)
			}
			return &domain.Product{ID: id, Name: "Test Phone"}, nil
		},
	}

	h := NewProductHandler(svc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/products/"+id.String(), nil)
	c.Params = gin.Params{{Key: "id", Value: id.String()}}

	h.GetByID(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if len(svc.GetByIDCalls()) != 1 {
		t.Errorf("expected 1 call to GetByID, got %d", len(svc.GetByIDCalls()))
	}
}

func TestProductHandler_GetByID_InvalidUUID(t *testing.T) {
	svc := &mocks.ProductServiceMock{}
	h := NewProductHandler(svc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/products/bad-id", nil)
	c.Params = gin.Params{{Key: "id", Value: "bad-id"}}

	h.GetByID(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestProductHandler_GetByID_NotFound(t *testing.T) {
	svc := &mocks.ProductServiceMock{
		GetByIDFunc: func(_ context.Context, _ uuid.UUID) (*domain.Product, error) {
			return nil, sql.ErrNoRows
		},
	}

	h := NewProductHandler(svc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	id := uuid.New()
	c.Request = httptest.NewRequest(http.MethodGet, "/products/"+id.String(), nil)
	c.Params = gin.Params{{Key: "id", Value: id.String()}}

	h.GetByID(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestProductHandler_GetByID_InternalError(t *testing.T) {
	svc := &mocks.ProductServiceMock{
		GetByIDFunc: func(_ context.Context, _ uuid.UUID) (*domain.Product, error) {
			return nil, fmt.Errorf("unexpected error")
		},
	}

	h := NewProductHandler(svc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	id := uuid.New()
	c.Request = httptest.NewRequest(http.MethodGet, "/products/"+id.String(), nil)
	c.Params = gin.Params{{Key: "id", Value: id.String()}}

	h.GetByID(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestProductHandler_GetBrands_Success(t *testing.T) {
	svc := &mocks.ProductServiceMock{
		GetBrandsFunc: func(_ context.Context) ([]string, error) {
			return []string{"Apple", "Samsung", "Xiaomi"}, nil
		},
	}

	h := NewProductHandler(svc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/brands", nil)

	h.GetBrands(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var brands []string
	if err := json.Unmarshal(w.Body.Bytes(), &brands); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(brands) != 3 {
		t.Errorf("expected 3 brands, got %d", len(brands))
	}
	if len(svc.GetBrandsCalls()) != 1 {
		t.Errorf("expected 1 call to GetBrands, got %d", len(svc.GetBrandsCalls()))
	}
}

func TestProductHandler_GetBrands_Error(t *testing.T) {
	svc := &mocks.ProductServiceMock{
		GetBrandsFunc: func(_ context.Context) ([]string, error) {
			return nil, fmt.Errorf("db error")
		},
	}

	h := NewProductHandler(svc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/brands", nil)

	h.GetBrands(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}


func TestCategoryHandler_List_Success(t *testing.T) {
	svc := &mocks.CategoryServiceMock{
		GetAllFunc: func(_ context.Context) ([]domain.Category, error) {
			return []domain.Category{
				{Name: "Phones", Slug: "phones"},
				{Name: "Laptops", Slug: "laptops"},
			}, nil
		},
	}

	h := NewCategoryHandler(svc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/categories", nil)

	h.List(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var cats []domain.Category
	if err := json.Unmarshal(w.Body.Bytes(), &cats); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(cats) != 2 {
		t.Errorf("expected 2 categories, got %d", len(cats))
	}
	if len(svc.GetAllCalls()) != 1 {
		t.Errorf("expected 1 call to GetAll, got %d", len(svc.GetAllCalls()))
	}
}

func TestCategoryHandler_List_Error(t *testing.T) {
	svc := &mocks.CategoryServiceMock{
		GetAllFunc: func(_ context.Context) ([]domain.Category, error) {
			return nil, fmt.Errorf("db error")
		},
	}

	h := NewCategoryHandler(svc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/categories", nil)

	h.List(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestCategoryHandler_GetByID_Success(t *testing.T) {
	id := uuid.New()
	svc := &mocks.CategoryServiceMock{
		GetByIDFunc: func(_ context.Context, gotID uuid.UUID) (*domain.Category, error) {
			return &domain.Category{ID: gotID, Name: "Phones"}, nil
		},
	}

	h := NewCategoryHandler(svc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/categories/"+id.String(), nil)
	c.Params = gin.Params{{Key: "id", Value: id.String()}}

	h.GetByID(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if len(svc.GetByIDCalls()) != 1 {
		t.Errorf("expected 1 call to GetByID, got %d", len(svc.GetByIDCalls()))
	}
}

func TestCategoryHandler_GetByID_InvalidUUID(t *testing.T) {
	svc := &mocks.CategoryServiceMock{}
	h := NewCategoryHandler(svc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/categories/bad", nil)
	c.Params = gin.Params{{Key: "id", Value: "bad"}}

	h.GetByID(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCategoryHandler_GetByID_NotFound(t *testing.T) {
	svc := &mocks.CategoryServiceMock{
		GetByIDFunc: func(_ context.Context, _ uuid.UUID) (*domain.Category, error) {
			return nil, sql.ErrNoRows
		},
	}

	h := NewCategoryHandler(svc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	id := uuid.New()
	c.Request = httptest.NewRequest(http.MethodGet, "/categories/"+id.String(), nil)
	c.Params = gin.Params{{Key: "id", Value: id.String()}}

	h.GetByID(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}


func TestExchangeHandler_GetRate_Success(t *testing.T) {
	provider := &mocks.RateProviderMock{
		GetUSDTRateFunc: func(_ context.Context) (float64, error) {
			return 95.40, nil
		},
	}

	h := NewExchangeHandler(provider)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/exchange/rate", nil)

	h.GetRate(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp rateResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Rate != 95.40 {
		t.Errorf("expected rate 95.40, got %f", resp.Rate)
	}
	if len(provider.GetUSDTRateCalls()) != 1 {
		t.Errorf("expected 1 call to GetUSDTRate, got %d", len(provider.GetUSDTRateCalls()))
	}
}

func TestExchangeHandler_GetRate_Error(t *testing.T) {
	provider := &mocks.RateProviderMock{
		GetUSDTRateFunc: func(_ context.Context) (float64, error) {
			return 0, fmt.Errorf("grinex unavailable")
		},
	}

	h := NewExchangeHandler(provider)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/exchange/rate", nil)

	h.GetRate(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}


func TestErrorResponse_Format(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	errorResponse(c, http.StatusBadRequest, "test error")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}

	var resp ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Error != "test error" {
		t.Errorf("expected 'test error', got %q", resp.Error)
	}
}

func TestProductHandler_List_DefaultPagination(t *testing.T) {
	var capturedFilter domain.ProductFilter
	svc := &mocks.ProductServiceMock{
		GetByFilterFunc: func(_ context.Context, f domain.ProductFilter) (*domain.ProductList, error) {
			capturedFilter = f
			return &domain.ProductList{Products: []domain.Product{}, Total: 0, Page: 1, PageSize: 24}, nil
		},
	}

	h := NewProductHandler(svc)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/products", nil)

	h.List(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if capturedFilter.Limit != 24 {
		t.Errorf("expected default limit 24, got %d", capturedFilter.Limit)
	}
	if capturedFilter.Offset != 0 {
		t.Errorf("expected offset 0, got %d", capturedFilter.Offset)
	}
}
