# Используем официальный образ Go
FROM golang:1.24.4 AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Собираем приложение
RUN go build -o backend ./cmd/main.go

# Минимальный runtime образ
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/backend .

# Копируем конфиг
COPY .env .

EXPOSE 8080

CMD ["./backend"]
