PG_DSN ?= "host=localhost port=54320 user=postgres password=postgres dbname=store_scraper sslmode=disable"
MIGRATIONS_DIR = ./backend/migrations
MOQ = $(shell which moq)

.PHONY: init up down rebuild rebuild-api rebuild-parser rebuild-frontend logs logs-api logs-parser logs-frontend migrate migrate-down migrate-status swagger test lint clean generate-mocks test-frontend lint-frontend format-frontend

init:
	@test -f .env || (cp .env.example .env && echo "Created .env from .env.example")
	docker compose up --build -d postgres redis
	@echo "Waiting for postgres to be ready..."
	@until docker compose exec -T postgres pg_isready -U postgres > /dev/null 2>&1; do sleep 1; done
	goose -dir=$(MIGRATIONS_DIR) postgres $(PG_DSN) up -v
	docker compose up --build -d
	@echo ""
	@echo "Ready!"
	@echo "  Frontend: http://localhost:3000"
	@echo "  API:      http://localhost:38080"
	@echo "  Swagger:  http://localhost:38080/swagger/index.html"

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

test-frontend:
	cd frontend && npx vitest run

lint-frontend:
	cd frontend && npx eslint src/

format-frontend:
	cd frontend && npx prettier --write src/
