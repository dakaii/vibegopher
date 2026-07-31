# Roadmap & missing features

VibeGopher is intentionally small: a public Twitter knockoff with an AI critic. Gaps below are prioritized for the demo, not for a full social network.

## Worth doing next (demo quality)

| Feature | Why |
|---------|-----|
| Cursor-based pagination for feed/comments | Hard `Limit(100/200)` will hide history |
| Mute / hide `@vibe_critic` per user | Persona docs already call this out |
| Rate limiting on auth + write endpoints | Abuse protection for a public deploy |
| Structured logging + request IDs | Debug Cloud Run / worker failures |
| Separate Cloud Run service for the critic worker | Avoid coupling LLM latency to API process |
| Frontend refresh of critic replies via SSE/WebSocket | Polling works; push is nicer |
| Basic profile page (`/u/:username`) | Users exist but have no surface |
| Empty / error states polish in the SPA | First-run UX |

## Nice-to-have social features

- Likes / reactions
- Follow graph + home timeline vs global feed
- Media uploads (needs object storage + virus scanning)
- Edit history / soft-delete UX
- Notifications

## Critic improvements

- Topic skip-lists in code (not “shadow ban” language in-prompt)
- Dead-letter / max-age for stuck `bot_jobs`
- Vertex AI via ADC instead of raw API keys
- Opt-out flag on posts (“don’t summon the critic”)
- Stronger link-fetch quotas and robots.txt respect

## Explicitly out of scope for this public repo

- Slack SaaS, workspace billing, seat management
- Multi-tenant private deployments as a product
- Content takedown / moderation tooling that suppresses posts

Those belong in a **private** repository if/when the critic is productized for Slack.

## Docs / ops hygiene already expected

- Goose migrations only (no AutoMigrate drift)
- Secrets in GCP Secret Manager + GitHub Actions sync ([`infra/SECRETS.md`](../infra/SECRETS.md))
- Password auth gated behind `ENABLE_PASSWORD_AUTH` and forbidden in production
