package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/burbble/marketplace/internal/domain"
	"github.com/burbble/marketplace/pkg/db"
)

type CategoryRepository interface {
	Upsert(ctx context.Context, categories []domain.Category) error
	GetAll(ctx context.Context) ([]domain.Category, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Category, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error)
}

type categoryRepo struct {
	conn *db.Connection
}

func NewCategoryRepo(conn *db.Connection) CategoryRepository {
	return &categoryRepo{conn: conn}
}

func (r *categoryRepo) Upsert(ctx context.Context, categories []domain.Category) error {
	if len(categories) == 0 {
		return nil
	}

	q := r.conn.Builder.
		Insert("categories").
		Columns("name", "slug", "url", "updated_at")

	now := time.Now()
	for _, c := range categories {
		q = q.Values(c.Name, c.Slug, c.URL, now)
	}

	q = q.Suffix("ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name, url = EXCLUDED.url, updated_at = EXCLUDED.updated_at")

	query, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("build upsert categories: %w", err)
	}

	_, err = r.conn.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec upsert categories: %w", err)
	}

	return nil
}

func (r *categoryRepo) GetAll(ctx context.Context) ([]domain.Category, error) {
	query, args, err := r.conn.Builder.
		Select("id", "name", "slug", "url", "created_at", "updated_at").
		From("categories").
		OrderBy("name ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select categories: %w", err)
	}

	var categories []domain.Category
	if err := r.conn.DB.SelectContext(ctx, &categories, query, args...); err != nil {
		return nil, fmt.Errorf("select categories: %w", err)
	}

	return categories, nil
}

func (r *categoryRepo) GetBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	query, args, err := r.conn.Builder.
		Select("id", "name", "slug", "url", "created_at", "updated_at").
		From("categories").
		Where("slug = ?", slug).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select category by slug: %w", err)
	}

	var cat domain.Category
	if err := r.conn.DB.GetContext(ctx, &cat, query, args...); err != nil {
		return nil, fmt.Errorf("get category by slug: %w", err)
	}

	return &cat, nil
}

func (r *categoryRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	query, args, err := r.conn.Builder.
		Select("id", "name", "slug", "url", "created_at", "updated_at").
		From("categories").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select category by id: %w", err)
	}

	var cat domain.Category
	if err := r.conn.DB.GetContext(ctx, &cat, query, args...); err != nil {
		return nil, fmt.Errorf("get category by id: %w", err)
	}

	return &cat, nil
}
