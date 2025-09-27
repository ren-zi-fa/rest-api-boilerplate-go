#!/bin/bash
set -e  

echo ">> Loading environment variables..."
if [ -f .env ]; then
  source .env
fi

echo ">> Building Go application..."
go build -o bin/main ./cmd

echo ">> Building Docker containers..."
docker compose -f docker-compose.prod.yml build

echo ">> Starting Docker containers..."
docker compose -f docker-compose.prod.yml up -d

echo ">> Waiting for MySQL to be ready..."
until docker exec db-prod mysqladmin ping -h"localhost" --silent; do
    echo "Waiting for database..."
    sleep 2
done

echo ">> Running database migrations..."
go run cmd/migrate/main.go up

echo ">> Application started!"
echo "Visit: http://localhost:${APP_PORT:-8081}/api/v1/posts"
