#!/bin/sh
set -e

WORKDIR_APP="/app"

MIGRATIONS="$WORKDIR_APP/migrations"
DB_DRIVER=${DB_DRIVER}

case "$DB_DRIVER" in
  postgres)
    DATABASE_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"
    ;;
  sqlite)
    DATABASE_URL="sqlite3://$SQLITE_FILE"
    ;;
  *)
    echo "Unsupported DB_DRIVER: $DB_DRIVER"
    exit 1
    ;;
esac

echo "Running migrations for $DB_DRIVER"
migrate -path "$MIGRATIONS" -database "$DATABASE_URL" up

echo "Миграции применены, запускаем сервис…"
exec "$WORKDIR_APP/shortener"
