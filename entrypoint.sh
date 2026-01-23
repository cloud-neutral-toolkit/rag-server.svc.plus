#!/usr/bin/env bash
set -euo pipefail

CONFIG_FILE="${CONFIG_PATH:-/etc/rag-server/server.yaml}"
DEFAULT_CONFIG="/etc/rag-server/server.yaml"
mkdir -p "$(dirname "${CONFIG_FILE}")"

if [ ! -f "${CONFIG_FILE}" ]; then
  cp "${DEFAULT_CONFIG}" "${CONFIG_FILE}"
fi

# Port configuration is now handled natively by the application via PORT env var
# and other settings via REDIS_ADDR, DATABASE_URL etc.
# escaping fragile awk logic.

exec /usr/local/bin/rag-server --config "${CONFIG_FILE}" "$@"
