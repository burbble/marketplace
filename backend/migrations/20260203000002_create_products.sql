-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products (
    id             UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    external_id    TEXT         NOT NULL UNIQUE,
    sku            TEXT         NOT NULL DEFAULT '',
    name           TEXT         NOT NULL,
    original_price INTEGER      NOT NULL DEFAULT 0,
    price          INTEGER      NOT NULL DEFAULT 0,
    image_url      TEXT         NOT NULL DEFAULT '',
    product_url    TEXT         NOT NULL DEFAULT '',
    brand          TEXT         NOT NULL DEFAULT '',
    category_id    UUID         NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_products_category_id ON products (category_id);
CREATE INDEX idx_products_brand ON products (brand);
CREATE INDEX idx_products_price ON products (price);
CREATE INDEX idx_products_external_id ON products (external_id);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS products;
