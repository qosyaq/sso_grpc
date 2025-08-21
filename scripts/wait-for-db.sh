#!/usr/bin/env bash
# usage: wait-for-db.sh <host> <port> [timeout_seconds]
set -e

HOST=${1:-db}
PORT=${2:-5432}
TIMEOUT=${3:-60}
START_TS=$(date +%s)

echo "Waiting for postgres at ${HOST}:${PORT} (timeout ${TIMEOUT}s)..."

while true; do
  if pg_isready -h "${HOST}" -p "${PORT}" -U "sso_user" >/dev/null 2>&1; then
    echo "Postgres is ready."
    exit 0
  fi
  NOW_TS=$(date +%s)
  ELAPSED=$((NOW_TS - START_TS))
  if [ "${ELAPSED}" -ge "${TIMEOUT}" ]; then
    echo "Timed out waiting for Postgres (${ELAPSED}s >= ${TIMEOUT}s)." >&2
    exit 1
  fi
  sleep 1
done
