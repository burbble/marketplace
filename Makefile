PG_DSN ?= "host=localhost port=54320 user=postgres password=postgres dbname=store_scraper sslmode=disable"
MIGRATIONS_DIR = ./backend/migrations

.PHONY: up down rebuild rebuild-api rebuild-parser rebuild-frontend logs logs-api logs-parser logs-frontend migrate migrate-down migrate-status swagger test clean

up:
	docker compose up --build -d

down:
	docker compose down

rebuild:
	docker compose down
	docker compose up --build -d

rebuild-api:
	docker compose up --build -d api

rebuild-parser:
	docker compose up --build -d parser

rebuild-frontend:
	docker compose up --build -d frontend

clean:
	docker compose down -v --rmi all

logs:
	docker compose logs -f

logs-api:
	docker compose logs -f api

logs-parser:
	docker compose logs -f parser

logs-frontend:
	docker compose logs -f frontend

migrate:
	goose -dir=$(MIGRATIONS_DIR) postgres $(PG_DSN) up -v

migrate-down:
	goose -dir=$(MIGRATIONS_DIR) postgres $(PG_DSN) down -v

migrate-status:
	goose -dir=$(MIGRATIONS_DIR) postgres $(PG_DSN) status

swagger:
	cd backend && swag init -g cmd/api/main.go -o docs

test:
	cd backend && go test ./... -v -count=1
