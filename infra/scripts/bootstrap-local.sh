#!/usr/bin/env bash
# One-time local Pulumi bootstrap (GCS backend + WIF).
# Does not print account emails, tokens, or secret values.
set -euo pipefail

PROJECT_ID="${GCP_PROJECT_ID:?set GCP_PROJECT_ID}"
BACKEND_URL="${PULUMI_BACKEND_URL:?set PULUMI_BACKEND_URL (gs://…)}"
STACK="${PULUMI_STACK:-dev}"
GITHUB_OWNER="${GITHUB_OWNER:?set GITHUB_OWNER (GitHub user/org for WIF)}"
GITHUB_REPO="${GITHUB_REPO:-vibegopher}"
REGION="${GCP_REGION:-us-central1}"
CORS_ORIGIN="${CORS_ORIGIN:-http://localhost:5173}"

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
INFRA="${ROOT}/infra"

case "${BACKEND_URL}" in
  gs://*) ;;
  *)
    echo "error: PULUMI_BACKEND_URL must start with gs://" >&2
    exit 1
    ;;
esac

BUCKET="${BACKEND_URL#gs://}"
BUCKET="${BUCKET%%/*}"

if ! command -v gcloud >/dev/null 2>&1; then
  echo "error: gcloud not found" >&2
  exit 1
fi
if ! command -v pulumi >/dev/null 2>&1; then
  echo "error: pulumi not found" >&2
  exit 1
fi
if ! command -v go >/dev/null 2>&1; then
  echo "error: go not found on PATH (required to compile infra/)" >&2
  echo "  install: brew install go" >&2
  exit 1
fi
if ! command -v gh >/dev/null 2>&1; then
  echo "error: gh not found (needed to set GitHub variables)" >&2
  exit 1
fi

# Verify ADC/user creds without printing identity.
if ! gcloud auth print-access-token >/dev/null 2>&1; then
  echo "error: gcloud not logged in — run:" >&2
  echo "  GCP_PROJECT_ID=${PROJECT_ID} ${ROOT}/infra/scripts/gcloud-login.sh" >&2
  exit 1
fi
if ! gcloud auth application-default print-access-token >/dev/null 2>&1; then
  echo "error: ADC missing — run gcloud-login.sh" >&2
  exit 1
fi

gcloud config set project "${PROJECT_ID}" --quiet

echo "Ensuring private state bucket exists…"
if ! gcloud storage buckets describe "gs://${BUCKET}" --project="${PROJECT_ID}" >/dev/null 2>&1; then
  gcloud storage buckets create "gs://${BUCKET}" \
    --project="${PROJECT_ID}" \
    --location="${REGION}" \
    --uniform-bucket-level-access \
    --quiet
  echo "created gs://${BUCKET}"
else
  echo "bucket already exists"
fi

# Best-effort: block public ACLs (ignore if org policy already enforces this).
gcloud storage buckets update "gs://${BUCKET}" --public-access-prevention --quiet 2>/dev/null || true

cd "${INFRA}"
if [[ ! -f Pulumi.dev.yaml ]]; then
  cp Pulumi.dev.yaml.example Pulumi.dev.yaml
fi

echo "Logging Pulumi into GCS backend…"
pulumi login "${BACKEND_URL}"

if ! pulumi stack select "${STACK}" >/dev/null 2>&1; then
  pulumi stack init "${STACK}"
fi

pulumi config set gcp:project "${PROJECT_ID}"
pulumi config set vibegopher:region "${REGION}"
pulumi config set vibegopher:githubOwner "${GITHUB_OWNER}"
pulumi config set vibegopher:githubRepo "${GITHUB_REPO}"
pulumi config set vibegopher:corsOrigin "${CORS_ORIGIN}"
pulumi config set vibegopher:enableCloudRun false
pulumi config set vibegopher:protectSecrets true

echo "Running pulumi up (enableCloudRun=false)…"
pulumi up --yes

WIF="$(pulumi stack output workloadIdentityProvider)"
DEPLOY_SA="$(pulumi stack output deployServiceAccount)"

if [[ -z "${WIF}" || -z "${DEPLOY_SA}" ]]; then
  echo "error: missing stack outputs (githubOwner/repo may be unset so WIF was not created)" >&2
  exit 1
fi

echo "Granting deploy SA objectAdmin on state bucket…"
gcloud storage buckets add-iam-policy-binding "gs://${BUCKET}" \
  --member="serviceAccount:${DEPLOY_SA}" \
  --role="roles/storage.objectAdmin" \
  --quiet >/dev/null

echo "Writing GitHub Actions variables (no secrets printed)…"
gh variable set GCP_WORKLOAD_IDENTITY_PROVIDER --body "${WIF}"
gh variable set GCP_DEPLOY_SERVICE_ACCOUNT --body "${DEPLOY_SA}"
gh variable set PULUMI_BACKEND_URL --body "${BACKEND_URL}"
gh variable set GCP_PROJECT_ID --body "${PROJECT_ID}"
gh variable set GCP_REGION --body "${REGION}"
gh variable set PULUMI_STACK --body "${STACK}"
gh variable set CORS_ORIGIN --body "${CORS_ORIGIN}"

echo
echo "ok: bootstrap complete"
echo "next:"
echo "  1. Merge the GCS deploy workflow PR (if not already on main)"
echo "  2. Run GitHub Action: Deploy to GCP"
echo "  3. Optional later: GOOGLE_OAUTH_CLIENT_ID, GEMINI_API_KEY, prod CORS_ORIGIN"
