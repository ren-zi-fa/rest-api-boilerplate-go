#!/usr/bin/env bash

# Load environment variables
if [ -f .env.prod ]; then
  source .env.prod
fi

set -e 

CMD=$1
shift   

case "$CMD" in
  up)
    echo ">> Setting up database if not exists... $DB_NAME"
    mysql -u "$DB_USER" -p"$DB_PASSWORD" -h "$DB_HOST" -P "$DB_PORT" \
      -e "CREATE SCHEMA IF NOT EXISTS $DB_NAME;"

    echo ">> Migrating up (inside container go_app)..."
   docker exec -it go_app ./migrate up "$@"
    ;;
  down)
    echo ">> Dropping database schema... $DB_NAME"
    mysql -u "$DB_USER" -p"$DB_PASSWORD" -h "$DB_HOST" -P "$DB_PORT" \
      -e "DROP SCHEMA IF EXISTS $DB_NAME;"
    ;;
  *)
    echo "Usage: $0 {up|down}"
    exit 1
    ;;
esac
