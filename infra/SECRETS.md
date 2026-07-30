# Secrets & config placement

Rule of thumb:

- **GCP Secret Manager** — runtime secrets the Cloud Run service (and later the AI bot) need at request time
- **Pulumi config** — non-secret infrastructure settings (and encrypted bootstrap values only if you must)
- **GitHub Actions secrets / vars** — CI credentials and the source of truth used to *sync* secret payloads into GCP
- **Neon** — database lives in Neon; only the connection string is stored in GCP

Pulumi creates Secret Manager **secret containers + IAM**. It does **not** store real secret payloads in Pulumi state when you use the sync script / deploy workflow. Placeholder `REPLACE_ME` versions may exist until the first sync.

## GCP Secret Manager (runtime)

| Secret ID | Purpose | Who reads it |
|-----------|---------|--------------|
| `DATABASE_URL` | Neon pooled Postgres URL (`postgres://...@.../neondb?sslmode=require`). App prefers this over `POSTGRES_*`. | Cloud Run API |
| `AUTH_SECRET` | App session JWT signing key (if you mint tokens after Google Auth) | Cloud Run API |
| `GOOGLE_OAUTH_CLIENT_ID` | Google OAuth / GIS client ID (token `aud` check) | Cloud Run API |
| `GOOGLE_OAUTH_CLIENT_SECRET` | Optional; only if you use authorization-code exchange | Cloud Run API |
| `GEMINI_API_KEY` | AI bot LLM key (skip later if you move to Vertex ADC) | Cloud Run / bot worker |

These secrets (and the Secret Manager API enablement) are **protected by default** (`vibegopher:protectSecrets=true`).  
`pulumi destroy --exclude-protected` tears down Cloud Run / AR / deploy wiring but **keeps Secret Manager**.

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

**Do not** put Neon URLs, API keys, or OAuth client secrets in Pulumi config if you can avoid it. If you ever must (`pulumi config set --secret`), treat that as bootstrap-only and rotate into Secret Manager afterward.

## GitHub Actions (CI)

### Repository secrets

| Name | Purpose |
|------|---------|
| `PULUMI_ACCESS_TOKEN` | Pulumi Cloud (or self-hosted backend) access |
| `DATABASE_URL` | Synced into GCP SM on deploy |
| `AUTH_SECRET` | Synced into GCP SM on deploy |
| `GOOGLE_OAUTH_CLIENT_ID` | Synced into GCP SM on deploy |
| `GOOGLE_OAUTH_CLIENT_SECRET` | Optional; synced if present |
| `GEMINI_API_KEY` | Optional; synced if present |

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
| Connection string | GitHub secret → synced to GCP `DATABASE_URL` |
| DB schema | Atlas migrations in this repo |

## Destroy behavior

| Action | Result |
|--------|--------|
| **Destroy** workflow (default) | `pulumi destroy --exclude-protected` — removes Cloud Run, Artifact Registry (images may remain per GCP policy), SAs, WIF; **keeps Secret Manager API + secrets** |
| **Destroy** with `destroy_secrets=true` | Unprotects secret resources, then full destroy — **irreversible for SM secrets in this stack** |

## Syncing secret values

```bash
# Local (requires gcloud auth + project set)
./infra/scripts/sync-secrets.sh

# Or let .github/workflows/deploy.yml do it after auth
```

The script adds a new Secret Manager version only when the value changed; it never prints secret values.
