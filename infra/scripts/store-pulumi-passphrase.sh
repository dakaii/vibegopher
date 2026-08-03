#!/usr/bin/env bash
# Store PULUMI_CONFIG_PASSPHRASE in GitHub Secrets + GCP Secret Manager (recovery).
# Does not print the passphrase.
set -euo pipefail

PROJECT_ID="${GCP_PROJECT_ID:?set GCP_PROJECT_ID}"
PASSPHRASE="${PULUMI_CONFIG_PASSPHRASE:?set PULUMI_CONFIG_PASSPHRASE}"
SM_SECRET_ID="${PULUMI_PASSPHRASE_SECRET_ID:-vibegopher-pulumi-config-passphrase}"

if ! command -v gh >/dev/null 2>&1; then
  echo "error: gh not found" >&2
  exit 1
fi
if ! command -v gcloud >/dev/null 2>&1; then
  echo "error: gcloud not found" >&2
  exit 1
fi

echo "Writing GitHub secret PULUMI_CONFIG_PASSPHRASE…"
printf '%s' "${PASSPHRASE}" | gh secret set PULUMI_CONFIG_PASSPHRASE

echo "Writing GCP Secret Manager secret ${SM_SECRET_ID} (operator recovery; not used by Cloud Run)…"
if gcloud secrets describe "${SM_SECRET_ID}" --project="${PROJECT_ID}" >/dev/null 2>&1; then
  printf '%s' "${PASSPHRASE}" | gcloud secrets versions add "${SM_SECRET_ID}" \
    --project="${PROJECT_ID}" \
    --data-file=- >/dev/null
else
  printf '%s' "${PASSPHRASE}" | gcloud secrets create "${SM_SECRET_ID}" \
    --project="${PROJECT_ID}" \
    --replication-policy=automatic \
    --data-file=- >/dev/null
fi

echo "ok: passphrase stored in GitHub Secrets + GCP Secret Manager"
echo "retrieve later:"
echo "  gcloud secrets versions access latest --secret=${SM_SECRET_ID} --project=${PROJECT_ID}"
