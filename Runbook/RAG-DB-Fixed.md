
# RAG Server 数据库连接修复方案

## 问题诊断

### 当前错误
```
ERROR: cannot parse `admin_password`: failed to parse as keyword/value (invalid keyword/value)
WARN: postgres cache disabled; no database connection
```

### 根本原因
1. `DATABASE_URL` 环境变量指向 Secret Manager 引用 `admin_password`，而不是实际的数据库连接字符串
2. RAG 服务器无法连接到 PostgreSQL 数据库
3. 导致 `/api/rag/query` 返回 404

## 架构说明

```
Cloud Run (RAG Server)
  ↓
Stunnel (127.0.0.1:5432)
  ↓ TLS Tunnel
  ↓
postgresql.onwalk.net:443
  ↓
PostgreSQL Database
```

## 修复步骤

### 方案 1: 使用 Stunnel (推荐)

#### 1. 更新 Cloud Run 环境变量

```bash
gcloud run services update rag-server-svc-plus \
  --region asia-northeast1 \
  --project xzerolab-480008 \
  --set-env-vars="
DB_TLS_HOST=postgresql.onwalk.net,
DB_TLS_PORT=443,
POSTGRES_USER=postgres,
POSTGRES_PASSWORD=otdcRLTJamszk3AE,
POSTGRES_DB=knowledge_db,
NVIDIA_API_KEY=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"
```

**说明**:
- `DB_TLS_HOST` 和 `DB_TLS_PORT` 会触发 entrypoint.sh 启动 Stunnel
- Stunnel 会在 `127.0.0.1:5432` 监听
- `config/rag-server.yaml` 中的 `pgurl` 会自动使用 `127.0.0.1:5432`

#### 2. 验证配置

部署后检查日志：
```bash
gcloud logging read "resource.type=cloud_run_revision AND resource.labels.service_name=rag-server-svc-plus" \
  --limit 20 --project xzerolab-480008
```

应该看到：
```
Starting Stunnel...
Stunnel is up!
```

### 方案 2: 直接数据库连接 (不推荐，需要公网 IP)

如果 PostgreSQL 有公网 IP，可以直接连接：

```bash
# 更新 Secret Manager
echo "postgres://postgres:otdcRLTJamszk3AE@<PUBLIC_IP>:5432/knowledge_db?sslmode=require" | \
  gcloud secrets versions add DATABASE_URL --data-file=- --project xzerolab-480008
```

## 当前状态

### 工作的功能
- ✅ `/api/askai` - 直接 AI 回答（不依赖数据库）
- ✅ CORS 配置
- ✅ NVIDIA API 集成

### 待修复的功能
- ❌ `/api/rag/query` - RAG 检索（需要数据库连接）
- ❌ 向量搜索
- ❌ 知识库同步

## 前端行为

前端已经实现了优雅降级：
1. 首先尝试 `/api/rag/query` (RAG 检索)
2. 如果失败，回退到 `/api/askai` (直接 AI)
3. 用户仍然可以获得答案，只是没有知识库上下文

## 建议

**立即执行方案 1**，因为：
1. 安全性更好（TLS 加密）
2. 不需要暴露数据库公网 IP
3. 符合现有架构设计
4. entrypoint.sh 已经支持 Stunnel

## 部署命令

```bash
# 完整部署命令
gcloud run services update rag-server-svc-plus \
  --region asia-northeast1 \
  --project xzerolab-480008 \
  --update-env-vars="
DB_TLS_HOST=postgresql.onwalk.net,
DB_TLS_PORT=443,
POSTGRES_USER=postgres,
POSTGRES_PASSWORD=otdcRLTJamszk3AE,
POSTGRES_DB=knowledge_db,
NVIDIA_API_KEY=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" \
  --clear-env-vars=DATABASE_URL
```

注意：
- 使用 `--update-env-vars` 添加新变量
- 使用 `--clear-env-vars` 移除错误的 `DATABASE_URL`

