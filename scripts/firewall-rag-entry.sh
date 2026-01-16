#!/usr/bin/env bash
set -euo pipefail

DB_GATEWAY_IP="203.0.113.10"
STUNNEL_PORT="8443"

if command -v ufw >/dev/null 2>&1; then
  ufw allow 443/tcp
  ufw allow out "${STUNNEL_PORT}"/tcp to "${DB_GATEWAY_IP}"
  ufw deny in "${STUNNEL_PORT}"/tcp
  exit 0
fi

if command -v firewall-cmd >/dev/null 2>&1; then
  firewall-cmd --permanent --add-service=https
  firewall-cmd --permanent --add-rich-rule="rule family=ipv4 destination address=${DB_GATEWAY_IP} port protocol=tcp port=${STUNNEL_PORT} accept"
  firewall-cmd --permanent --add-rich-rule="rule family=ipv4 port protocol=tcp port=${STUNNEL_PORT} drop"
  firewall-cmd --reload
  exit 0
fi

iptables -A INPUT -p tcp --dport 443 -j ACCEPT
iptables -A OUTPUT -p tcp -d "${DB_GATEWAY_IP}" --dport "${STUNNEL_PORT}" -j ACCEPT
iptables -A INPUT -p tcp --dport "${STUNNEL_PORT}" -j DROP
