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
- `gcloud` authenticated to the target project (`./scripts/gcloud-login.sh`)
- A Pulumi backend (`pulumi login`)

### Pulumi backend: Cloud token vs GCS

The deploy workflow currently expects a **Pulumi Cloud** token (`PULUMI_ACCESS_TOKEN`). That is the usual setup when state lives in Pulumi Cloud.

You can **avoid that token** by using a **self-managed backend** (e.g. `pulumi login gs://YOUR_STATE_BUCKET`) and relying on GitHub → GCP **Workload Identity** for both `gcloud` and Pulumi state access. That is a solid GCP-centric practice and one less SaaS secret — but CI must `pulumi login` to the same bucket, and the deploy SA needs object Admin (or equivalent) on that bucket. Do not mix backends for the same stack.

## First-time setup

Bootstrap **locally** once (Workload Identity and Artifact Registry must exist before GitHub Actions can deploy):

```bash
# From repo root — browser login; does not print account emails
GCP_PROJECT_ID=YOUR_GCP_PROJECT_ID ./infra/scripts/gcloud-login.sh

cd infra
cp Pulumi.dev.yaml.example Pulumi.dev.yaml
# edit gcp:project and vibegopher:githubOwner

pulumi login                 # Pulumi Cloud, or: pulumi login gs://YOUR_STATE_BUCKET
pulumi stack init dev        # or select an existing stack
pulumi config set gcp:project YOUR_PROJECT_ID
pulumi config set vibegopher:githubOwner YOUR_GH_OWNER
# No image in Artifact Registry yet:
pulumi config set vibegopher:enableCloudRun false
pulumi up
```

Copy stack outputs into GitHub:

- `workloadIdentityProvider` → repo variable `GCP_WORKLOAD_IDENTITY_PROVIDER`
- `deployServiceAccount` → repo variable `GCP_DEPLOY_SERVICE_ACCOUNT`
- set `GCP_PROJECT_ID`, `GCP_REGION`, `PULUMI_STACK`, `CORS_ORIGIN` (easy ones may already be set)
- set secrets in [SECRETS.md](./SECRETS.md); `PULUMI_ACCESS_TOKEN` only if using Pulumi Cloud

Sync runtime secret values:

```bash
export GCP_PROJECT_ID=...
export DATABASE_URL='postgres://...'   # Neon pooled URL
export AUTH_SECRET='...'
export GOOGLE_OAUTH_CLIENT_ID='...'
./scripts/sync-secrets.sh
```

Then run the **Deploy** GitHub Action (or push to `main`). It pushes the image, re-syncs secrets, sets `enableCloudRun=true`, and runs `pulumi up`.

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
