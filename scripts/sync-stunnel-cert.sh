#!/usr/bin/env bash
set -euo pipefail

SRC_CERT="/etc/stunnel/certs/rag-server.svc.plus.pem"
SRC_KEY="/etc/stunnel/certs/rag-server.svc.plus.key"
DEST_HOST="db-gateway.svc.plus"
DEST_DIR="/etc/stunnel/certs"

rsync -az --chmod=600 "$SRC_CERT" "$SRC_KEY" "root@${DEST_HOST}:${DEST_DIR}/"
ssh "root@${DEST_HOST}" "systemctl reload stunnel-rag-server"
