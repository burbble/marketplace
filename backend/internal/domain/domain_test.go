package domain

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestProduct_JSONSerialization(t *testing.T) {
	id := uuid.New()
	catID := uuid.New()
	now := time.Now().Truncate(time.Second)

	p := Product{
		ID:            id,
		ExternalID:    "ext-123",
		SKU:           "SKU001",
		Name:          "Test Product",
		OriginalPrice: 50000,
		Price:         49000,
		ImageURL:      "/img/test.jpg",
		ProductURL:    "/products/test/",
		Brand:         "TestBrand",
		Description:   "Product description",
		CategoryID:    catID,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded Product
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if decoded.ID != id {
		t.Errorf("ID: expected %s, got %s", id, decoded.ID)
	}
	if decoded.ExternalID != "ext-123" {
		t.Errorf("ExternalID: expected 'ext-123', got %q", decoded.ExternalID)
	}
	if decoded.SKU != "SKU001" {
		t.Errorf("SKU: expected 'SKU001', got %q", decoded.SKU)
	}
	if decoded.Name != "Test Product" {
		t.Errorf("Name: expected 'Test Product', got %q", decoded.Name)
	}
	if decoded.OriginalPrice != 50000 {
		t.Errorf("OriginalPrice: expected 50000, got %d", decoded.OriginalPrice)
	}
	if decoded.Price != 49000 {
		t.Errorf("Price: expected 49000, got %d", decoded.Price)
	}
	if decoded.Description != "Product description" {
		t.Errorf("Description: expected 'Product description', got %q", decoded.Description)
	}
	if decoded.Brand != "TestBrand" {
		t.Errorf("Brand: expected 'TestBrand', got %q", decoded.Brand)
	}
	if decoded.CategoryID != catID {
		t.Errorf("CategoryID: expected %s, got %s", catID, decoded.CategoryID)
	}
}

func TestProduct_JSONFieldNames(t *testing.T) {
	p := Product{
		ID:          uuid.New(),
		ExternalID:  "ext",
		Name:        "Test",
		Description: "desc",
		CategoryID:  uuid.New(),
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	expectedFields := []string{
		"id", "external_id", "sku", "name", "original_price", "price",
		"image_url", "product_url", "brand", "description", "category_id",
		"created_at", "updated_at",
	}

	for _, field := range expectedFields {
		if _, ok := raw[field]; !ok {
			t.Errorf("missing JSON field %q", field)
		}
	}
}

func TestCategory_JSONSerialization(t *testing.T) {
	id := uuid.New()
	now := time.Now().Truncate(time.Second)

	c := Category{
		ID:        id,
		Name:      "Phones",
		Slug:      "phones",
		URL:       "/phones/",
		CreatedAt: now,
		UpdatedAt: now,
	}

	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded Category
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if decoded.ID != id {
		t.Errorf("ID: expected %s, got %s", id, decoded.ID)
	}
	if decoded.Name != "Phones" {
		t.Errorf("Name: expected 'Phones', got %q", decoded.Name)
	}
	if decoded.Slug != "phones" {
		t.Errorf("Slug: expected 'phones', got %q", decoded.Slug)
	}
	if decoded.URL != "/phones/" {
		t.Errorf("URL: expected '/phones/', got %q", decoded.URL)
	}
}

func TestProductList_JSONSerialization(t *testing.T) {
	list := ProductList{
		Products: []Product{
			{Name: "A", Price: 100},
			{Name: "B", Price: 200},
		},
		Total:    2,
		Page:     1,
		PageSize: 24,
	}

	data, err := json.Marshal(list)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded ProductList
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if decoded.Total != 2 {
		t.Errorf("Total: expected 2, got %d", decoded.Total)
	}
	if decoded.Page != 1 {
		t.Errorf("Page: expected 1, got %d", decoded.Page)
	}
	if decoded.PageSize != 24 {
		t.Errorf("PageSize: expected 24, got %d", decoded.PageSize)
	}
	if len(decoded.Products) != 2 {
		t.Errorf("Products length: expected 2, got %d", len(decoded.Products))
	}
}

func TestProductList_EmptyProducts(t *testing.T) {
	list := ProductList{
		Products: []Product{},
		Total:    0,
		Page:     1,
		PageSize: 24,
	}

	data, err := json.Marshal(list)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded ProductList
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if decoded.Products == nil {
		t.Error("expected non-nil empty slice, got nil")
	}
	if len(decoded.Products) != 0 {
		t.Errorf("expected 0 products, got %d", len(decoded.Products))
	}
}

func TestProductFilter_ZeroValue(t *testing.T) {
	f := ProductFilter{}

	if f.CategoryID != nil {
		t.Error("expected nil CategoryID")
	}
	if f.Brand != nil {
		t.Error("expected nil Brand")
	}
	if f.MinPrice != nil {
		t.Error("expected nil MinPrice")
	}
	if f.MaxPrice != nil {
		t.Error("expected nil MaxPrice")
	}
	if f.Search != nil {
		t.Error("expected nil Search")
	}
	if f.Limit != 0 {
		t.Errorf("expected Limit 0, got %d", f.Limit)
	}
	if f.Offset != 0 {
		t.Errorf("expected Offset 0, got %d", f.Offset)
	}
}
