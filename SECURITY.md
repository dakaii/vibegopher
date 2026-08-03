# Security Policy

## Supported versions

This project is a small public demo. Fixes land on `main`; there are no long-term support branches.

## Reporting a vulnerability

Please **do not** open a public GitHub issue for security problems.

Email the maintainer via the contact listed on the [GitHub profile](https://github.com/dakaii), or open a private [GitHub security advisory](https://github.com/dakaii/vibegopher/security/advisories/new) if available.

Include:

- Affected component (API, critic worker, frontend, infra)
- Steps to reproduce
- Impact assessment
- Any suggested fix

## What not to commit

Never commit:

- Real Neon / Postgres URLs
- `AUTH_SECRET`, Gemini keys, Google OAuth client secrets
- Pulumi stack files with encrypted secrets (`infra/Pulumi.*.yaml` except examples)
- Local `.env` files (`frontend/.env`, root `.env`)
- Production credentials of any kind

Local Compose defaults in `config/*.conf` are intentionally fake and for development/tests only.

## Production expectations

- Set a strong `AUTH_SECRET` (≥16 chars; not the sample values)
- Set an explicit `CORS_ORIGIN` (never `*` in production)
- Keep `ENABLE_PASSWORD_AUTH` unset/false in production
- Prefer Google Sign-In for the SPA
- Treat critic link-fetching as untrusted input (SSRF protections exist; still rate-limit externally)
