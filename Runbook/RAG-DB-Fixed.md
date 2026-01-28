# RAG Server 数据库连接修复 Runbook

## 📋 概述

**问题**: RAG Server 无法连接到 PostgreSQL 数据库，导致 `/api/rag/query` 返回 404  
**影响范围**: RAG 检索功能、向量搜索、知识库同步  
**修复时间**: ~5 分钟  
**风险等级**: 🟡 中等（前端已实现降级，用户体验影响有限）

---

## 🔍 问题诊断

### 错误症状

```
ERROR: cannot parse `admin_password`: failed to parse as keyword/value (invalid keyword/value)
WARN: postgres cache disabled; no database connection
```

### 根本原因

1. ❌ `DATABASE_URL` 环境变量指向 Secret Manager 引用 `admin_password`，而不是实际的数据库连接字符串
2. ❌ RAG 服务器无法连接到 PostgreSQL 数据库
3. ❌ 导致 `/api/rag/query` 返回 404

### 架构图

```
┌─────────────────────┐
│  Cloud Run          │
│  (RAG Server)       │
└──────────┬──────────┘
           │
           ↓
┌─────────────────────┐
│  Stunnel            │
│  127.0.0.1:5432     │
└──────────┬──────────┘
           │ TLS Tunnel
           ↓
┌─────────────────────┐
│  postgresql         │
│  .svc.plus:443       │
└──────────┬──────────┘
           │
           ↓
┌─────────────────────┐
│  PostgreSQL         │
│  Database           │
└─────────────────────┘
```

---

## ✅ 修复方案

### 方案 1: 使用 Stunnel (推荐) ⭐

#### 优势
- ✅ 安全性更好（TLS 加密）
- ✅ 不需要暴露数据库公网 IP
- ✅ 符合现有架构设计
- ✅ entrypoint.sh 已经支持 Stunnel

#### 步骤 1: 更新 Cloud Run 环境变量

```bash
gcloud run services update rag-server-svc-plus \
  --region asia-northeast1 \
  --project xzerolab-480008 \
  --update-env-vars="\
DB_TLS_HOST=postgresql.svc.plus,\
DB_TLS_PORT=443,\
POSTGRES_USER=postgres,\
POSTGRES_PASSWORD=otdcRLTJamszk3AE,\
POSTGRES_DB=knowledge_db,\
NVIDIA_API_KEY=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" \
  --clear-env-vars=DATABASE_URL
```

**关键参数说明**:
- `DB_TLS_HOST` 和 `DB_TLS_PORT`: 触发 entrypoint.sh 启动 Stunnel
- Stunnel 会在 `127.0.0.1:5432` 监听
- `config/rag-server.yaml` 中的 `pgurl` 会自动使用 `127.0.0.1:5432`
- `--clear-env-vars=DATABASE_URL`: 移除错误的环境变量

#### 步骤 2: 验证部署

**检查日志**:
```bash
gcloud logging read \
  "resource.type=cloud_run_revision AND resource.labels.service_name=rag-server-svc-plus" \
  --limit 20 \
  --project xzerolab-480008 \
  --format="table(timestamp,textPayload)"
```

**期望输出**:
```
Starting Stunnel...
Stunnel is up!
INFO: Connected to PostgreSQL at 127.0.0.1:5432
```

#### 步骤 3: 功能测试

**测试 RAG 查询**:
```bash
curl -X POST https://rag-server-svc-plus-HASH-an.a.run.app/api/rag/query \
  -H "Content-Type: application/json" \
  -d '{"query": "test query", "top_k": 5}'
```

**期望响应**: HTTP 200，包含检索结果

---

### 方案 2: 直接数据库连接 (不推荐)

⚠️ **仅在 PostgreSQL 有公网 IP 且无法使用 Stunnel 时使用**

```bash
# 更新 Secret Manager
echo "postgres://postgres:otdcRLTJamszk3AE@<PUBLIC_IP>:5432/knowledge_db?sslmode=require" | \
  gcloud secrets versions add DATABASE_URL --data-file=- --project xzerolab-480008

# 更新 Cloud Run 使用 Secret
gcloud run services update rag-server-svc-plus \
  --region asia-northeast1 \
  --project xzerolab-480008 \
  --update-secrets=DATABASE_URL=DATABASE_URL:latest
```

---

## 📊 当前状态

### ✅ 工作的功能
- `/api/askai` - 直接 AI 回答（不依赖数据库）
- CORS 配置
- NVIDIA API 集成

### ❌ 待修复的功能
- `/api/rag/query` - RAG 检索（需要数据库连接）
- 向量搜索
- 知识库同步

---

## 🎯 前端降级策略

前端已经实现了优雅降级：

1. **首选**: 尝试 `/api/rag/query` (RAG 检索)
2. **降级**: 如果失败，回退到 `/api/askai` (直接 AI)
3. **结果**: 用户仍然可以获得答案，只是没有知识库上下文

---

## 🔄 回滚计划

如果修复后出现问题，执行以下回滚：

```bash
# 恢复原有 DATABASE_URL
gcloud run services update rag-server-svc-plus \
  --region asia-northeast1 \
  --project xzerolab-480008 \
  --update-secrets=DATABASE_URL=admin_password:latest \
  --remove-env-vars=DB_TLS_HOST,DB_TLS_PORT,POSTGRES_USER,POSTGRES_PASSWORD,POSTGRES_DB
```

---

## 📝 验证清单

- [ ] Cloud Run 环境变量已更新
- [ ] `DATABASE_URL` 已清除
- [ ] 日志显示 "Stunnel is up!"
- [ ] `/api/rag/query` 返回 200
- [ ] 前端 RAG 功能正常
- [ ] 向量搜索可用

---

## 📚 相关文档

- [Stunnel 配置文档](../docs/stunnel-setup.md)
- [PostgreSQL 连接指南](../docs/postgres-connection.md)
- [Cloud Run 环境变量管理](https://cloud.google.com/run/docs/configuring/environment-variables)

---

## 🆘 故障排查

### 问题: Stunnel 无法启动

**检查**:
```bash
# 查看 Stunnel 日志
gcloud logging read "textPayload=~\"stunnel\"" --limit 50 --project xzerolab-480008
```

**可能原因**:
- DNS 解析失败
- 防火墙阻止 443 端口
- 证书验证失败

### 问题: 数据库连接超时

**检查**:
```bash
# 测试网络连通性
gcloud run services describe rag-server-svc-plus \
  --region asia-northeast1 \
  --project xzerolab-480008 \
  --format="value(status.url)"
```

**可能原因**:
- PostgreSQL 服务未运行
- 密码错误
- 数据库名称错误

---

**最后更新**: 2026-01-26  
**负责人**: DevOps Team  
**审核人**: Tech Lead
