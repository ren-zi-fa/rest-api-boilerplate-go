#!/usr/bin/env bash

# Load environment variables
if [ -f .env ]; then
  source .env
fi

set -e 

CMD=$1
shift   

case "$CMD" in
  build)
    echo ">> Building..."
    go build -o bin/app cmd/main.go
      echo ">> Done...."
    ;;
  
  test)
    echo ">> Running tests..."
    go test -v ./...
    ;;
  
  run)
    echo ">> Building..."
    go build -o bin/app cmd/main.go
    echo ">> Running..."
    ./bin/app
    ;;

  run-dev)
    echo ">>Running On Air ..."
    air
    echo ">> Ready....."
    ;;

  migration)
    if [ -z "$1" ]; then
      echo "Usage: ./run.sh migration <name>"
      exit 1
    fi
    echo ">> Creating migration $1..."
    migrate create -ext sql -dir cmd/migrate/migrations "$1"
    ;;
  
  migrate-up)
    echo ">> Migrating up..."
    go run cmd/migrate/main.go up
    ;;
  
  migrate-down)
    echo ">> Migrating down..."
    go run cmd/migrate/main.go down
    ;;

  fix-version)
    if [ -z "$1" ]; then
      echo "Usage: ./run.sh fix-version <version>"
      exit 1
    fi
    VERSION=$1
    echo ">> Forcing database version to $VERSION (cleaning dirty flag)..."
    mysql -u "$DB_USER" -p"$DB_PASSWORD" -h "$DB_HOST" -P "$DB_PORT" "$DB_NAME" \
      -e "UPDATE schema_migrations SET dirty = FALSE WHERE version = '$VERSION';"
    echo ">> Done."
    ;;
  
  *)
    echo "Usage: $0 {build|test|run|run-dev|migration|migrate-up|migrate-down|fix-version}"
    exit 1
    ;;
esac
