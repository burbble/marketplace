package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID            uuid.UUID `db:"id" json:"id"`
	ExternalID    string    `db:"external_id" json:"external_id"`
	SKU           string    `db:"sku" json:"sku"`
	Name          string    `db:"name" json:"name"`
	OriginalPrice int       `db:"original_price" json:"original_price"`
	Price         int       `db:"price" json:"price"`
	ImageURL      string    `db:"image_url" json:"image_url"`
	ProductURL    string    `db:"product_url" json:"product_url"`
	Brand         string    `db:"brand" json:"brand"`
	CategoryID    uuid.UUID `db:"category_id" json:"category_id"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

type ProductFilter struct {
	CategoryID *uuid.UUID `form:"category_id"`
	Brand      *string    `form:"brand"`
	MinPrice   *int       `form:"min_price"`
	MaxPrice   *int       `form:"max_price"`
	Search     *string    `form:"search"`
	Page       int        `form:"page"`
	PageSize   int        `form:"page_size"`
}

type ProductList struct {
	Products []Product `json:"products"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
}
