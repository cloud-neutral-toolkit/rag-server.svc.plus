# RAG Server

The **XControl RAG Server** (`rag-server`) is a high-performance, modular backend designed to power Retrieval-Augmented Generation (RAG) applications. It functions as the core "Knowledge Engine" of the [console.svc.plus](https://svc.plus) platform.

## 🚀 Key Features

*   **Store & Retrieve**: Efficient vector storage and semantic search using **PostgreSQL + pgvector**.
*   **Knowledge Sync**: Git-Ops style knowledge management. Sync content directly from Git repositories.
*   **Model Agnostic**: Compatible with OpenAI API format (Chutes.ai, Ollama, etc.).
*   **Cloud Native**: Designed for serverless deployment (Google Cloud Run).
*   **Secure**: Integrated authentication middleware.

## 🛠 Tech Stack

- **Language**: Go 1.25+
- **Web Framework**: Gin
- **Database**: PostgreSQL 16 (pgvector extension required)
- **Cache**: PostgreSQL (hstore + unlogged cache table)
- **Authentication**: JWT / XControl Auth Service

## 📦 Getting Started

### Prerequisites

- **Go** 1.24 or higher
- **PostgreSQL** with `vector`, `pg_jieba`, and `hstore` extensions enabled.

### Local Development

1.  **Clone and Init**:
    ```bash
    make init
    ```

2.  **Database Setup**:
    Ensure you have a Postgres instance running. Then initialize the schema:
    ```bash
    # Update Makefile DB credentials if necessary or set via env vars
    make init-db
    ```

3.  **Run Locally**:
    ```bash
    # Run with hot-reload (requires air) or standard go run
    make dev
    ```

4.  **Configuration**:
    The server looks for `config/server.yaml`. You can override defaults using environment variables (see below).

## ⚙️ Configuration

The application uses `server.yaml` for base configuration but prioritizes Environment Variables—making it ideal for **Cloud Run** or K8s.

| Setting | Env Variable | Default | Description |
| :--- | :--- | :--- | :--- |
| **Port** | `PORT` | `8080` | HTTP listening port. |
| **Database URL** | `DATABASE_URL` / `PG_URL` | - | Full Postgres connection URL (e.g., `postgres://user:pass@host:5432/db`). |
| **LLM Token** | `CHUTES_API_TOKEN` | - | API Token for the LLM provider. |
| **LLM Endpoint** | `CHUTES_API_URL` | - | Base URL for LLM chat completions. |
| **LLM Model** | `CHUTES_API_MODEL` | - | Model name to use (e.g., `deepseek-r1:8b`). |

### `server.yaml` Example
See `config/server.yaml` for the complete schema including detailed chunking, embedding, and auth settings.

## ☁️ Deployment (Cloud Run)

This project includes a `Dockerfile` optimized for Cloud Run.

```bash
# Build (or use GitHub Actions provided in .github/workflows)
docker build -t rag-server .

# Run
docker run -p 8080:8080 -e PORT=8080 -e DATABASE_URL="..." rag-server
```

**Cloud Run Tips**:
- Map the Cloud SQL instance using the Cloud Run SQL connection.
- Set `DATABASE_URL` to the socket path or private IP.
- Mount secrets for API keys.

## 🔌 API Reference

## 📚 Documentation

Detailed documentation is available in the [`docs/`](./docs) directory:

*   [**Getting Started**](./docs/getting-started.md): Installation and local development guide.
*   [**Configuration**](./docs/configuration.md): Environment variables and `server.yaml` settings.
*   [**Deployment**](./docs/deployment.md): Docker and Cloud Run deployment instructions.
*   [**API Reference**](./docs/api-reference.md): details on RAG, Sync, and System endpoints.

## 📜 License

This project is licensed under the MIT License.
