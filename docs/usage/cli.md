# CLI

## rag-server

Starts the HTTP API server.

```
rag-server [--config <path>] [--log-level <debug|info|warn|error>]
```

- `--config`: path to YAML configuration file.
- `--log-level`: overrides `log.level` from config.

## rag-cli

Synchronizes repositories and ingests markdown files into the server via API.

```
rag-cli --config <path> [--file <path>] [--log-level <level>]
```

Environment:

- `SERVER_URL`: base URL for the API (default: from config, or `http://localhost:8080`).

Behavior:

- If `--file` is provided, only that file is embedded and upserted.
- Otherwise, it iterates over `global.datasources` in config and ingests each repo.

## ingest (batch tool)

Direct ingestion into Postgres (no HTTP) using the RAG config:

```
ingest --config <path> [--only-repo <name>] [--dry-run] [--max-files <n>] [--migrate-dim] [--concurrency <n>]
```

This tool uses the same chunking and embedding config as the server.

## ragbench (optional)

Benchmark runner for retrieval quality. It expects a YAML input file describing queries.

```
ragbench --in queries.yaml --api http://localhost:8080 --out report.md
```

Note: `ragbench` assumes a response shape that is not identical to the current `/api/rag/query` response. Use it with a compatible adapter or update the tool before relying on the metrics.
