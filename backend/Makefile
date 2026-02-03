PG_DSN ?= "host=localhost port=54320 user=postgres password=postgres dbname=store_scraper sslmode=disable"
MIGRATIONS_DIR = ./migrations

.PHONY: migrate migrate-down migrate-status new-migration build-api build-parser swagger

migrate:
	goose -dir=$(MIGRATIONS_DIR) postgres $(PG_DSN) up -v

migrate-down:
	goose -dir=$(MIGRATIONS_DIR) postgres $(PG_DSN) down -v

migrate-status:
	goose -dir=$(MIGRATIONS_DIR) postgres $(PG_DSN) status

new-migration:
	goose -dir=$(MIGRATIONS_DIR) create $(name) sql

swagger:
	# go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g cmd/api/main.go -o docs

build-api:
	go build -o bin/api ./cmd/api

build-parser:
	go build -o bin/parser ./cmd/parser
