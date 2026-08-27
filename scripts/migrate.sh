#!/bin/sh
set -e

DATABASE_URL="postgres://${DB__USER}:${DB__PASSWORD}@${DB__HOST}:${DB__PORT}/${DB__NAME}?sslmode=${DB__SSLMODE}"

exec /app/migrator --postgres-dsn="${DATABASE_URL}" --migrations-path=/app/migrations "$@"