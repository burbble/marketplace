# Marketplace

Парсер товаров с store77.net с веб-интерфейсом. Go (API + парсер), Next.js (фронтенд), PostgreSQL, Redis.

## Требования

- Docker, Docker Compose
- Go 1.23+ (для локальной разработки)
- [goose](https://github.com/pressly/goose) (для миграций)

## Быстрый старт

1. Скопировать `.env.example` → `.env`:

```bash
cp .env.example .env
```

2. Запустить все сервисы:

```bash
make up
```

3. Применить миграции:

```bash
make migrate
```

Фронтенд: http://localhost:3000
API: http://localhost:38080

## Переменные окружения

Настраиваются в `.env` (корень проекта). Основные:

| Переменная | По умолчанию | Описание |
|---|---|---|
| `PG_HOST` | postgres | Хост PostgreSQL |
| `PG_PORT` | 5432 | Порт PostgreSQL (внутри Docker) |
| `PG_USER` | postgres | Пользователь БД |
| `PG_PASSWORD` | postgres | Пароль БД |
| `PG_DB_NAME` | store_scraper | Имя БД |
| `REDIS_HOST` | redis | Хост Redis |
| `REDIS_PORT` | 6379 | Порт Redis (внутри Docker) |
| `HTTP_PORT` | 8080 | Порт API |
| `GIN_MODE` | debug | Режим Gin (debug/release) |
| `RATE_LIMIT_RPS` | 100 | Лимит запросов в секунду |
| `SCRAPE_INTERVAL` | 10m | Интервал между циклами парсинга |
| `SCRAPE_WORKERS` | 5 | Количество параллельных воркеров парсера |
| `BACKEND_URL` | http://api:8080 | URL бэкенда для фронтенда |

## Makefile команды

```
make up                 — запустить всё
make down               — остановить
make rebuild            — пересобрать всё
make rebuild-api        — пересобрать только API
make rebuild-parser     — пересобрать только парсер
make rebuild-frontend   — пересобрать только фронтенд
make clean              — остановить и удалить volumes/images
make logs               — логи всех сервисов
make logs-api           — логи API
make logs-parser        — логи парсера
make logs-frontend      — логи фронтенда
make migrate            — применить миграции
make migrate-down       — откатить миграцию
make migrate-status     — статус миграций
make swagger            — сгенерировать Swagger
make test               — запустить тесты
```

## Структура проекта

```
backend/
├── cmd/api/          — точка входа API сервера
├── cmd/parser/       — точка входа парсера
├── internal/
│   ├── config/       — конфигурация (viper)
│   ├── domain/       — доменные модели
│   ├── handler/      — HTTP хэндлеры (Gin)
│   ├── repository/   — работа с БД (sqlx + squirrel)
│   ├── service/      — бизнес-логика
│   └── scraper/      — парсинг store77.net (rod)
├── pkg/
│   ├── postgres/     — подключение к PostgreSQL
│   ├── ratelimit/    — rate limiter (Redis)
│   ├── pagination/   — пагинация и сортировка
│   └── zapx/         — настройка логгера
├── migrations/       — SQL миграции (goose)
└── docs/             — Swagger

frontend/             — Next.js 15, FSD архитектура
├── src/
│   ├── app/          — страницы (каталог, детальная)
│   ├── entities/     — product, category, exchange
│   ├── features/     — фильтры, сортировка
│   ├── widgets/      — header, product-grid, pagination
│   └── shared/       — API клиент, утилиты, UI
```

## API

Swagger UI доступен по адресу http://localhost:38080/swagger/index.html (при `GIN_MODE=debug`).

Основные эндпоинты:

```
GET  /api/v1/products          — список товаров (фильтры, пагинация, сортировка)
GET  /api/v1/products/:id      — товар по ID
GET  /api/v1/brands            — список брендов
GET  /api/v1/categories        — список категорий
GET  /api/v1/categories/:id    — категория по ID
GET  /api/v1/exchange/rate     — курс USDT/RUB
GET  /health                   — healthcheck
```

## Локальная разработка (без Docker)

Для бэкенда скопировать `backend/.env.example` → `backend/.env` и запустить:

```bash
cd backend && go run ./cmd/api
```

Для фронтенда скопировать `frontend/.env.example` → `frontend/.env.local` и запустить:

```bash
cd frontend && npm install && npm run dev
```

Для локальной разработки PostgreSQL и Redis должны быть запущены (порты 54320 и 63790).
