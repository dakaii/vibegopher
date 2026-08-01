# AGENTS.md

## Cursor Cloud specific instructions

VibeGopher = Go REST API (`cmd/server`, :8081) + Vue 3 SPA (`frontend/`, :5173) + Postgres, with an
async `@vibe_critic` worker. See `README.md`, `CONTRIBUTING.md`, and `docs/API.md` for product/API details.

The README/Makefile document a **Docker Compose** dev flow (`make up`, `make test`). Docker is **not**
available in this cloud VM — everything runs **natively** instead. The notes below capture the non-obvious
differences from the documented Docker flow.

### Toolchain
- Go 1.24 lives at `/usr/local/go/bin` (added to `PATH` in `~/.bashrc`). The repo requires Go 1.24; the
  system `/usr/bin/go` is 1.22 and fails with `toolchain not available`. Ensure `/usr/local/go/bin` is first
  on `PATH` (login shells already do this).
- Node 22 + npm are preinstalled. Frontend deps: `npm --prefix frontend install` (done by the update script).

### PostgreSQL (native, not Docker)
- Postgres 16 runs natively on `localhost:5432`, role `postgres` / password `postgres`.
- It is **not auto-started on boot** — start it each session: `sudo pg_ctlcluster 16 main start`
  (check with `pg_lsclusters`).
- Databases already created: `vibegopher_development` and `vibegopher_test`.
- Gotcha: `config/development.conf` / `config/test.conf` set `POSTGRES_HOST=postgresql-dev` /
  `postgresql-test` (Docker service names). When running natively you **must override
  `POSTGRES_HOST=localhost`** (and `POSTGRES_PORT=5432`). Do **not** point the app at those conf files as-is.

### Running the backend (natively)
```bash
export POSTGRES_HOST=localhost POSTGRES_PORT=5432 POSTGRES_USER=postgres POSTGRES_PASSWORD=postgres \
       POSTGRES_DB=vibegopher_development APP_ENV=development PORT=8081 \
       AUTH_SECRET=dev-secret-key-change-me CORS_ORIGIN=http://localhost:5173 \
       CRITIC_WORKER_ENABLED=true ENABLE_PASSWORD_AUTH=true HASH_COST=8
go run ./cmd/migrate up      # apply goose migrations (re-run after new migration files)
go run ./cmd/server          # API on http://localhost:8081
```
- The critic worker only produces real `@vibe_critic` replies when `GEMINI_API_KEY` is set; without it the
  worker still runs but `bot_jobs` fail (harmless for core testing).

### Running the frontend (natively)
- `npm --prefix frontend run dev` → `http://localhost:5173` (Vite proxies `/api` → `:8081`).
- **CORS gotcha:** leave `VITE_API_BASE_URL` **empty/unset** so the SPA calls relative `/api` through the
  Vite proxy (same-origin). The backend does **not** answer CORS preflight (`OPTIONS` → 404), so setting
  `VITE_API_BASE_URL=http://localhost:8081` (as `frontend/.env.example` does) breaks all SPA API calls. A
  fresh clone with no `frontend/.env` works because the default base URL is empty.

### Auth for local/manual testing
- The SPA login is Google Sign-In only and needs a real `VITE_GOOGLE_CLIENT_ID`. For local testing without
  Google OAuth, run the API with `ENABLE_PASSWORD_AUTH=true` and use `POST /api/signup` + `POST /api/login`
  (see `docs/API.md`). To view the authenticated SPA feed, put the returned JWT in `localStorage` under key
  `vibegopher_token`, then reload.

### Lint / test / build
- Backend: `go vet ./...`; unit tests `go test ./internal/... -count=1` (no DB needed).
- Backend integration tests live in `testing/` and need a migrated `vibegopher_test` DB with
  `APP_ENV=test ENABLE_PASSWORD_AUTH=true CORS_ORIGIN=*` and `POSTGRES_HOST=localhost`.
  - **Known pre-existing flake (not an environment issue):** `testing/factory` seeds faker via
    `time.Now().Unix()` (second resolution), so any test creating ≥2 users within the same wall-clock second
    collides on `users.username` (currently ~6 failures, e.g. `TestUpdatePostUnauthorized`). 21 tests pass.
- Frontend (from `frontend/`): `npm run lint` (Biome), `npm run typecheck`, `npm run build`.
