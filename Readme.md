# Outfit Server

## Setup
1. Install dependencies: `go mod tidy`
2. Configure `config.yaml` or `.env`.
3. Run SQL script: `psql -U app_user -d outfit -f schema.sql`
4. Generate Swagger: `make swag`
5. Run: `make run`
6. Swagger: http://localhost:8080/swagger/index.html

## Стек
- Go 1.21, Gin, GORM, PostgreSQL 18, JWT, Viper, Zap, Swagger.

## БД
- Схемы: schema.sql.

## Запуск
make run