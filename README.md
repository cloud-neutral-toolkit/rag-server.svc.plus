# XControl RAG Server

The **XControl RAG Server** (`rag-server`) is a high-performance, modular backend designed to power Retrieval-Augmented Generation (RAG) applications. It functions as the core "Knowledge Engine" of the [XControl](https://svc.plus) platform.

## 🚀 Key Features

*   **Store & Retrieve**: Efficient vector storage and semantic search using **PostgreSQL + pgvector**.
*   **Knowledge Sync**: Git-Ops style knowledge management. Sync content directly from Git repositories.
*   **Model Agnostic**: Compatible with OpenAI API format (Chutes.ai, Ollama, etc.).
*   **Cloud Native**: Designed for serverless deployment (Google Cloud Run).
*   **Secure**: Integrated authentication middleware.

## 🛠 Tech Stack

*   **Language**: Go 1.25+
*   **Framework**: Gin
*   **Database**: PostgreSQL 16 (pgvector extension)
*   **Cache**: Redis

## 📚 Documentation

Detailed documentation is available in the [`docs/`](./docs) directory:

*   [**Getting Started**](./docs/getting-started.md): Installation and local development guide.
*   [**Configuration**](./docs/configuration.md): Environment variables and `server.yaml` settings.
*   [**Deployment**](./docs/deployment.md): Docker and Cloud Run deployment instructions.
*   [**API Reference**](./docs/api-reference.md): details on RAG, Sync, and System endpoints.

## 📜 License

This project is licensed under the MIT License.
