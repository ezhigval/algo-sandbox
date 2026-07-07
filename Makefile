.PHONY: run test bench docker-up docker-down build

run:
	REDIS_ADDR=localhost:6382 go run ./cmd/server

test:
	go test ./... -race -count=1

bench:
	go test ./internal/algo/... -bench=. -benchmem -count=1

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

build:
	go build -o bin/server ./cmd/server
	go build -o bin/algo ./cmd/algo
