# Configuration

The application uses `server.yaml` for base configuration but prioritizes Environment Variables—making it ideal for **Cloud Run** or K8s.

## Environment Variables

| Setting | Env Variable | Default | Description |
| :--- | :--- | :--- | :--- |
| **Port** | `PORT` | `8080` | HTTP listening port. |
| **Redis Address** | `REDIS_ADDR` | `127.0.0.1:6379` | Redis connection string. |
| **Redis Password** | `REDIS_PASSWORD` | - | Redis password. |
| **Database URL** | `DATABASE_URL` / `PG_URL` | - | Full Postgres connection URL (e.g., `postgres://user:pass@host:5432/db`). |
| **LLM Token** | `CHUTES_API_TOKEN` | - | API Token for the LLM provider. |
| **LLM Endpoint** | `CHUTES_API_URL` | - | Base URL for LLM chat completions. |
| **LLM Model** | `CHUTES_API_MODEL` | - | Model name to use (e.g., `deepseek-r1:8b`). |

## `server.yaml`

See [`config/server.yaml`](../config/server.yaml) for the complete schema including detailed chunking, embedding, and auth settings.
