PG_DSN ?= "host=localhost port=54320 user=postgres password=postgres dbname=store_scraper sslmode=disable"
MIGRATIONS_DIR = ./backend/migrations
MOQ = $(shell which moq)

.PHONY: up down rebuild rebuild-api rebuild-parser rebuild-frontend logs logs-api logs-parser logs-frontend migrate migrate-down migrate-status swagger test lint clean generate-mocks

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
	go test -C backend ./... -v -count=1

lint:
	cd backend && golangci-lint run ./...

generate-mocks:
	cd backend && $(MOQ) -pkg mocks -out internal/mocks/service_mock.go internal/service ProductService CategoryService
	cd backend && $(MOQ) -pkg mocks -out internal/mocks/repository_mock.go internal/repository/postgres ProductRepository CategoryRepository
	cd backend && $(MOQ) -pkg mocks -out internal/mocks/exchange_mock.go internal/exchange RateProvider
