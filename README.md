# VibeGopher

**Public template** for a Twitter-style demo: Go REST API, Vue 3 SPA, Neon Postgres, `@vibe_critic` (Gemini), and Pulumi → GCP (Cloud Run + GCS state + GitHub Actions WIF).

Fork or use as a reference stack. A production Slack SaaS / alternate auth product should live in a **separate private repo** that references this template — do not add Slack product code here.

## Template at a glance

| You get | You bring |
|---------|-----------|
| App + critic worker + Vue SPA | Neon project (pooled + direct URLs) |
| Pulumi infra + Deploy/Destroy workflows | GCP project + one-time local bootstrap |
| Google GIS auth on `main` | GitHub secrets/vars (see [`docs/DEPLOY.md`](./docs/DEPLOY.md)) |

Full checklist (bootstrap → deploy → destroy leftovers): [`docs/DEPLOY.md`](./docs/DEPLOY.md).

## Features

- Google Sign-In (GIS) + app JWT sessions (primary on `main`)
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
| DB (cloud) | Neon |
| Infra | Pulumi (Go) → GCP Cloud Run, Artifact Registry, Secret Manager, WIF |

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
cp .env.example .env       # VITE_GOOGLE_CLIENT_ID, VITE_API_BASE_URL
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
| `AUTH_SECRET` | JWT signing key (strong value required in production) |
| `DATABASE_URL` / `POSTGRES_*` | Database (Neon pooled URL in cloud) |
| `POSTGRES_SSLMODE` | DSN sslmode when not using `DATABASE_URL` (default `disable`) |
| `GOOGLE_OAUTH_CLIENT_ID` | Google GIS audience |
| `GEMINI_API_KEY` | Critic LLM key |
| `CRITIC_WORKER_ENABLED` | Start in-process worker (`BOT_WORKER_ENABLED` still works) |
| `CORS_ORIGIN` | Comma-separated SPA origins (required explicit in production) |
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

- **users** — username, optional password, `google_sub`, `email`, `is_bot`
- **posts** / **comments** — content + author
- **bot_jobs** — async critic queue
- Seeded bot user: `@vibe_critic` (`is_bot=true`)

## GCP deploy / destroy

See the template checklist: [`docs/DEPLOY.md`](./docs/DEPLOY.md).

- **Deploy:** push to `main` or `gh workflow run "Deploy to GCP"` — build → secret sync → `pulumi up` → goose migrate → Cloud Run bump.
- **Destroy:** `gh workflow run "Destroy GCP infrastructure" -f confirm=destroy` — default keeps Secret Manager; does **not** touch Neon or the GCS state bucket. Details: [`infra/README.md`](./infra/README.md).

Secret placement (GCP Secret Manager vs Pulumi vs GitHub): [`infra/SECRETS.md`](./infra/SECRETS.md).

Auth note: `main` uses Google GIS. Clerk migration is optional/unmerged ([PR #14](https://github.com/dakaii/vibegopher/pull/14)); prefer a private product repo for Clerk/Slack.

## License

[MIT](./LICENSE)
