# Neon + GCP first deploy checklist

## 1. Neon

1. Create a Neon project + database.
2. Copy two connection strings:
   - **Pooled** → GitHub secret `DATABASE_URL` (and GCP Secret Manager via deploy sync)
   - **Direct** (no `-pooler` in host) → GitHub secret `DATABASE_URL_MIGRATE`

## 2. Google OAuth

1. Google Cloud Console → OAuth client (Web)
2. Authorized JS origins: local Vite (`http://localhost:5173`) + production frontend origin
3. Client ID → GitHub `GOOGLE_OAUTH_CLIENT_ID`, API secret sync, and `frontend/.env` `VITE_GOOGLE_CLIENT_ID`

## 3. Local Pulumi bootstrap (once)

```bash
cd infra
cp Pulumi.dev.yaml.example Pulumi.dev.yaml
pulumi stack init dev
pulumi config set gcp:project YOUR_PROJECT_ID
pulumi config set vibegopher:githubOwner YOUR_GH_USER_OR_ORG
pulumi config set vibegopher:githubRepo vibegopher
pulumi config set vibegopher:enableCloudRun false
# When you later enable Cloud Run, also set:
# pulumi config set vibegopher:corsOrigin https://your-frontend-origin
pulumi up
```

Copy stack outputs into GitHub **variables**:

- `GCP_WORKLOAD_IDENTITY_PROVIDER`
- `GCP_DEPLOY_SERVICE_ACCOUNT`
- also set `GCP_PROJECT_ID`, `GCP_REGION`, `PULUMI_STACK`, `CORS_ORIGIN`, `CORS_ORIGIN`

## 4. GitHub secrets

| Secret | Notes |
|--------|--------|
| `PULUMI_ACCESS_TOKEN` | Pulumi Cloud |
| `DATABASE_URL` | Neon pooled |
| `DATABASE_URL_MIGRATE` | Neon direct |
| `AUTH_SECRET` | long random string |
| `GOOGLE_OAUTH_CLIENT_ID` | Web client ID |
| `GEMINI_API_KEY` | AI critic |

## 5. Deploy

Push to `main` or run **Deploy to GCP**. Flow: image → `pulumi up` → sync secrets → goose migrate → Cloud Run revision.

Health: `GET https://YOUR_CLOUD_RUN_URL/api/health`

## 6. Frontend

```bash
cd frontend
cp .env.example .env
# VITE_API_BASE_URL=https://YOUR_CLOUD_RUN_URL
# VITE_GOOGLE_CLIENT_ID=...
npm run dev   # local
npm run build # static host later (Firebase / GCS+CDN)
```

## 7. Critic worker

Cloud Run sets `APP_ENV=production`, `CRITIC_WORKER_ENABLED=true`, and `minScale=1` so `@vibe_critic` polls `bot_jobs`.  
Set repository variable `CORS_ORIGIN` to your SPA origin before deploy (required).  
Standalone worker: `go run ./cmd/bot` with the same DB + `GEMINI_API_KEY`.
