# VibeGopher web (Vue 3 + TypeScript 7)

Sign in with Google, browse the feed, post, and comment. `@vibe_critic` replies asynchronously from the Go bot worker.

Tooling:

- **TypeScript 7** for app types / `tsc`
- **Biome** for lint + format
- **vue-tsc** typechecks SFCs via `@typescript/typescript6` (TS 7 does not yet export the classic Node API `vue-tsc` needs)

## Setup

```bash
cp .env.example .env
# set VITE_GOOGLE_CLIENT_ID (same OAuth client ID as the API)
# set VITE_API_BASE_URL=http://localhost:8081

npm install
npm run dev
```

Open `http://localhost:5173`.

```bash
npm run lint
npm run typecheck
npm run build
```

## Google Cloud OAuth

1. Create an OAuth 2.0 Client ID (Web application).
2. Authorized JavaScript origins: `http://localhost:5173` (and your prod frontend origin).
3. Use the **Client ID** as both `VITE_GOOGLE_CLIENT_ID` and API `GOOGLE_OAUTH_CLIENT_ID`.

## Build

```bash
npm run build
```

Static output lands in `dist/` — host later on Firebase Hosting or Cloud Storage + CDN.

## Note on TypeScript 7 + vue-tsc

The project depends on **typescript@7**. Vue SFC typechecking still goes through `vue-tsc`, which needs the classic Node `typescript/lib/tsc` export. Until that lands in TS 7.x, `npm run typecheck` bridges via `@typescript/typescript6` (`scripts/vue-tsc.cjs`). App code and editor tooling still target TypeScript 7.
