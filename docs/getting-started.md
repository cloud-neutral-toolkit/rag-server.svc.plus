# Getting Started

## Prerequisites

- **Go** 1.24 or higher
- **PostgreSQL** with `vector` and `zhparser` extensions enabled.
- **Redis**

## Local Development

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
    The server looks for `config/server.yaml`. You can override defaults using environment variables. See [Configuration](./configuration.md) for details.
