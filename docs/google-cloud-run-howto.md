# Cloud Run Deployment Setup Guide

This guide details the configuration required to enable the automated Cloud Run deployment workflow.

## 1. Google Cloud Configuration

### APIs
Ensure the following APIs are enabled in your Google Cloud Project:
- Cloud Run API
- Artifact Registry API
- IAM Service Account Credentials API

### Artifact Registry
Create a Docker repository in Artifact Registry:
- **Name**: `my-repo` (or update `REPOSITORY` in `deploy.yml`)
- **Region**: `asia-northeast1` (or update `REGION` in `deploy.yml`)
- **Format**: Docker

### Workload Identity Federation (WIF)
1. **Create a Workload Identity Pool**:
   ```bash
   gcloud iam workload-identity-pools create "my-pool" \
     --project="YOUR_PROJECT_ID" \
     --location="global" \
     --display-name="GitHub Actions Pool"
   ```

2. **Create a Workload Identity Provider**:
   ```bash
   gcloud iam workload-identity-pools providers create-oidc "my-provider" \
     --project="YOUR_PROJECT_ID" \
     --location="global" \
     --workload-identity-pool="my-pool" \
     --display-name="My GitHub Provider" \
     --attribute-mapping="google.subject=assertion.sub,attribute.actor=assertion.actor,attribute.repository=assertion.repository" \
     --issuer-uri="https://token.actions.githubusercontent.com"
   ```

3. **Create a Service Account**:
   Create a service account (e.g., `my-service-account`) and grant it the following roles:
   - **Cloud Run Admin** (`roles/run.admin`): To deploy services.
   - **Artifact Registry Writer** (`roles/artifactregistry.writer`): To push images.
   - **Service Account User** (`roles/iam.serviceAccountUser`): To act as the service account.

4. **Bind Service Account to WIF**:
   ```bash
   gcloud iam service-accounts add-iam-policy-binding "my-service-account@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
     --project="YOUR_PROJECT_ID" \
     --role="roles/iam.workloadIdentityUser" \
     --member="principalSet://iam.googleapis.com/projects/YOUR_PROJECT_NUMBER/locations/global/workloadIdentityPools/my-pool/attribute.repository/YOUR_GITHUB_USERNAME/REPO_NAME"
   ```

## 2. GitHub Configuration

### Secrets
Go to **Settings > Secrets and variables > Actions** in your GitHub repository and add the following secrets:

| Secret Name | Description | Example Value |
|-------------|-------------|---------------|
| `WIF_PROVIDER` | The full resource name of the WIF provider. | `projects/123456789/locations/global/workloadIdentityPools/my-pool/providers/my-provider` |
| `WIF_SERVICE_ACCOUNT` | The email of the service account. | `my-service-account@your-project-id.iam.gserviceaccount.com` |

### Workflow Environment Variables
Edit `.github/workflows/deploy.yml` and update the `env` section with your specific values:

```yaml
env:
  PROJECT_ID: your-project-id # RELACE THIS
  REGION: asia-northeast1     # Modify if using a different region
  SERVICE_NAME: rag-server    # The name of your Cloud Run service
  REPOSITORY: my-repo         # The name of your Artifact Registry repo
```

## 3. Verify Deployment
Push a change to the `main` branch. The "Build and Deploy to Cloud Run" workflow should trigger automatically.
