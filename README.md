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

## Prerequisites

- **Docker** and **Docker Compose ≥ 2.24** (needed for optional `env_file` / `depends_on` with `required: false`)
- **Make**

## Quick start

```bash
make build
make up                    # Postgres + migrate + API → http://localhost:8081
```

Frontend:

```bash
cd frontend
cp .env.example .env       # VITE_CLERK_PUBLISHABLE_KEY; leave VITE_API_BASE_URL empty locally
npm install && npm run dev # http://localhost:5173
```

Tests: `make test`  
Frontend check: `cd frontend && npm run check`

## Local database

- `config/development.conf` — default Docker Postgres
- `config/test.conf` — test env
- `config/local.conf.example` — template for machine-local overrides (e.g. Neon)
- `config/local.conf` — optional gitignored override

By default, Compose loads `config/development.conf` for **backend** and **migrator**. If `config/local.conf` exists, it is loaded after and overrides matching variables. The local Postgres service always uses `development.conf` only.

Local Docker Postgres is behind the `local-db` profile. `make run-db` / `make up` / `make migrate` enable that profile automatically when `config/local.conf` is absent. Backend and migrator declare an optional dependency on `postgresql-dev`, so Neon users are not forced to start local Postgres.

**Local Docker Postgres (default):**

```bash
make run-db
make migrate
docker compose --profile local-db up backend
# or: make up
```

**Use Neon (or another remote DB) locally:**

```bash
cp config/local.conf.example config/local.conf
# Edit config/local.conf with your credentials (POSTGRES_SSLMODE=require for Neon)
make migrate                 # uses local.conf; does not start postgresql-dev
docker compose up backend    # no local-db profile → no local Postgres
```

## Configuration

| Variable | Purpose |
|----------|---------|
| `CLERK_SECRET_KEY` | Clerk Backend API secret (required in production; verifies SPA session JWTs) |
| `AUTH_SECRET` | HS256 signing for password auth (tests/local; not required for Clerk prod) |
| `DATABASE_URL` / `POSTGRES_*` | Database (Neon pooled URL in cloud) |
| `POSTGRES_SSLMODE` | DSN sslmode when not using `DATABASE_URL` (default `disable`) |
| `GEMINI_API_KEY` | Critic LLM key |
| `CRITIC_WORKER_ENABLED` | Start in-process worker (`BOT_WORKER_ENABLED` still works) |
| `CORS_ORIGIN` | Comma-separated SPA origins (required explicit in production; also used as Clerk `azp` allow-list) |
| `ENABLE_PASSWORD_AUTH` | Expose `/api/signup` + `/api/login` (tests/local only; **forbidden in production**) |
| `APP_ENV` | `development` \| `test` \| `production` |

## Docs

- [API](./docs/API.md)
- [Deploy](./docs/DEPLOY.md)
- [Bot persona](./docs/BOT_PERSONA.md)
- [Roadmap / missing features](./docs/ROADMAP.md)
- [Infra & secrets](./infra/README.md) · [SECRETS.md](./infra/SECRETS.md)
- [Contributing](./CONTRIBUTING.md) · [Security](./SECURITY.md)

## Schema (high level)

- **users** — username, optional password, `clerk_user_id`, legacy `google_sub`, `email`, `is_bot`
- **posts** / **comments** — content + author
- **bot_jobs** — async critic queue
- Seeded bot user: `@vibe_critic` (`is_bot=true`)

## GCP deploy

See [`docs/DEPLOY.md`](./docs/DEPLOY.md). Push to `main` runs build → `pulumi up` → secret sync → goose migrate → Cloud Run bump. Destroy keeps Secret Manager protected by default.

Secret placement (GCP Secret Manager vs Pulumi vs GitHub): [`infra/SECRETS.md`](./infra/SECRETS.md).

## License

[MIT](./LICENSE)
