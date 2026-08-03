# VibeGopher GCP infrastructure (Pulumi)

Part of the **public template**: fork this repo, bring your own Neon + GCP project, bootstrap once locally, then deploy/destroy from GitHub Actions. Production Slack/auth products should live in a **private** derivative — not in this open template.

Provisions:

- Cloud Run API service
- Artifact Registry (Docker)
- Secret Manager secret **containers** + IAM (payloads synced separately)
- Runtime + deploy service accounts
- GitHub Actions Workload Identity Federation (when `githubOwner` is set)

Does **not** provision: Neon, the GCS state bucket lifecycle (bootstrap creates/uses it), or GitHub secrets/vars.

## Prerequisites

- [Pulumi CLI](https://www.pulumi.com/docs/install/)
- Go 1.24+ (infra program)
- `gcloud` authenticated to the target project (`./infra/scripts/gcloud-login.sh`)
- A Pulumi backend (`pulumi login` to `gs://…`)

### Pulumi backend: GCS (default for this repo)

State lives in a **private GCS bucket**. GitHub Actions uses Workload Identity (no Pulumi Cloud token). Local and CI must use the same `PULUMI_BACKEND_URL` (`gs://…`). The deploy SA needs `roles/storage.objectAdmin` on that bucket.

## First-time setup

Bootstrap **locally** once (WIF + Artifact Registry must exist before GitHub Actions can deploy):

```bash
# Browser login (does not print account emails)
GCP_PROJECT_ID=YOUR_GCP_PROJECT_ID ./infra/scripts/gcloud-login.sh

export PULUMI_CONFIG_PASSPHRASE='…strong random value…'
./infra/scripts/store-pulumi-passphrase.sh   # GitHub Secrets + GCP SM recovery copy

# Bucket + pulumi up (Cloud Run off) + GitHub WIF vars + bucket IAM
GCP_PROJECT_ID=YOUR_GCP_PROJECT_ID \
PULUMI_BACKEND_URL=gs://YOUR_STATE_BUCKET \
GITHUB_OWNER=YOUR_GH_USER_OR_ORG \
./infra/scripts/bootstrap-local.sh
```

Passphrase storage / rotation: [SECRETS.md](./SECRETS.md).

Full checklist (what you bring, Google auth, destroy leftovers): [`docs/DEPLOY.md`](../docs/DEPLOY.md).

Then run **Deploy to GCP** (or push to `main`). Deploy logs into the GCS backend, sets `enableCloudRun=true`, syncs secrets, migrates, and updates Cloud Run.

## Deploy / destroy via GitHub Actions

| Workflow | File | How to run |
|----------|------|------------|
| **Deploy** | `.github/workflows/deploy.yml` | Push to `main` or `gh workflow run "Deploy to GCP"` |
| **Destroy** | `.github/workflows/destroy.yml` | `gh workflow run "Destroy GCP infrastructure" -f confirm=destroy` |

Destroy defaults:

- `confirm` must be exactly `destroy`
- `destroy_secrets=false` — keeps protected Secret Manager (and SAs still referenced by secret IAM); **removes** Cloud Run, Artifact Registry, WIF, and non-protected project IAM
- Does **not** delete Neon, the GCS state bucket, or GitHub secrets/vars
- Deploy SA includes `roles/resourcemanager.projectIamAdmin` so CI can delete project `IAMMember` bindings; if destroy fails after WIF is gone, finish with local owner ADC (see [`docs/DEPLOY.md`](../docs/DEPLOY.md))

After destroy, bootstrap locally again before the next CI deploy (WIF is gone). Optional full SM wipe: `-f destroy_secrets=true`.

```bash
gh workflow run "Destroy GCP infrastructure" -f confirm=destroy
gh run watch
```

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
| `scripts/gcloud-login.sh` | Quiet browser ADC login |
| `scripts/bootstrap-local.sh` | State bucket + first `pulumi up` + WIF GitHub vars |
| `scripts/store-pulumi-passphrase.sh` | GitHub secret + SM recovery copy |
| `scripts/sync-secrets.sh` | Push secret payloads to GCP |
| `SECRETS.md` | What goes where |
