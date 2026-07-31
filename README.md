# VibeGopher

Public Twitter-style demo: Go REST API, Vue 3 SPA, and `@vibe_critic` (Gemini) that replies to posts and comments.

This repository is intentionally **open**. A future Slack SaaS built around the same critic idea should live in a **separate private repo** — do not add Slack product code here.

## Features

- Google Sign-In (primary) + app JWT sessions
- Posts & comments (280 chars)
- Async critic worker on `bot_jobs` (in-process or `cmd/bot`)
- PostgreSQL via goose migrations (Neon in cloud)
- Pulumi → GCP Cloud Run deploy + destroy workflows

## Stack

| Layer | Choice |
|-------|--------|
| API | Go, Gorilla Mux, GORM (queries only) |
| Migrations | [goose](https://github.com/pressly/goose) SQL in `db/migrations` |
| Frontend | Vue 3 + TypeScript + Biome (`frontend/`) |
| AI | Gemini via `internal/critic` |
| Infra | Pulumi (Go) → GCP |

## Quick start

```bash
make build
make up                    # Postgres + migrate + API → http://localhost:8081

cd frontend
cp .env.example .env       # VITE_GOOGLE_CLIENT_ID, VITE_API_BASE_URL
npm install && npm run dev # http://localhost:5173
```

Tests: `make test`  
Frontend check: `cd frontend && npm run check`

## Configuration

| Variable | Purpose |
|----------|---------|
| `AUTH_SECRET` | JWT signing key (strong value required in production) |
| `DATABASE_URL` / `POSTGRES_*` | Database (Neon pooled URL in cloud) |
| `GOOGLE_OAUTH_CLIENT_ID` | Google GIS audience |
| `GEMINI_API_KEY` | Critic LLM key |
| `CRITIC_WORKER_ENABLED` | Start in-process worker (`BOT_WORKER_ENABLED` still works) |
| `CORS_ORIGIN` | Comma-separated SPA origins (required explicit in production) |
| `ENABLE_PASSWORD_AUTH` | Expose `/api/signup` + `/api/login` (tests/local only; **forbidden in production**) |
| `APP_ENV` | `development` \| `test` \| `production` |

Local Compose env files: `config/development.conf`, `config/test.conf` (fake secrets only).

## Docs

- [API](./docs/API.md)
- [Deploy](./docs/DEPLOY.md)
- [Bot persona](./docs/BOT_PERSONA.md)
- [Roadmap / missing features](./docs/ROADMAP.md)
- [Infra & secrets](./infra/README.md) · [SECRETS.md](./infra/SECRETS.md)
- [Contributing](./CONTRIBUTING.md) · [Security](./SECURITY.md)

## Schema (high level)

- **users** — username, optional password, `google_sub`, `email`, `is_bot`
- **posts** / **comments** — content + author
- **bot_jobs** — async critic queue
- Seeded bot user: `@vibe_critic` (`is_bot=true`)

## GCP deploy

See [`docs/DEPLOY.md`](./docs/DEPLOY.md). Push to `main` runs build → `pulumi up` → secret sync → goose migrate → Cloud Run bump. Destroy keeps Secret Manager protected by default.

## License

[MIT](./LICENSE)
