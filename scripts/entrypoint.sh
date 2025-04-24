#!/bin/sh
set -e

WORKDIR_APP="/app"
echo "Миграции"

migrate-up.sh

echo "Миграции применены, запускаем сервис"
exec "$WORKDIR_APP/shortener