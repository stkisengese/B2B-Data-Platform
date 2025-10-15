.PHONY: build test test-unit test-integration test-coverage test-bench lint run run-collector migrate up down clean

# Go parameters
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_TEST=$(GO_CMD) test
GO_LINT=golangci-lint run
GO_RUN=$(GO_CMD) run

# Docker parameters
DOCKER_COMPOSE=docker-compose

# Service parameters
SERVER_MAIN=cmd/server/main.go
COLLECTOR_MAIN=cmd/collector/main.go
MIGRATOR_MAIN=cmd/migrator/main.go

# Build all binaries
build:
	$(GO_BUILD) -o bin/server $(SERVER_MAIN)
	$(GO_BUILD) -o bin/collector $(COLLECTOR_MAIN)
	$(GO_BUILD) -o bin/migrator $(MIGRATOR_MAIN)

# Run all tests
test: test-unit

# Run unit tests
test-unit:
	$(GO_TEST) -v ./internal/...

# Run integration tests
test-integration:
	$(GO_TEST) -v -tags=integration ./...

# Generate test coverage report
test-coverage:
	$(GO_TEST) -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmark tests
test-bench:
	$(GO_TEST) -bench=. -benchmem ./...

# Run linter
lint:
	$(GO_LINT) ./...

# Run API server
run:
	$(GO_RUN) $(SERVER_MAIN)

# Run data collector service
run-collector:
	$(GO_RUN) $(COLLECTOR_MAIN)

# Run database migrations
migrate:
	$(GO_RUN) $(MIGRATOR_MAIN)

# Start all services with Docker Compose
up:
	$(DOCKER_COMPOSE) up -d

# Stop all Docker services
down:
	$(DOCKER_COMPOSE) down

# Clean build artifacts and test files
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f *.db

# Install dependencies
deps:
	$(GO_CMD) mod download
	$(GO_CMD) mod tidy

# Install development tools
dev-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Quick development setup
dev-setup: deps dev-tools migrate
	@echo "Development environment ready!"

# Verify everything is working
verify: lint test build
	@echo "All checks passed!"