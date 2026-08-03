# VibeGopher GCP infrastructure (Pulumi)

Provisions:

- Cloud Run API service
- Artifact Registry (Docker)
- Secret Manager secret **containers** + IAM (payloads synced separately)
- Runtime + deploy service accounts
- GitHub Actions Workload Identity Federation (when `githubOwner` is set)

## Prerequisites

- [Pulumi CLI](https://www.pulumi.com/docs/install/)
- Go 1.22+
- `gcloud` authenticated to the target project (`./infra/scripts/gcloud-login.sh`)
- A Pulumi backend (`pulumi login`)

### Pulumi backend: GCS (default for this repo)

State lives in a **private GCS bucket**. GitHub Actions uses Workload Identity (no Pulumi Cloud token). Local and CI must use the same `PULUMI_BACKEND_URL` (`gs://…`). The deploy SA needs `roles/storage.objectAdmin` on that bucket.

## First-time setup

Bootstrap **locally** once (WIF + Artifact Registry must exist before GitHub Actions can deploy):

```bash
# Browser login (does not print account emails)
GCP_PROJECT_ID=YOUR_GCP_PROJECT_ID ./infra/scripts/gcloud-login.sh

export PULUMI_CONFIG_PASSPHRASE='…from password manager…'

# Bucket + pulumi up (Cloud Run off) + GitHub WIF vars + bucket IAM
# Also writes PULUMI_CONFIG_PASSPHRASE to GitHub Secrets when the env var is set.
GCP_PROJECT_ID=YOUR_GCP_PROJECT_ID \
PULUMI_BACKEND_URL=gs://YOUR_STATE_BUCKET \
GITHUB_OWNER=YOUR_GH_USER_OR_ORG \
./infra/scripts/bootstrap-local.sh
```

Or follow the manual steps in [`docs/DEPLOY.md`](../docs/DEPLOY.md).

Then run **Deploy to GCP** (or push to `main`). Deploy logs into the GCS backend, sets `enableCloudRun=true`, syncs secrets, migrates, and updates Cloud Run.

## Deploy / destroy via GitHub Actions

- **Deploy** — `.github/workflows/deploy.yml` on push to `main` (build/push → `pulumi up` → sync secrets → goose migrate → Cloud Run revision bump)
- **Destroy** — `.github/workflows/destroy.yml` (`workflow_dispatch`); by default keeps Secret Manager but **removes WIF/deploy SA** — bootstrap locally again before the next CI deploy

## Protect Secret Manager (default)

`vibegopher:protectSecrets` defaults to `true`:

- Secret Manager API enablement is protected / retain-on-delete
- Each secret (+ placeholder versions / IAM on secrets) is protected / retain-on-delete
- Destroy workflow uses `pulumi destroy --yes --exclude-protected`

To also delete secrets, run the Destroy workflow with `destroy_secrets=true` (requires typing `destroy`).

## Layout

| File | Role |
|------|------|
| `main.go` | Stack entrypoint / exports |
| `config.go` | Pulumi config + protect helpers |
| `apis.go` | Enable GCP APIs |
| `secrets.go` | Secret Manager containers |
| `artifact_registry.go` | Docker repo |
| `iam.go` | Runtime / deploy service accounts |
| `cloudrun.go` | Cloud Run service |
| `github_wif.go` | GitHub OIDC → deploy SA |
| `scripts/sync-secrets.sh` | Push secret payloads to GCP |
| `SECRETS.md` | What goes where |
