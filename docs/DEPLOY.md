# Neon + GCP first deploy checklist

State backend: **GCS** (`PULUMI_BACKEND_URL=gs://…`). No Pulumi Cloud token.

## Already done (typical)

- Neon project + GitHub secrets `DATABASE_URL` / `DATABASE_URL_MIGRATE` / `AUTH_SECRET`
- GitHub vars `GCP_PROJECT_ID`, `GCP_REGION`, `PULUMI_STACK`, `CORS_ORIGIN`

## 1. One-time local bootstrap

```bash
# Browser login (does not print account emails)
GCP_PROJECT_ID=YOUR_GCP_PROJECT_ID ./infra/scripts/gcloud-login.sh

# Strong passphrase for the GCS stack (store in GitHub + GCP SM — see infra/SECRETS.md)
export PULUMI_CONFIG_PASSPHRASE='…strong random value…'
./infra/scripts/store-pulumi-passphrase.sh

# Creates/uses private state bucket, pulumi up (no Cloud Run yet),
# sets WIF GitHub vars, grants deploy SA access to the bucket
GCP_PROJECT_ID=YOUR_GCP_PROJECT_ID \
PULUMI_BACKEND_URL=gs://YOUR_STATE_BUCKET \
GITHUB_OWNER=YOUR_GH_USER_OR_ORG \
./infra/scripts/bootstrap-local.sh
```

If the passphrase was ever pasted into chat or logs, **rotate it** before relying on deploy (`pulumi stack change-secrets-provider passphrase` — details in [`infra/SECRETS.md`](../infra/SECRETS.md)).

## 2. Google OAuth (can wait until after first API deploy)

1. Google Cloud Console → OAuth client (Web)
2. Authorized JS origins: `http://localhost:5173` + production frontend origin
3. Client ID → GitHub `GOOGLE_OAUTH_CLIENT_ID` (+ `frontend/.env` `VITE_GOOGLE_CLIENT_ID`)

## 3. GitHub configuration

### Secrets

| Secret | Notes |
|--------|--------|
| `DATABASE_URL` | Neon pooled |
| `DATABASE_URL_MIGRATE` | Neon direct (non-pooler) |
| `AUTH_SECRET` | long random string |
| `PULUMI_CONFIG_PASSPHRASE` | GCS stack config passphrase (also keep recovery copy in GCP SM) |
| `GOOGLE_OAUTH_CLIENT_ID` | optional until SPA auth |
| `GEMINI_API_KEY` | optional until critic |

No `PULUMI_ACCESS_TOKEN` when using the GCS backend.

### Variables

| Variable | Notes |
|----------|--------|
| `GCP_PROJECT_ID` | Target project |
| `GCP_REGION` | e.g. `us-central1` |
| `PULUMI_STACK` | e.g. `dev` |
| `PULUMI_BACKEND_URL` | `gs://YOUR_STATE_BUCKET` |
| `GCP_WORKLOAD_IDENTITY_PROVIDER` | from bootstrap |
| `GCP_DEPLOY_SERVICE_ACCOUNT` | from bootstrap |
| `CORS_ORIGIN` | SPA origin(s); `http://localhost:5173` OK for local SPA → Cloud Run API |

## 4. Deploy

Push to `main` or run **Deploy to GCP**. Flow: image → `pulumi login gs://…` → `pulumi up` → sync secrets → goose migrate → Cloud Run revision.

Health: `GET https://YOUR_CLOUD_RUN_URL/api/health`

## 5. Frontend

```bash
cd frontend
cp .env.example .env
# VITE_API_BASE_URL=https://YOUR_CLOUD_RUN_URL
# VITE_GOOGLE_CLIENT_ID=...
npm run dev   # local
npm run build # static host later (Firebase / GCS+CDN)
```

## 6. Critic worker

Cloud Run sets `APP_ENV=production`, `CRITIC_WORKER_ENABLED=true`, and `minScale=1` so `@vibe_critic` polls `bot_jobs`.  
Standalone worker: `go run ./cmd/bot` with the same DB + `GEMINI_API_KEY`.
