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
- `gcloud` authenticated to the target project
- A Pulumi backend (`pulumi login`)

## First-time setup

Bootstrap **locally** once (Workload Identity and Artifact Registry must exist before GitHub Actions can deploy):

```bash
cd infra
cp Pulumi.dev.yaml.example Pulumi.dev.yaml
# edit gcp:project and vibegopher:githubOwner

pulumi stack init dev   # or select an existing stack
pulumi config set gcp:project YOUR_PROJECT_ID
pulumi config set vibegopher:githubOwner YOUR_GH_OWNER
# No image in Artifact Registry yet:
pulumi config set vibegopher:enableCloudRun false
pulumi up
```

Copy stack outputs into GitHub:

- `workloadIdentityProvider` → repo variable `GCP_WORKLOAD_IDENTITY_PROVIDER`
- `deployServiceAccount` → repo variable `GCP_DEPLOY_SERVICE_ACCOUNT`
- set `GCP_PROJECT_ID`, `GCP_REGION`, `PULUMI_STACK`
- set secrets listed in [SECRETS.md](./SECRETS.md) including `PULUMI_ACCESS_TOKEN`

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

- **Deploy** — `.github/workflows/deploy.yml` on push to `main` (build → Artifact Registry → sync secrets → `pulumi up`)
- **Destroy** — `.github/workflows/destroy.yml` (`workflow_dispatch`); by default keeps Secret Manager

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
