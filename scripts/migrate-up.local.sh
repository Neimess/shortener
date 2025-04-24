#!/bin/bash
set -e

WORKDIR="$WORKDIR"

MIGRATIONS="$WORKDIR/migrations"

source "$WORKDIR/.env.local"

case "$DRIVER" in
  postgres)
    DATABASE_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"
    ;;
  sqlite)
    DATABASE_URL="sqlite3://$SQLITE_FILE"
    ;;
  *)
    echo "Unsupported DRIVER: $DRIVER"
    exit 1
    ;;
esac

echo "Running migrations for $DRIVER"
migrate -path "$MIGRATIONS" -database "$DATABASE_URL" up

