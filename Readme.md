# Outfit Server

## Setup
1. Install dependencies: `go mod tidy`
2. Configure `config.yaml` or `.env`.
3. Run SQL script: `psql -U app_user -d outfit -f schema.sql`
4. Generate Swagger: `make swag`
5. Run: `make run`
6. Swagger: http://localhost:8080/swagger/index.html

## Стек
- Go, Gin, GORM, PostgreSQL, JWT, Viper, Zap, Swagger.

## БД
- Схема: `schema.sql`.

## Запуск
```bash
make run
```

## Health
- `GET /healthz` (liveness)
- `GET /readyz` (readiness + DB check)

## Публичный API-контракт (`/api`)

### Auth
- `POST /api/auth/register`
- `POST /api/auth/login`
- `POST /api/auth/refresh`
- `POST /api/auth/logout` *(JWT)*
- `GET /api/auth/profile` *(JWT)*

### Catalog
- `GET /api/products`
- `GET /api/products/:id`
- `GET /api/categories`
- `GET /api/categories/tree`

### Cart *(JWT)*
- `GET /api/cart`
- `POST /api/cart/items`
- `DELETE /api/cart/items/:id`

### Orders *(JWT)*
- `POST /api/orders`
- `GET /api/orders`
- `GET /api/orders/my` *(compat route)*
- `GET /api/orders/:id`

## Стандартизированные response envelopes
- products list: `{ "products": [...] }`
- categories list/tree: `{ "categories": [...] }`
- order by id: `{ "order": {...} }`
- orders list: `{ "orders": [...] }`

## Server tuning (optional config)
`server` section supports runtime tuning:
- `body_limit_bytes`
- `read_timeout_sec`
- `read_header_timeout_sec`
- `write_timeout_sec`
- `idle_timeout_sec`
- `shutdown_timeout_sec`
- `rate_limit_rps`
- `rate_limit_burst`

All values have safe defaults in code if omitted.
