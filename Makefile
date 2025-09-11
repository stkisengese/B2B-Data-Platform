.PHONY: build test lint run up down

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

build:
	$(GO_BUILD) -o bin/server $(SERVER_MAIN)
	$(GO_BUILD) -o bin/collector $(COLLECTOR_MAIN)
	$(GO_BUILD) -o bin/migrator $(MIGRATOR_MAIN)

test:
	$(GO_TEST) ./...

lint:
	$(GO_LINT) ./...

run:
	$(GO_RUN) $(SERVER_MAIN)

up:
	$(DOCKER_COMPOSE) up -d

down:
	$(DOCKER_COMPOSE) down