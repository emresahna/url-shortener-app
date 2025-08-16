# Version
VERSION ?= 0.1.0

# Build directories
BIN_DIR = bin
DOCS_DIR = internal/docs

# Build targets
.PHONY: all build build-api build-scheduler build-worker test run-api run-scheduler run-worker clean lint swagger help

all: clean swagger build test

# Build commands
build: build-api build-scheduler build-worker

build-api:
	go build -o $(BIN_DIR)/api -ldflags "-X main.Version=$(VERSION)" ./cmd/api

build-scheduler:
	go build -o $(BIN_DIR)/scheduler -ldflags "-X main.Version=$(VERSION)" ./cmd/scheduler

build-worker:
	go build -o $(BIN_DIR)/worker -ldflags "-X main.Version=$(VERSION)" ./cmd/worker

# Test commands
test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Run commands
run-api:
	$(BIN_DIR)/api

run-scheduler:
	$(BIN_DIR)/scheduler

run-worker:
	$(BIN_DIR)/worker

# Development commands
dev-api:
	go run ./cmd/api

dev-scheduler:
	go run ./cmd/scheduler

dev-worker:
	go run ./cmd/worker

# Clean commands
clean:
	rm -rf $(BIN_DIR)/ coverage.out

# Lint commands
lint:
	golangci-lint run

# Swagger commands
swagger-install:
	go install github.com/swaggo/swag/cmd/swag@latest

swagger: swagger-install
	swag init -g cmd/api/main.go -o $(DOCS_DIR)

# Docker commands
docker-build:
	docker build -t url-shortener:$(VERSION) .

docker-run:
	docker run -p 8080:8080 url-shortener:$(VERSION)

# Help command
help:
	@echo "Available targets:"
	@echo "  all            : Clean, generate swagger docs, build, and test"
	@echo "  build          : Build all binaries"
	@echo "  build-api      : Build API server"
	@echo "  build-scheduler: Build scheduler"
	@echo "  build-worker   : Build worker"
	@echo "  test           : Run tests"
	@echo "  test-coverage  : Run tests with coverage report"
	@echo "  run-api        : Run API server from binary"
	@echo "  run-scheduler  : Run scheduler from binary"
	@echo "  run-worker     : Run worker from binary"
	@echo "  dev-api        : Run API server directly with go run"
	@echo "  dev-scheduler  : Run scheduler directly with go run"
	@echo "  dev-worker     : Run worker directly with go run"
	@echo "  clean          : Remove build artifacts"
	@echo "  lint           : Run linter"
	@echo "  swagger        : Generate Swagger documentation"
	@echo "  docker-build   : Build Docker image"
	@echo "  docker-run     : Run Docker container"
	@echo "  help           : Show this help message"