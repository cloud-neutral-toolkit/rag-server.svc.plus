# Configuration

The server reads a YAML config file and applies environment variable overrides. The same file is used by `rag-server` and `rag-cli`.

## Config resolution order

1. `--config <path>` flag
2. `config/rag-server.yaml` if present
3. `config/server.yaml`

> Note: the repo includes `config/rag-server.yaml` which uses legacy keys (for example `vector_db`). The **current** schema uses `global.vectordb` and `server.addr`. Use `example/config/server.yaml` as the canonical reference.

## Minimal example

```yaml
server:
  addr: ":8080"
  readTimeout: 30s
  writeTimeout: 30s
  publicUrl: "https://rag-server.svc.plus"
  allowedOrigins:
    - "https://rag-server.svc.plus"

auth:
  enable: false
  authUrl: "https://accounts.svc.plus"
  apiBaseUrl: "https://api.svc.plus"
  publicToken: ""

global:
  cache:
    table: "cache_kv"
    defaultTTL: 5m
  vectordb:
    pgurl: "postgres://user:pass@127.0.0.1:5432/knowledge_db?sslmode=disable"
  proxy: ""
  datasources:
    - name: knowledge
      repo: https://github.com/example/knowledge.git
      path: /

sync:
  repo:
    proxy: "socks5://127.0.0.1:1080"

models:
  embedder:
    provider: "chutes"     # or "ollama"
    models: ["bge-m3"]
    endpoint: "http://127.0.0.1:9000/v1/embeddings"
    token: "${NVIDIA_API_KEY}"
  generator:
    provider: "chutes"
    models: ["deepseek-r1:8b"]
    endpoint: "http://127.0.0.1:9000/v1/chat/completions"
    token: "${NVIDIA_API_KEY}"

embedding:
  max_batch: 64
  dimension: 1024
  max_chars: 8000
  rate_limit_tpm: 120000

chunking:
  max_tokens: 800
  overlap_tokens: 80
  include_exts: [".md", ".mdx"]
  ignore_dirs: [".git", "node_modules", "dist", "build"]

retrieval:
  alpha: 0.5
  candidates: 50

api:
  askai:
    timeout: 100
    retries: 3
```

## Key sections

### server

- `addr`: listen address (default `:8080` if empty).
- `readTimeout`, `writeTimeout`: HTTP timeouts.
- `publicUrl`: used to derive default CORS origin if none configured.
- `allowedOrigins`: CORS allowlist. Use `*` to allow all origins (credentials disabled).

### auth

- `enable`: toggle auth middleware.
- `authUrl`, `apiBaseUrl`: currently informational; middleware does not call external services.
- `publicToken`: passed into the auth client (JWT verification uses internal secrets).

### global

- `cache`: Postgres-backed cache table used for token caching.
- `vectordb`: Postgres connection info (either `pgurl` or discrete fields).
- `proxy`: optional outbound proxy for HTTP and Git operations.
- `datasources`: list of Git repos for ingestion.

### models

- `embedder`: embedding provider used by RAG queries and ingestion.
- `generator`: chat completion provider used by `/api/askai`.
- `reranker` (optional, in RAG config): supports reranking if `endpoint` is set.

### embedding / chunking

Tune chunk size, overlap, and embedding limits. See `internal/rag/config` for full options (including `embed_toc`, `embed_headings`, `by_paragraph`, and `additional_max_tokens`).

### retrieval

- `alpha`: blend between vector and text scores (0..1).
- `candidates`: number of candidates retrieved before reranking.

## Environment variables

- `PORT`: overrides server listen port.
- `DATABASE_URL` or `PG_URL`: overrides Postgres DSN.
- `SERVER_URL`: base URL for `rag-cli` when not provided in config.
- `CHUTES_API_URL`, `CHUTES_API_MODEL`, `CHUTES_API_TOKEN`: override AskAI model settings for OpenAI-compatible providers.

Environment variables are expanded inside the YAML file (for example `${NVIDIA_API_KEY}`).
