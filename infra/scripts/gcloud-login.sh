#!/usr/bin/env bash
# Interactive gcloud (+ ADC) login for local Pulumi / sync-secrets.
# Does not print account emails, tokens, or other identity details.
set -euo pipefail

PROJECT_ID="${GCP_PROJECT_ID:-}"

usage() {
  cat <<'EOF'
Usage:
  GCP_PROJECT_ID=your-gcp-project-id ./infra/scripts/gcloud-login.sh

What it does:
  1. Opens browser flows for user login + Application Default Credentials
  2. Sets the active gcloud project (from GCP_PROJECT_ID only)
  3. Verifies credentials work without printing account identity

Notes:
  - Nothing is written to the repo; credentials stay in your local gcloud/ADC store.
  - Re-run when you see invalid_grant / RAPT reauth errors from Pulumi or gcloud.
EOF
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

if [[ -z "${PROJECT_ID}" ]]; then
  echo "error: set GCP_PROJECT_ID in the environment (do not put accounts/emails in this script)" >&2
  usage >&2
  exit 1
fi

if ! command -v gcloud >/dev/null 2>&1; then
  echo "error: gcloud CLI not found" >&2
  exit 1
fi

echo "Step 1/3: gcloud user login (browser)…"
# Omit --update-adc: older CLIs reject --update-adc=false; ADC is step 2.
# Redirect stdout to avoid account-list noise; errors still show on stderr.
gcloud auth login >/dev/null

echo "Step 2/3: Application Default Credentials (browser)…"
# Pulumi GCP provider and many local tools use ADC, not only gcloud user creds.
gcloud auth application-default login >/dev/null

echo "Step 3/3: set project + verify (no identity printed)…"
gcloud config set project "${PROJECT_ID}" --quiet

# Verify without printing account email or token value.
if ! gcloud auth print-access-token >/dev/null 2>&1; then
  echo "error: user credentials did not produce an access token" >&2
  exit 1
fi
if ! gcloud auth application-default print-access-token >/dev/null 2>&1; then
  echo "error: application-default credentials did not produce an access token" >&2
  exit 1
fi
if ! gcloud projects describe "${PROJECT_ID}" --format='value(projectId)' >/dev/null 2>&1; then
  echo "error: cannot describe project ${PROJECT_ID} (check ID / IAM)" >&2
  exit 1
fi

echo "ok: gcloud + ADC ready for project ${PROJECT_ID}"
echo "next:"
echo "  GCP_PROJECT_ID=${PROJECT_ID} \\"
echo "  PULUMI_BACKEND_URL=gs://YOUR_STATE_BUCKET \\"
echo "  GITHUB_OWNER=YOUR_GH_USER_OR_ORG \\"
echo "  ./infra/scripts/bootstrap-local.sh"
