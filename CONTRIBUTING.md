# Contributing

Thanks for helping with VibeGopher. This repo is the **public Twitter-style demo** (Go API + Vue SPA + `@vibe_critic`). A future Slack SaaS product is intentionally **out of scope** here and should live in a private repository.

## Development setup

1. Install Docker, Docker Compose, Make, Go 1.24+, and Node 22+.
2. Backend:

   ```bash
   make build
   make up                 # postgres + migrate + API on :8081
   ```

3. Frontend:

   ```bash
   cd frontend
   cp .env.example .env    # set VITE_GOOGLE_CLIENT_ID
   npm install
   npm run dev             # http://localhost:5173
   ```

4. Tests:

   ```bash
   make test
   cd frontend && npm run lint && npm run typecheck
   ```

## Code guidelines

- Prefer small, focused changes with clear commit messages.
- Keep secrets out of git (see [`SECURITY.md`](./SECURITY.md)).
- Schema changes go through goose SQL in `db/migrations` — do not rely on GORM AutoMigrate.
- Password `/api/signup` and `/api/login` are for tests/local only (`ENABLE_PASSWORD_AUTH=true`). Production auth is Google Sign-In.
- Critic logic lives in `internal/critic`. Persona prompt source of truth: `internal/critic/persona.go`.
- Frontend: TypeScript + Biome; prefer `user.is_bot` over hard-coding usernames.

## Pull requests

1. Fork / branch from `main`.
2. Add or update tests when behavior changes.
3. Run `make test` and frontend `npm run lint && npm run typecheck`.
4. Open a PR with a short summary of *what* and *why*.

## Scope notes

| In this repo | Not in this repo |
|--------------|------------------|
| Public feed, posts, comments | Slack SaaS / billing / multi-tenant workspaces |
| `@vibe_critic` Gemini worker | Private critic-engine extraction (optional later) |
| Pulumi GCP deploy for the demo | Customer-specific infra |

If a change only makes sense for a paid Slack product, keep it out of this public tree.
