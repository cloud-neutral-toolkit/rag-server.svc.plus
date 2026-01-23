# Deployment (Cloud Run)

This project includes a `Dockerfile` optimized for Cloud Run.

## Build and Run

```bash
# Build (or use GitHub Actions provided in .github/workflows)
docker build -t rag-server .

# Run
docker run -p 8080:8080 -e PORT=8080 -e DATABASE_URL="..." rag-server
```

## Cloud Run Tips

- **Database Connection**: Map the Cloud SQL instance using the Cloud Run SQL connection feature.
- **Connection String**: Set `DATABASE_URL` to the socket path (e.g., `host=/cloudsql/project:region:instance user=...`) or private IP depending on your VPC configuration.
- **Secrets**: Mount sensitive values like `CHUTES_API_TOKEN` using Cloud Secret Manager.

For a more detailed guide, see [Google Cloud Run How-to](./google-cloud-run-howto.md).
