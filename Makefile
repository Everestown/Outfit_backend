.PHONY: build run test swag clean

build:
	go build -o bin/server ./cmd/server

run:
	go run cmd/server/main.go

test:
	go test ./...

swag:
	swag init -g cmd/server/main.go -o internal/docs

clean:
	rm -rf bin/