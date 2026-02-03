-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS categories (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT         NOT NULL,
    slug       TEXT         NOT NULL UNIQUE,
    url        TEXT         NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_categories_slug ON categories (slug);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS categories;
