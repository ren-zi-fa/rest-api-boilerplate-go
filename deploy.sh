#!/bin/bash
set -e  

echo ">> Loading environment variables..."
if [ -f .env.prod ]; then
  source .env.prod
fi

echo $DB_PORT
echo ">> Building Go application..."
go build -o bin/main ./cmd

echo ">> Building Docker containers..."
docker compose -p prod --env-file .env.prod -f docker-compose.prod.yml build

echo ">> Starting Docker containers..."
docker compose -p prod --env-file .env.prod -f docker-compose.prod.yml up -d

echo ">> Waiting for MySQL to be ready..."
until docker exec db-prod mysqladmin ping -h"localhost" --silent; do
    echo "Waiting for database..."
    sleep 2
done

echo ">> Running database migrations..."
docker exec -it go_app ./migrate up "$@"

echo ">> Application started!"
echo "Visit: http://localhost:${APP_PORT:-8081}/api/v1/posts"
