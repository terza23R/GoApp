APP_NAME := goapp
CMD_PATH := ./cmd/main.go

.PHONY: help config docker-up docker-down docker-logs run build test fmt tidy

help:
	@echo "Available targets:"
	@echo "  make config       - create config.yaml from config.yaml.example (if missing)"
	@echo "  make docker-up    - start MySQL via docker compose"
	@echo "  make docker-down  - stop containers"
	@echo "  make docker-logs  - tail docker compose logs"
	@echo "  make run          - run the app"
	@echo "  make build        - build binary into ./bin/$(APP_NAME)"
	@echo "  make test         - run tests"
	@echo "  make fmt          - format code"
	@echo "  make tidy         - go mod tidy"

config:
	@test -f config.yaml || cp config.yaml.example config.yaml
	@echo "config.yaml is ready"

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f --tail=100

run:
	-go run $(CMD_PATH)

build:
	mkdir -p bin
	go build -o bin/$(APP_NAME) $(CMD_PATH)

test:
	go test ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy
