# Secrets & config placement

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
| `PULUMI_ACCESS_TOKEN` | yes | Pulumi Cloud access |
| `DATABASE_URL` | yes | Neon **pooled** URL → synced to GCP SM for Cloud Run |
| `DATABASE_URL_MIGRATE` | strongly recommended | Neon **direct** (non-pooler) URL for goose; supports session advisory locks |
| `AUTH_SECRET` | yes | Synced into GCP SM; deploy fails if missing/`REPLACE_ME` |
| `GOOGLE_OAUTH_CLIENT_ID` | optional | Synced if present |
| `GOOGLE_OAUTH_CLIENT_SECRET` | optional | Synced if present |
| `GEMINI_API_KEY` | optional | Synced if present |

If `DATABASE_URL_MIGRATE` is unset, deploy falls back to `DATABASE_URL` but **rejects** URLs containing `-pooler`.

### Repository variables

| Name | Purpose |
|------|---------|
| `GCP_PROJECT_ID` | Target project |
| `GCP_REGION` | e.g. `us-central1` |
| `GCP_WORKLOAD_IDENTITY_PROVIDER` | Full WIF provider resource name (from stack output) |
| `GCP_DEPLOY_SERVICE_ACCOUNT` | Deploy SA email (from stack output) |
| `PULUMI_STACK` | e.g. `dev` or `prod` |

No long-lived GCP JSON keys — deploy uses **Workload Identity Federation**.

## Neon

| Item | Where |
|------|--------|
| Neon project / branch / roles | Neon console (or Neon API), not Pulumi GCP |
| App connection string | GitHub `DATABASE_URL` (pooled) → GCP SM `DATABASE_URL` |
| Migrate connection string | GitHub `DATABASE_URL_MIGRATE` (direct host) — CI only, not stored in GCP |
| DB schema | goose migrations in `db/migrations` (applied by deploy workflow, not Pulumi) |

## Destroy behavior

| Action | Result |
|--------|--------|
| **Destroy** workflow (default) | Removes Cloud Run, Artifact Registry wiring, deploy/runtime SAs, **and WIF**; **keeps Secret Manager** |
| **Destroy** with `destroy_secrets=true` | Also destroys Secret Manager secrets in this stack |

After a default destroy, GitHub Actions **cannot** redeploy until you run a local bootstrap `pulumi up` and refresh `GCP_WORKLOAD_IDENTITY_PROVIDER` / `GCP_DEPLOY_SERVICE_ACCOUNT`.

## Syncing secret values

Deploy order: **build/push → pulumi up (secret containers) → sync payloads → migrate → bump Cloud Run revision**.

```bash
export GCP_PROJECT_ID=...
export DATABASE_URL='postgres://...@...-pooler.../db?sslmode=require'
export AUTH_SECRET='...'
./infra/scripts/sync-secrets.sh
```

The script never prints secret values; required secrets must be set and must not be `REPLACE_ME`.
