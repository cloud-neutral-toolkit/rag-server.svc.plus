# Endpoints

## POST /api/askai

Generate an answer from the configured chat completion provider. If RAG is initialized, retrieved chunks are also returned.

Request:

```json
{ "question": "How do I deploy?" }
```

Response:

```json
{
  "answer": "...",
  "chunks": [
    {"repo":"...","path":"...","chunk_id":0,"content":"...","metadata":{}}
  ]
}
```

Errors: `500` with `{ "error": "...", "config": {"timeout": <sec>, "retries": <n>} }`.

## POST /api/rag/query

Hybrid retrieval only (no generation). The server currently returns top 5 chunks.

Request:

```json
{ "question": "backup policy" }
```

Response:

```json
{ "chunks": [ {"repo":"...","path":"...","chunk_id":0,"content":"...","metadata":{}} ] }
```

Errors:

- `400` on invalid JSON
- `500` on unexpected failures
- `4xx/5xx` propagated from embedding provider

## POST /api/rag/upsert

Upsert pre-embedded documents.

Request (shape only):

```json
{
  "docs": [
    {
      "repo": "...",
      "path": "...",
      "chunk_id": 0,
      "content": "...",
      "embedding": [0.1, 0.2],
      "metadata": {"heading": "Intro"},
      "content_sha": "<sha256>"
    }
  ]
}
```

Response:

```json
{ "rows": 123 }
```

Errors: `503` if the vector store is unavailable.

## POST /api/sync

Sync a Git repository to a local directory.

Request:

```json
{ "repo_url": "https://github.com/org/repo.git", "local_path": "/tmp/repo" }
```

Response:

```json
{ "status": "synced" }
```

## GET /api/users

Returns users from the service database.

## GET /api/nodes

Returns nodes from the service database.

## GET /api/admin/settings

Requires `X-User-Role: admin|operator`.

Response:

```json
{ "version": 3, "matrix": {"moduleA": {"admin": true, "user": false}} }
```

## POST /api/admin/settings

Requires `X-User-Role: admin|operator`.

Request:

```json
{ "version": 3, "matrix": {"moduleA": {"admin": true}} }
```

Response:

```json
{ "version": 4, "matrix": {"moduleA": {"admin": true}} }
```

Conflict:

- `409` with `{ "error": "admin settings version conflict", "version": <current>, "matrix": <current> }`
