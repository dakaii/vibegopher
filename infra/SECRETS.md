# Secrets & config placement

Public **template** repo: you bring Neon + GCP + GitHub secrets/vars. Do not commit secret values. A private Slack/auth product should keep product secrets in that private repo — not here.

Auth on `main` uses **Google GIS** + app JWT (`AUTH_SECRET` / `GOOGLE_OAUTH_CLIENT_ID`). Clerk is not required; see optional unmerged [PR #14](https://github.com/dakaii/vibegopher/pull/14) only as a reference.

Rule of thumb:

- **GCP Secret Manager** — runtime secrets the Cloud Run service (and later the AI bot) need at request time
- **Pulumi config** — non-secret infrastructure settings (and encrypted bootstrap values only if you must)
- **GitHub Actions secrets / vars** — CI credentials and the source of truth used to *sync* secret payloads into GCP
- **Neon** — database lives in Neon; pooled URL for the app, **direct** URL for migrations

Pulumi creates Secret Manager **secret containers + IAM**. It does **not** store real secret payloads in Pulumi state when you use the sync script / deploy workflow. Placeholder `REPLACE_ME` versions may exist until the first sync; deploy **fails** if required secrets are still placeholders after sync.

## GCP Secret Manager (runtime)

| Secret ID | Purpose | Who reads it |
|-----------|---------|--------------|
| `DATABASE_URL` | Neon **pooled** Postgres URL for the API (`…-pooler…` is OK here) | Cloud Run API |
| `AUTH_SECRET` | App session JWT signing key (required) | Cloud Run API |
| `GOOGLE_OAUTH_CLIENT_ID` | Google OAuth / GIS client ID (token `aud` check) | Cloud Run API |
| `GOOGLE_OAUTH_CLIENT_SECRET` | Optional; only if you use authorization-code exchange | Cloud Run API |
| `GEMINI_API_KEY` | AI bot LLM key (skip later if you move to Vertex ADC) | Cloud Run / bot worker |

Operator-only (created by `store-pulumi-passphrase.sh`, **not** wired into Cloud Run):

| Secret ID | Purpose |
|-----------|---------|
| `vibegopher-pulumi-config-passphrase` | Recovery copy of `PULUMI_CONFIG_PASSPHRASE` |

These secrets (and the Secret Manager API enablement) are **protected by default** (`vibegopher:protectSecrets=true`).  
`pulumi destroy --exclude-protected` tears down Cloud Run / AR / deploy wiring / WIF but **keeps Secret Manager**.

## Pulumi config (infra, mostly non-secret)

Set with `pulumi config set` (from `infra/`):

| Key | Secret? | Purpose |
|-----|---------|---------|
| `gcp:project` | no | GCP project ID |
| `gcp:region` / `vibegopher:region` | no | Deploy region (default `us-central1`) |
| `vibegopher:serviceName` | no | Cloud Run service name |
| `vibegopher:artifactRepoId` | no | Artifact Registry repo id |
| `vibegopher:imageName` / `imageTag` | no | Image coordinates |
| `vibegopher:githubOwner` / `githubRepo` | no | WIF binding for GitHub Actions |
| `vibegopher:protectSecrets` | no | Default `true` — protect Secret Manager on destroy |
| `vibegopher:createPlaceholderSecretVersions` | no | Default `true` — bootstrap `REPLACE_ME` versions |

**Do not** put Neon URLs, API keys, or OAuth client secrets in Pulumi config if you can avoid it.

## GitHub Actions (CI)

### Repository secrets

| Name | Required | Purpose |
|------|----------|---------|
| `DATABASE_URL` | yes | Neon **pooled** URL → synced to GCP SM for Cloud Run |
| `DATABASE_URL_MIGRATE` | strongly recommended | Neon **direct** (non-pooler) URL for goose; supports session advisory locks |
| `AUTH_SECRET` | yes | Synced into GCP SM; deploy fails if missing/`REPLACE_ME` |
| `PULUMI_CONFIG_PASSPHRASE` | yes | Unlocks encrypted Pulumi stack config on the **GCS** backend (local + CI must match) |
| `GOOGLE_OAUTH_CLIENT_ID` | optional | Synced if present |
| `GOOGLE_OAUTH_CLIENT_SECRET` | optional | Synced if present |
| `GEMINI_API_KEY` | optional | Synced if present |

If `DATABASE_URL_MIGRATE` is unset, deploy falls back to `DATABASE_URL` but **rejects** URLs containing `-pooler`.

**No `PULUMI_ACCESS_TOKEN`** — this repo uses a **GCS Pulumi backend** (`PULUMI_BACKEND_URL`). CI authenticates with GitHub → GCP Workload Identity.

### `PULUMI_CONFIG_PASSPHRASE` (where to keep it)

This unlocks encrypted Pulumi stack config on the GCS backend. It is **not** a Pulumi Cloud token. Losing it can lock you out of encrypted stack config.

You do **not** need a dedicated password manager. Use both of these:

| Place | Role |
|-------|------|
| **GitHub secret** `PULUMI_CONFIG_PASSPHRASE` | Required for CI deploy/destroy |
| **GCP Secret Manager** `vibegopher-pulumi-config-passphrase` | Operator recovery copy (humans / laptop). **Not** mounted into Cloud Run |

```bash
# After choosing a strong passphrase (rotate if it was ever pasted into chat/logs):
export GCP_PROJECT_ID=YOUR_GCP_PROJECT_ID
export PULUMI_CONFIG_PASSPHRASE='…strong random value…'
./infra/scripts/store-pulumi-passphrase.sh
```

Retrieve later (prints the secret — do this in a private terminal):

```bash
gcloud secrets versions access latest \
  --secret=vibegopher-pulumi-config-passphrase \
  --project=YOUR_GCP_PROJECT_ID
```

**Do not** commit the passphrase, put it in `config/local.conf`, or store it only in chat/Notes.

#### Rotate the passphrase

Safe to do anytime after the stack exists:

```bash
cd infra
export PULUMI_CONFIG_PASSPHRASE='…current…'
pulumi login "$PULUMI_BACKEND_URL"   # e.g. gs://YOUR_STATE_BUCKET
pulumi stack select dev
pulumi stack change-secrets-provider passphrase
# enter the NEW passphrase when prompted

export PULUMI_CONFIG_PASSPHRASE='…new…'
./scripts/store-pulumi-passphrase.sh   # updates GitHub + GCP SM
```

### Repository variables

| Name | Purpose |
|------|---------|
| `GCP_PROJECT_ID` | Target project |
| `GCP_REGION` | e.g. `us-central1` |
| `PULUMI_BACKEND_URL` | `gs://YOUR_STATE_BUCKET` (private bucket; CI + local must use the same) |
| `GCP_WORKLOAD_IDENTITY_PROVIDER` | Full WIF provider resource name (from stack output) |
| `GCP_DEPLOY_SERVICE_ACCOUNT` | Deploy SA email (from stack output) |
| `PULUMI_STACK` | e.g. `dev` or `prod` |
| `CORS_ORIGIN` | Frontend origin(s) for the API (`https://app.example.com` or comma-separated). Synced into Pulumi as `vibegopher:corsOrigin` / Cloud Run `CORS_ORIGIN`. |

No long-lived GCP JSON keys — deploy uses **Workload Identity Federation**. The deploy SA also needs `roles/storage.objectAdmin` on the state bucket (bootstrap script grants this).

## Neon

| Item | Where |
|------|--------|
| Neon project / branch / roles | Neon console (or Neon API), not Pulumi GCP |
| App connection string | GitHub `DATABASE_URL` (pooled) → GCP SM `DATABASE_URL` |
| Migrate connection string | GitHub `DATABASE_URL_MIGRATE` (direct host) — CI only, not stored in GCP |
| DB schema | goose migrations in `db/migrations` (applied by deploy workflow, not Pulumi) |

## Destroy behavior

```bash
gh workflow run "Destroy GCP infrastructure" -f confirm=destroy
# optional: -f destroy_secrets=true
```

| Action | Result |
|--------|--------|
| **Destroy** workflow (default) | Removes Cloud Run, Artifact Registry wiring, deploy/runtime SAs, **and WIF**; **keeps Secret Manager** |
| **Destroy** with `destroy_secrets=true` | Also destroys Secret Manager secrets in this stack |
| **Never destroyed by the workflow** | Neon, GCS Pulumi state bucket, GitHub secrets/vars |

After a default destroy, GitHub Actions **cannot** redeploy until you run a local bootstrap `pulumi up` and refresh `GCP_WORKLOAD_IDENTITY_PROVIDER` / `GCP_DEPLOY_SERVICE_ACCOUNT`.

See [`docs/DEPLOY.md`](../docs/DEPLOY.md) § Destroy for the leftover checklist.

## Syncing secret values

Deploy order: **build/push → pulumi up (secret containers) → sync payloads → migrate → bump Cloud Run revision**.

```bash
export GCP_PROJECT_ID=...
export DATABASE_URL='postgres://...@...-pooler.../db?sslmode=require'
export AUTH_SECRET='...'
./infra/scripts/sync-secrets.sh
```

The script never prints secret values; required secrets must be set and must not be `REPLACE_ME`.
