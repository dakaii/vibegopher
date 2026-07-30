# VibeGopher web (Vue 3)

Sign in with Google, browse the feed, post, and comment. `@vibe_critic` replies asynchronously from the Go bot worker.

## Setup

```bash
cp .env.example .env
# set VITE_GOOGLE_CLIENT_ID (same OAuth client ID as the API)
# set VITE_API_BASE_URL=http://localhost:8081

npm install
npm run dev
```

Open `http://localhost:5173`.

## Google Cloud OAuth

1. Create an OAuth 2.0 Client ID (Web application).
2. Authorized JavaScript origins: `http://localhost:5173` (and your prod frontend origin).
3. Use the **Client ID** as both `VITE_GOOGLE_CLIENT_ID` and API `GOOGLE_OAUTH_CLIENT_ID`.

## Build

```bash
npm run build
```

Static output lands in `dist/` — host later on Firebase Hosting or Cloud Storage + CDN.
