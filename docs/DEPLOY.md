# Template deploy checklist (Neon + GCP + Pulumi)

This public repo is a **forkable template**: Go REST API, Vue 3 SPA, Neon Postgres, and Pulumi → GCP (Cloud Run, Artifact Registry, Secret Manager, GitHub Actions WIF) with state in a **private GCS bucket**.

A production Slack/auth product should be a **private fork or derivative** — do not add Slack SaaS code here. Auth on `main` is **Google Identity Services → `/api/auth/google` → app JWT**. [Clerk PR #14](https://github.com/dakaii/vibegopher/pull/14) is optional/unmerged reference only.

State backend: **GCS** (`PULUMI_BACKEND_URL=gs://…`). No Pulumi Cloud token.

---

## What you get vs what you bring

| Included in this template | You must bring |
|---------------------------|----------------|
| App source (API + Vue SPA + critic worker) | A **GCP project** with billing |
| Pulumi Go program under `infra/` | A **Neon** project (pooled + direct connection strings) |
| GitHub Actions: Deploy / Destroy / tests / Docker publish | GitHub repo **secrets** + **variables** (tables below) |
| Bootstrap scripts (`gcloud-login`, `bootstrap-local`, passphrase store, sync-secrets) | One-time **local** `gcloud` login + bootstrap (WIF cannot create itself from CI) |
| Google GIS auth path on `main` | Google OAuth **Web client ID** (when you want SPA sign-in) |
| Optional Gemini critic | `GEMINI_API_KEY` when you enable the bot |

Not managed by Pulumi/destroy: **Neon**, the **GCS state bucket** itself, and **GitHub** secrets/vars.

---

## 0. Prerequisites on your laptop

- [gcloud CLI](https://cloud.google.com/sdk/docs/install) + a GCP account that can create APIs/IAM in the target project
- [Pulumi CLI](https://www.pulumi.com/docs/install/)
- Go **1.24+** (bootstrap / local `pulumi up` compile the infra program)
- [GitHub CLI](https://cli.github.com/) (`gh`) authenticated to the fork (`gh auth login`)
- Docker (local app) and Make

---

## 1. Neon

1. Create a Neon project (or use an existing one).
2. Copy:
   - **Pooled** connection string → GitHub secret `DATABASE_URL` (and later GCP SM)
   - **Direct** (non-pooler) connection string → GitHub secret `DATABASE_URL_MIGRATE` (goose; session/advisory locks)
3. Schema is applied by the Deploy workflow (`go run ./cmd/migrate up`), not by Pulumi.

Neon is **never** destroyed by the Destroy workflow.

---

## 2. One-time local bootstrap (GCP + WIF)

```bash
# Browser login (does not print account emails)
GCP_PROJECT_ID=YOUR_GCP_PROJECT_ID ./infra/scripts/gcloud-login.sh

# Strong passphrase for the GCS stack (store in GitHub + GCP SM — see infra/SECRETS.md)
export PULUMI_CONFIG_PASSPHRASE='…strong random value…'
./infra/scripts/store-pulumi-passphrase.sh

# Creates/uses private state bucket, pulumi up (Cloud Run off),
# sets WIF GitHub vars, grants deploy SA access to the bucket
GCP_PROJECT_ID=YOUR_GCP_PROJECT_ID \
PULUMI_BACKEND_URL=gs://YOUR_STATE_BUCKET \
GITHUB_OWNER=YOUR_GH_USER_OR_ORG \
./infra/scripts/bootstrap-local.sh
```

If the passphrase was ever pasted into chat or logs, **rotate it** before relying on deploy (`pulumi stack change-secrets-provider passphrase` — details in [`infra/SECRETS.md`](../infra/SECRETS.md)).

After bootstrap, confirm GitHub variables exist: `PULUMI_BACKEND_URL`, `GCP_WORKLOAD_IDENTITY_PROVIDER`, `GCP_DEPLOY_SERVICE_ACCOUNT`.

---

## 3. Google OAuth (can wait until after first API deploy)

Auth on `main`: Google Identity Services in the SPA → `POST /api/auth/google` → HS256 app JWT (`AUTH_SECRET`).

1. Google Cloud Console → APIs & Services → Credentials → OAuth client (Web)
2. Authorized JavaScript origins: `http://localhost:5173` + your production frontend origin
3. Client ID → GitHub secret `GOOGLE_OAUTH_CLIENT_ID` and `frontend/.env` as `VITE_GOOGLE_CLIENT_ID`

Password auth (`ENABLE_PASSWORD_AUTH=true`, `/api/signup` + `/api/login`) is for tests/local only and is **forbidden in production**.

Optional later: swap to Clerk (or another IdP) in a **private** product repo; see open [PR #14](https://github.com/dakaii/vibegopher/pull/14) as a design reference — not required for this template.

---

## 4. GitHub configuration

### Secrets

| Secret | Notes |
|--------|--------|
| `DATABASE_URL` | Neon pooled |
| `DATABASE_URL_MIGRATE` | Neon direct (non-pooler) |
| `AUTH_SECRET` | long random string (app JWT signing) |
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

Full placement rules: [`infra/SECRETS.md`](../infra/SECRETS.md).

---

## 5. First deploy

Push to `main` or run **Actions → Deploy to GCP** (`workflow_dispatch`).

Flow: build/push image → sync Secret Manager payloads → `pulumi login gs://…` → `pulumi up` (Cloud Run on) → goose migrate → Cloud Run revision bump.

Health: `GET https://YOUR_CLOUD_RUN_URL/api/health`

CLI:

```bash
gh workflow run "Deploy to GCP"
gh run watch
```

---

## 6. Frontend

```bash
cd frontend
cp .env.example .env
# VITE_API_BASE_URL=              # empty locally (Vite proxy); set to Cloud Run URL for prod builds
# VITE_GOOGLE_CLIENT_ID=....apps.googleusercontent.com
npm run dev   # local → http://localhost:5173
npm run build # static host later (Firebase / GCS+CDN)
```

**CORS tip:** leave `VITE_API_BASE_URL` empty for local Vite so the SPA calls relative `/api` through the proxy. Pointing the browser at Cloud Run from `localhost` requires `CORS_ORIGIN` to include `http://localhost:5173`.

---

## 7. Critic worker

Cloud Run sets `APP_ENV=production`, `CRITIC_WORKER_ENABLED=true`, and `minScale=1` so `@vibe_critic` polls `bot_jobs`.  
Standalone worker: `go run ./cmd/bot` with the same DB + `GEMINI_API_KEY`.

---

## 8. Destroy / teardown

Workflow: **Actions → Destroy GCP infrastructure** (`.github/workflows/destroy.yml`).

| Input | Meaning |
|-------|---------|
| `confirm` | Must be exactly `destroy` |
| `destroy_secrets` | Default **`false`** — keeps protected Secret Manager; set `true` only for a full SM wipe of this stack |

```bash
# Default: tear down Cloud Run / AR wiring / SAs / WIF; keep Secret Manager
gh workflow run "Destroy GCP infrastructure" -f confirm=destroy

# Full template reset of stack secrets (still does NOT touch Neon or the state bucket)
gh workflow run "Destroy GCP infrastructure" -f confirm=destroy -f destroy_secrets=true

gh run watch
```

### What remains after a default destroy

| Still there | Why |
|-------------|-----|
| Neon project + data | Out of scope; delete in Neon if you want |
| GCP Secret Manager secrets + related IAM / SAs | Protected (`protectSecrets=true` / `--exclude-protected`); runtime/deploy SAs stay while secret IAM bindings remain |
| Secret Manager + IAM API enablement | Protected with the secrets |
| GCS Pulumi state bucket | Not a Pulumi-managed app resource; delete manually if retiring the project |
| GitHub secrets / variables | Untouched; WIF vars become **stale** until you bootstrap again |
| Operator SM passphrase copy | `vibegopher-pulumi-config-passphrase` (from `store-pulumi-passphrase.sh`) |

Removed by a successful destroy: Cloud Run, Artifact Registry repo, WIF pool/provider, project-level deploy IAM roles (non-protected).

After destroy, CI **cannot** redeploy until you bootstrap locally again and refresh `GCP_WORKLOAD_IDENTITY_PROVIDER` / `GCP_DEPLOY_SERVICE_ACCOUNT` (see job summary and [`infra/README.md`](../infra/README.md)).

**If CI destroy fails mid-flight** (e.g. 403 removing project IAM after WIF is already gone): finish on a laptop with project-owner ADC:

```bash
GCP_PROJECT_ID=YOUR_GCP_PROJECT_ID ./infra/scripts/gcloud-login.sh
export PULUMI_CONFIG_PASSPHRASE='…'   # or read from SM recovery secret
cd infra && pulumi login "$PULUMI_BACKEND_URL"
pulumi stack select "$PULUMI_STACK"
pulumi destroy --yes --exclude-protected
```

The deploy SA is granted `roles/resourcemanager.projectIamAdmin` so a **future** CI destroy can remove its own project IAM bindings (requires a bootstrap/`pulumi up` that applied that role before destroy).

---

## 9. Making this a GitHub template

Optional polish when publishing:

1. GitHub → **Settings → General → Template repository** (so forks start clean).
2. Clear or document example GitHub vars/secrets (never commit real values).
3. Leave Cloud Run stopped/destroyed; keep Neon only if you still need a personal demo.
4. Keep product auth/Slack work in a **private** repo that vendors or forks this template.
