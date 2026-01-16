#!/usr/bin/env bash
set -euo pipefail

DOMAIN="rag-server.svc.plus"
DATA_DIR="${CADDY_DATA_DIR:-/var/lib/caddy}"
EXPORT_DIR="/etc/ssl/exported/rag-server"

cert_path=$(find "${DATA_DIR}" -type f -name "${DOMAIN}.crt" 2>/dev/null | head -n 1 || true)
key_path=$(find "${DATA_DIR}" -type f -name "${DOMAIN}.key" 2>/dev/null | head -n 1 || true)

if [[ -z "${cert_path}" || -z "${key_path}" ]]; then
  echo "certificate for ${DOMAIN} not found under ${DATA_DIR}; 需先访问站点触发签发/或 data-dir 不同" >&2
  exit 1
fi

install -d -m 0755 "${EXPORT_DIR}"
install -m 0644 "${cert_path}" "${EXPORT_DIR}/fullchain.pem"
install -m 0600 "${key_path}" "${EXPORT_DIR}/privkey.pem"

systemctl reload stunnel-pg || true
