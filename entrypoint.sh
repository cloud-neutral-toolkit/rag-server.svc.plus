#!/usr/bin/env bash
set -euo pipefail

# -----------------------------------------------------------------------------
# Stunnel Configuration
# -----------------------------------------------------------------------------

USE_STUNNEL=0
STUNNEL_CONF="/etc/stunnel/stunnel.conf"

# 1. Try to load Stunnel config from file
if [ -f "/app/deploy/stunnel/rag-db-client.conf" ]; then
    mkdir -p /etc/stunnel
    cp "/app/deploy/stunnel/rag-db-client.conf" "${STUNNEL_CONF}"
    USE_STUNNEL=1
elif [ -f "deploy/stunnel/rag-db-client.conf" ]; then
    mkdir -p /etc/stunnel
    cp "deploy/stunnel/rag-db-client.conf" "${STUNNEL_CONF}"
    USE_STUNNEL=1
fi

# 2. Dynamic Override: if env vars are present, generate config
if [ -n "${DB_TLS_HOST:-}" ] && [ -n "${DB_TLS_PORT:-}" ]; then
  USE_STUNNEL=1
  echo "Configuring Stunnel from environment variables..."
  mkdir -p /etc/stunnel

  # Write certificate files from env vars if present
  if [ -n "${DB_CA:-}" ]; then
    echo "${DB_CA}" > /etc/stunnel/ca.pem
  fi

  # Generate stunnel.conf
  cat > "${STUNNEL_CONF}" <<EOF
foreground = no
pid = /var/run/stunnel/stunnel.pid
socket = l:TCP_NODELAY=1
socket = r:TCP_NODELAY=1
debug = 5
output = /var/log/stunnel.log

[postgres]
client = yes
accept = 127.0.0.1:5432
connect = ${DB_TLS_HOST}:${DB_TLS_PORT}
verify = 2
EOF

  if [ -f "/etc/stunnel/ca.pem" ]; then
    echo "CAfile = /etc/stunnel/ca.pem" >> "${STUNNEL_CONF}"
  elif [ -f "/etc/ssl/certs/ca-certificates.crt" ]; then
    echo "CAfile = /etc/ssl/certs/ca-certificates.crt" >> "${STUNNEL_CONF}"
  fi
fi

# 3. Start Stunnel if enabled
if [ "${USE_STUNNEL}" -eq 1 ]; then
  echo "Starting Stunnel..."
  mkdir -p /var/run/stunnel
  # Check if stunnel is installed
  if ! command -v stunnel &> /dev/null; then
      echo "Error: stunnel is not installed. Skipping Stunnel start."
      # Don't exit 1, maybe running locally without stunnel but config present?
      # But USE_STUNNEL=1 implies we need it.
      # Proceeding might fail connection.
  else
      stunnel "${STUNNEL_CONF}"
      
      # Wait for stunnel to be up (simple check)
      STUNNEL_UP=0
      for i in {1..30}; do
        if nc -z 127.0.0.1 5432; then
          echo "Stunnel is up!"
          STUNNEL_UP=1
          break
        fi
        echo "Waiting for Stunnel... ($i/30)"
        sleep 1
      done

      if [ "${STUNNEL_UP}" -eq 0 ]; then
        echo "Error: Stunnel failed to start or is not reachable on port 5432."
        exit 1
      fi
  fi
fi

# -----------------------------------------------------------------------------
# App Configuration
# -----------------------------------------------------------------------------

CONFIG_FILE="${CONFIG_PATH:-/etc/rag-server/rag-server.yaml}"
DEFAULT_CONFIG="/etc/rag-server/server.yaml"
mkdir -p "$(dirname "${CONFIG_FILE}")"

if [ ! -f "${CONFIG_FILE}" ]; then
  if [ -f "/app/config/rag-server.yaml" ]; then
      cp "/app/config/rag-server.yaml" "${CONFIG_FILE}"
  elif [ -f "/etc/rag-server/server.yaml" ]; then
      cp "/etc/rag-server/server.yaml" "${CONFIG_FILE}"
  elif [ -f "/app/config/server.yaml" ]; then
      cp "/app/config/server.yaml" "${CONFIG_FILE}"
  fi
fi

# DATABASE_URL override if Stunnel is used
if [ "$USE_STUNNEL" -eq 1 ]; then
    export DATABASE_URL="postgres://${DB_USER:-postgres}:${DB_PASSWORD:-password}@127.0.0.1:5432/${DB_NAME:-postgres}?sslmode=disable"
fi

exec /usr/local/bin/rag-server --config "${CONFIG_FILE}" "$@"
