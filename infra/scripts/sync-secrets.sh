#!/usr/bin/env bash
# Sync secret payloads into GCP Secret Manager without putting them in Pulumi state.
# Reads values from environment variables matching the secret IDs.
set -euo pipefail

PROJECT_ID="${GCP_PROJECT_ID:?GCP_PROJECT_ID is required}"

# Required for a safe deploy. Optional secrets may be skipped when unset.
REQUIRED_SECRETS=(
  DATABASE_URL
  AUTH_SECRET
)
OPTIONAL_SECRETS=(
  GOOGLE_OAUTH_CLIENT_ID
  GOOGLE_OAUTH_CLIENT_SECRET
  GEMINI_API_KEY
)

PLACEHOLDER="REPLACE_ME"

sync_one() {
  local name="$1"
  local required="$2"
  local value="${!name-}"

  if [[ -z "${value}" ]]; then
    if [[ "${required}" == "true" ]]; then
      echo "error: required secret ${name} is not set in the environment" >&2
      exit 1
    fi
    echo "skip ${name} (env not set)"
    return 0
  fi

  if [[ "${value}" == "${PLACEHOLDER}" ]]; then
    echo "error: ${name} must not be the placeholder value ${PLACEHOLDER}" >&2
    exit 1
  fi

  if ! gcloud secrets describe "${name}" --project="${PROJECT_ID}" >/dev/null 2>&1; then
    echo "error: secret ${name} does not exist yet — run pulumi up first" >&2
    exit 1
  fi

  local current
  current="$(gcloud secrets versions access latest --secret="${name}" --project="${PROJECT_ID}" 2>/dev/null || true)"
  if [[ "${current}" == "${value}" ]]; then
    echo "ok   ${name} (unchanged)"
    return 0
  fi

  printf '%s' "${value}" | gcloud secrets versions add "${name}" \
    --project="${PROJECT_ID}" \
    --data-file=- >/dev/null
  echo "ok   ${name} (new version added)"
}

echo "Syncing secrets to project ${PROJECT_ID} (values are never printed)"

for name in "${REQUIRED_SECRETS[@]}"; do
  sync_one "${name}" true
done

for name in "${OPTIONAL_SECRETS[@]}"; do
  sync_one "${name}" false
done

# Fail closed if GCP still has placeholders for required secrets.
for name in "${REQUIRED_SECRETS[@]}"; do
  current="$(gcloud secrets versions access latest --secret="${name}" --project="${PROJECT_ID}")"
  if [[ "${current}" == "${PLACEHOLDER}" ]]; then
    echo "error: GCP secret ${name} is still ${PLACEHOLDER} after sync" >&2
    exit 1
  fi
done

echo "required secrets synced"
