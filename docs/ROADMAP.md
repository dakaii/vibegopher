# Roadmap & missing features

VibeGopher is intentionally small: a public Twitter knockoff with an AI critic. Gaps below are prioritized for the demo, not for a full social network.

## Done recently

| Feature | Notes |
|---------|-------|
| Feed cursor pagination (`limit` + `before`) | Comments still hard-capped at 200 |
| Mute `@vibe_critic` (client `localStorage`) | Hide-only; jobs still enqueue |
| Per-IP rate limits on auth + writes | In-process; not shared across replicas |
| `X-Request-ID` + request access logs | Structured fields in `log.Printf` |
| Critic comment refresh poll | Reloads comments without remounting posts |

## Worth doing next (demo quality)

| Feature | Why |
|---------|-----|
| Comment pagination | Still `Limit(200)` |
| Server-side mute / skip enqueue | Save Gemini cost when muted |
| Separate Cloud Run service for the critic worker | Avoid coupling LLM latency to API process |
| SSE/WebSocket for critic replies | Better than 8s polling |
| Basic profile page (`/u/:username`) | Users exist but have no surface |
| Shared/redis rate limits | Needed if Cloud Run scales out |

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
