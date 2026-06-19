# PhishGuard web (SvelteKit BFF)

## Scripts

```bash
npm run dev
npm run check
npm test
npm run build
npm run format
npm run format:check
```

## Env

- `API_BASE_URL` — Go API base (default `http://localhost:8080`)

## Routes

| Path | Role |
|------|------|
| `/login` | public |
| `/campaigns` | admin |
| `/campaigns/new` | admin |
| `/campaigns/[id]` | admin |
| `/inbox` | employee |
| `/inbox/[id]` | employee |
| `/sim/[token]` | public simulation result (no credential fields) |

## Offline tests

`npm run check` and `npm test` do not need Docker, Postgres, Redis, or SQS.
