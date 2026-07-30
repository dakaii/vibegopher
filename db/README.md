# Database migrations

We use **[goose](https://github.com/pressly/goose)** for versioned SQL migrations.

| Concern | Tool |
|---------|------|
| Schema changes | goose (`db/migrations`) |
| Queries / models | GORM |
| Apply in CI / Neon | `go run ./cmd/migrate up` (deploy workflow) |
| Apply locally | `make migrate` |
| Cloud Run startup | **does not** migrate (avoids multi-instance races) |

## Commands

```bash
make migrate                 # up
make migrate-status
make migrate-down            # one step (prefer fix-forward in prod)
make create-migration NAME=add_google_sub
```

## Why not GORM AutoMigrate / Atlas?

- **GORM AutoMigrate** is fine for throwaway prototypes, not for production schema ownership (weak destructive changes, no reviewable history).
- **Atlas** is modern and powerful, but heavier than this app needs (declarative schema + migration dir + checksums + extra Docker image).
- **goose** is the boring default for small/medium Go + Postgres apps: plain SQL, embeddable, advisory locks, easy CI.
