# VibeGopher web (Vue 3 + TypeScript 7)

Sign in with Clerk, browse the feed, post, and comment. `@vibe_critic` replies asynchronously from the Go bot worker.

Tooling:

- **TypeScript 7** for app types / `tsc`
- **Biome** for lint + format
- **vue-tsc** typechecks SFCs via `@typescript/typescript6` (TS 7 does not yet export the classic Node API `vue-tsc` needs)

## Setup

```bash
cp .env.example .env
# set VITE_CLERK_PUBLISHABLE_KEY from the Clerk Dashboard
# leave VITE_API_BASE_URL empty so Vite proxies /api → :8081

npm install
npm run dev
```

Open `http://localhost:5173`.

```bash
npm run lint
npm run typecheck
npm run build
```

## Clerk

1. Create a Clerk application at [dashboard.clerk.com](https://dashboard.clerk.com).
2. Copy the **Publishable key** into `VITE_CLERK_PUBLISHABLE_KEY`.
3. Put the matching **Secret key** in the API env as `CLERK_SECRET_KEY`.
4. Allow `http://localhost:5173` (and your prod origin) in the Clerk Dashboard.
5. Optional: enable Google / Apple SSO connections inside Clerk.

## Build

```bash
npm run build
```

Static output lands in `dist/` — host later on Firebase Hosting or Cloud Storage + CDN.

## Note on TypeScript 7 + vue-tsc

The project depends on **typescript@7**. Vue SFC typechecking still goes through `vue-tsc`, which needs the classic Node `typescript/lib/tsc` export. Until that lands in TS 7.x, `npm run typecheck` bridges via `@typescript/typescript6` (`scripts/vue-tsc.cjs`). App code and editor tooling still target TypeScript 7.
