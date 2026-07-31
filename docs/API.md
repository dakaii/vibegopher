# API reference

Base URL (local): `http://localhost:8081/api`

Authenticate with:

```
Authorization: Bearer <jwt>
```

`expiresIn` on auth responses is **seconds until expiry** (24h), not a Unix timestamp.

## Auth

| Method | Path | Auth | Notes |
|--------|------|------|-------|
| `POST` | `/auth/google` | no | Body: `{ "id_token": "..." }` — primary SPA login |
| `POST` | `/signup` | no | Password register — only if `ENABLE_PASSWORD_AUTH=true` |
| `POST` | `/login` | no | Password login — only if `ENABLE_PASSWORD_AUTH=true` |
| `GET` | `/me` | yes | Current user (`id`, `username`, `email`, `is_bot`) |

### `POST /auth/google`

```json
{ "id_token": "google-gis-credential" }
```

```json
{
  "tokenType": "Bearer",
  "token": "...",
  "expiresIn": 86400
}
```

### Password auth (dev/tests only)

```json
{ "username": "alice", "password": "Secret123" }
```

Username: 3–50 chars, `[a-zA-Z0-9_-]`, no leading/trailing `_`/`-`.  
Password: ≥8 chars, upper + lower + digit.

## Posts

| Method | Path | Auth |
|--------|------|------|
| `GET` | `/posts` | no |
| `POST` | `/posts` | yes |
| `GET` | `/posts/{id}` | no |
| `GET` | `/posts/user/{userId}` | no |
| `PATCH` | `/posts/{id}` | yes (owner) |
| `DELETE` | `/posts/{id}` | yes (owner) |

Create/update body:

```json
{ "content": "string, 1–280 chars" }
```

Creating a post enqueues a critic job (`bot_jobs`) unless the author is the bot account.

## Comments

| Method | Path | Auth |
|--------|------|------|
| `POST` | `/comments` | yes |
| `GET` | `/comments/post/{postId}` | no |
| `GET` | `/comments/user/{userId}` | no |
| `GET` | `/comments/{id}` | no |
| `PATCH` | `/comments/{id}` | yes (owner) |
| `DELETE` | `/comments/{id}` | yes (owner) |

```json
{ "content": "string, 1–280 chars", "post_id": "<uuid>" }
```

## Health

`GET /api/health` → `{ "ok": true }`

## Errors

```json
{ "error": "human-readable message" }
```

Production responses avoid leaking internal error details for 5xx failures.
