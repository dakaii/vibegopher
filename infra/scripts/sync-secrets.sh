#!/usr/bin/env bash
# Sync secret payloads into GCP Secret Manager without putting them in Pulumi state.
# Reads values from environment variables matching the secret IDs.
set -euo pipefail

PROJECT_ID="${GCP_PROJECT_ID:?GCP_PROJECT_ID is required}"
SECRETS=(
  DATABASE_URL
  AUTH_SECRET
  GOOGLE_OAUTH_CLIENT_ID
  GOOGLE_OAUTH_CLIENT_SECRET
  GEMINI_API_KEY
)

echo "Syncing secrets to project ${PROJECT_ID} (values are never printed)"

for name in "${SECRETS[@]}"; do
  value="${!name-}"
  if [[ -z "${value}" ]]; then
    echo "skip ${name} (env not set)"
    continue
  fi

  if ! gcloud secrets describe "${name}" --project="${PROJECT_ID}" >/dev/null 2>&1; then
    echo "error: secret ${name} does not exist yet — run pulumi up first" >&2
    exit 1
  fi

  # Avoid creating a new version if the latest payload already matches.
  current="$(gcloud secrets versions access latest --secret="${name}" --project="${PROJECT_ID}" 2>/dev/null || true)"
  if [[ "${current}" == "${value}" ]]; then
    echo "ok   ${name} (unchanged)"
    continue
  fi

  printf '%s' "${value}" | gcloud secrets versions add "${name}" \
    --project="${PROJECT_ID}" \
    --data-file=- >/dev/null
  echo "ok   ${name} (new version added)"
done
