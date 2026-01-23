# API Reference

## RAG & AI

### Ask AI
`POST /api/askai`

Ask a question. Returns both the LLM-generated answer and the source chunks used for context.

**Body:**
```json
{
  "question": "How to configure XControl?"
}
```

### Semantic Search
`POST /api/rag/query`

Perform a semantic search against the knowledge base without generating an answer. Returns relevant document chunks.

**Body:**
```json
{
  "question": "backup policy"
}
```

### Upsert Documents
`POST /api/rag/upsert`

Manually index documents into the vector store.

## Knowledge Sync

### Sync Repository
`POST /api/sync`

Trigger a synchronization process from a remote Git repository.

**Body:**
```json
{
  "repo_url": "https://github.com/...",
  "local_path": "docs"
}
```

## System

### Health Checks
`GET /health`
`GET /healthz`

Returns 200 OK if the service is running. Use these for liveness and readiness probes.
