package postgres

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/burbble/marketplace/internal/domain"
	"github.com/burbble/marketplace/pkg/db"
)

type ProductRepository interface {
	Upsert(ctx context.Context, products []domain.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	GetByFilter(ctx context.Context, filter domain.ProductFilter) (*domain.ProductList, error)
	GetBrands(ctx context.Context) ([]string, error)
}

type productRepo struct {
	conn *db.Connection
}

func NewProductRepo(conn *db.Connection) ProductRepository {
	return &productRepo{conn: conn}
}

func (r *productRepo) Upsert(ctx context.Context, products []domain.Product) error {
	if len(products) == 0 {
		return nil
	}

	q := r.conn.Builder.
		Insert("products").
		Columns(
			"external_id", "sku", "name", "original_price", "price",
			"image_url", "product_url", "brand", "category_id", "updated_at",
		)

	now := time.Now()
	for _, p := range products {
		q = q.Values(
			p.ExternalID, p.SKU, p.Name, p.OriginalPrice, p.Price,
			p.ImageURL, p.ProductURL, p.Brand, p.CategoryID, now,
		)
	}

	q = q.Suffix(`ON CONFLICT (external_id) DO UPDATE SET
		sku = EXCLUDED.sku,
		name = EXCLUDED.name,
		original_price = EXCLUDED.original_price,
		price = EXCLUDED.price,
		image_url = EXCLUDED.image_url,
		product_url = EXCLUDED.product_url,
		brand = EXCLUDED.brand,
		category_id = EXCLUDED.category_id,
		updated_at = EXCLUDED.updated_at`)

	query, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("build upsert products: %w", err)
	}

	_, err = r.conn.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec upsert products: %w", err)
	}

	return nil
}

func (r *productRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	query, args, err := r.conn.Builder.
		Select(
			"id", "external_id", "sku", "name", "original_price", "price",
			"image_url", "product_url", "brand", "category_id", "created_at", "updated_at",
		).
		From("products").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select product by id: %w", err)
	}

	var p domain.Product
	if err := r.conn.DB.GetContext(ctx, &p, query, args...); err != nil {
		return nil, fmt.Errorf("get product by id: %w", err)
	}

	return &p, nil
}

func (r *productRepo) GetByFilter(ctx context.Context, filter domain.ProductFilter) (*domain.ProductList, error) {
	where := buildProductWhere(filter)

	countQ, countArgs, err := r.conn.Builder.
		Select("COUNT(*)").
		From("products").
		Where(where).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build count products: %w", err)
	}

	var total int
	if err := r.conn.DB.GetContext(ctx, &total, countQ, countArgs...); err != nil {
		return nil, fmt.Errorf("count products: %w", err)
	}

	q := r.conn.Builder.
		Select(
			"id", "external_id", "sku", "name", "original_price", "price",
			"image_url", "product_url", "brand", "category_id", "created_at", "updated_at",
		).
		From("products").
		Where(where).
		Limit(filter.Limit).
		Offset(filter.Offset)

	if len(filter.SortBy) > 0 {
		q = q.OrderBy(filter.SortBy...)
	} else {
		q = q.OrderBy("created_at DESC")
	}

	dataQ, dataArgs, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select products: %w", err)
	}

	var products []domain.Product
	if err := r.conn.DB.SelectContext(ctx, &products, dataQ, dataArgs...); err != nil {
		return nil, fmt.Errorf("select products: %w", err)
	}

	page := uint64(1)
	if filter.Limit > 0 {
		page = filter.Offset/filter.Limit + 1
	}

	return &domain.ProductList{
		Products: products,
		Total:    total,
		Page:     page,
		PageSize: filter.Limit,
	}, nil
}

func (r *productRepo) GetBrands(ctx context.Context) ([]string, error) {
	query, args, err := r.conn.Builder.
		Select("DISTINCT brand").
		From("products").
		Where("brand != ''").
		OrderBy("brand ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select brands: %w", err)
	}

	var brands []string
	if err := r.conn.DB.SelectContext(ctx, &brands, query, args...); err != nil {
		return nil, fmt.Errorf("select brands: %w", err)
	}

	return brands, nil
}

func buildProductWhere(f domain.ProductFilter) sq.And {
	var conds sq.And

	if f.CategoryID != nil {
		conds = append(conds, sq.Eq{"category_id": *f.CategoryID})
	}
	if f.Brand != nil && *f.Brand != "" {
		conds = append(conds, sq.Eq{"brand": *f.Brand})
	}
	if f.MinPrice != nil {
		conds = append(conds, sq.GtOrEq{"price": *f.MinPrice})
	}
	if f.MaxPrice != nil {
		conds = append(conds, sq.LtOrEq{"price": *f.MaxPrice})
	}
	if f.Search != nil && *f.Search != "" {
		conds = append(conds, sq.ILike{"name": "%" + *f.Search + "%"})
	}

	return conds
}
