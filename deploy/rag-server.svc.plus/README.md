```bash
sudo apt-get update
sudo apt-get install -y caddy stunnel4
make build
sudo install -m 0755 xcontrol-server /usr/local/bin/rag-server
sudo install -d /etc/rag-server /etc/stunnel /etc/caddy
sudo install -m 0644 deploy/rag-server.svc.plus/config/server.systemd.yaml /etc/rag-server/server.yaml
sudo install -m 0644 deploy/rag-server.svc.plus/etc/stunnel/rag-db-client.conf /etc/stunnel/rag-db-client.conf
sudo install -m 0644 deploy/rag-server.svc.plus/etc/caddy/Caddyfile.systemd /etc/caddy/Caddyfile
sudo install -m 0644 deploy/rag-server.svc.plus/systemd/stunnel-pg.service /etc/systemd/system/stunnel-pg.service
sudo install -m 0644 deploy/rag-server.svc.plus/systemd/rag-server.service /etc/systemd/system/rag-server.service
sudo systemctl daemon-reload
sudo systemctl enable --now stunnel-pg rag-server caddy
ss -lntp | rg ':(443|8080|15432)'
curl -fsS https://rag-server.svc.plus/health
psql "postgres://user:pass@127.0.0.1:15432/rag?sslmode=disable" -c 'select 1;'
nc -vz postgresql.svc.plus 5432

docker compose -f deploy/rag-server.svc.plus/docker-compose.yaml up -d --build
docker compose -f deploy/rag-server.svc.plus/docker-compose.yaml exec rag-server ss -lntp | rg ':(8080|15432)'
curl -fsS https://rag-server.svc.plus/health
psql "postgres://user:pass@127.0.0.1:15432/rag?sslmode=disable" -c 'select 1;'
nc -vz postgresql.svc.plus 5432
```
