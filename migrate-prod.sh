#!/usr/bin/env bash

# Load environment variables
if [ -f .env.prod ]; then
  source .env.prod
fi

echo ">> Setting up database if not exists... $DB_NAME"

echo ">> Migrating up..."
mysql -u "$DB_USER" -p"$DB_PASSWORD" -h "$DB_HOST" -P "$DB_PORT" \
  -e "CREATE SCHEMA IF NOT EXISTS $DB_NAME;"

go run cmd/migrate/main.go up
