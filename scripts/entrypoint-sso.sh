#!/usr/bin/env bash
set -e

# Проверяем, передан ли путь к конфигу через переменную окружения или флаг
CONFIG_PATH="${CONFIG_PATH:-./config/local_pg.yaml}"

if [ ! -f "$CONFIG_PATH" ]; then
  echo "Config file not found: $CONFIG_PATH"
  exit 1
fi

exec /app/sso --config="$CONFIG_PATH"
