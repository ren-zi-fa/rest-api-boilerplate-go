APP_PORT ?= 8080
DOCKER_COMPOSE_FILE = docker-compose.prod.yml
GO_BUILD_OUTPUT = bin/main

.PHONY: build docker up wait migrate run-prod clean


ENV:
	@if [ -f .env ]; then \
		export $$(cat .env | xargs); \
	fi

build:
	@echo ">> Building Go application..."
	@go build -o $(GO_BUILD_OUTPUT) ./cmd


docker:
	@echo ">> Building Docker containers..."
	@docker compose -f $(DOCKER_COMPOSE_FILE) build


up:
	@echo ">> Starting Docker containers..."
	@docker compose -f $(DOCKER_COMPOSE_FILE) up -d


wait:
	@echo ">> Waiting for MySQL to be ready..."
	@until docker exec my-mysql-db mysqladmin ping -h"localhost" --silent; do \
		echo "Waiting for database..."; \
		sleep 2; \
	done


migrate:
	@echo ">> Running database migrations..."
	@go run cmd/migrate/main.go up


run-prod: ENV build docker up wait migrate
	@echo ">> Application started!"
	@echo "Visit: http://localhost:$(APP_PORT)/api/posts"

clean:
	@echo ">> Cleaning build and docker containers..."
	@rm -f $(GO_BUILD_OUTPUT)
	@docker compose -f $(DOCKER_COMPOSE_FILE) down
