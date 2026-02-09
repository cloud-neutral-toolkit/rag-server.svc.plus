# rag-server.svc.plus

生产就绪的 RAG 后端服务：文档同步、向量检索、全文检索、问答生成 (Retrieval-Augmented Generation).

> A production-oriented RAG backend for ingestion, hybrid retrieval (vector + full-text), and answer generation.

## 部署要求 (Deployment Requirements)

| 维度 | 要求 / 规格 | 说明 |
|---|---|---|
| 网络 | API 可访问 | 默认 `:8080` |
| 端口 | `:8080` | HTTP API |
| 数据库 | PostgreSQL + pgvector + zhparser | `knowledge_db`（可配置） |
| 最低 | 1 CPU / 2GB RAM | 开发/小规模 |
| 推荐 | 2 CPU / 4GB RAM | 生产建议 |

## 快速开始 (Quickstart)

### 一键初始化 (Setup Script)

```bash
curl -fsSL "https://raw.githubusercontent.com/cloud-neutral-toolkit/rag-server.svc.plus/main/scripts/setup.sh?$(date +%s)" \
  | bash -s -- rag-server.svc.plus
```

### 本地运行 (Local Dev)

```bash
cp .env.example .env
make init-db
make dev
```

## 核心特性 & 技术栈 (Features & Tech Stack)

核心特性：
- 文档同步与入库：从 Git 仓库同步 markdown，chunk + embedding + upsert
- 混合检索：pgvector 相似度 + PostgreSQL 全文检索（含中文分词）
- 模型可插拔：OpenAI-compatible embeddings/chat-completions API（也可接入其他 provider）
- 可运维：单体无状态服务，后端依赖主要是 PostgreSQL

技术栈：
- Go
- PostgreSQL + pgvector + zhparser
- 可选：Cloud Run 部署示例（见 `deploy/` 与 `docs/`）

## 说明文档 (Docs)

- 文档入口：`docs/README.md`
- 快速开始：`docs/getting-started/quickstart.md`
- 配置：`docs/configuration.md` / `config/rag-server.yaml`
- API：`docs/api/overview.md`, `docs/api/endpoints.md`
- Token/JWT：`docs/TOKEN_AUTH_MANUAL.md`
- 部署：`docs/deployment.md`, `docs/google-cloud-run-howto.md`
- Runbooks：`docs/Runbook/`

