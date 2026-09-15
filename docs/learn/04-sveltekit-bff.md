# 04 — SvelteKit as a BFF

## Files

- [`apps/web/src/lib/server/api.ts`](../../apps/web/src/lib/server/api.ts)
- [`apps/web/src/routes/campaigns/new/+page.server.ts`](../../apps/web/src/routes/campaigns/new/+page.server.ts)
- [`apps/web/src/hooks.server.ts`](../../apps/web/src/hooks.server.ts)

## Pattern

The browser talks only to SvelteKit. Server `load` functions and form actions call the Go API and forward the session cookie. Benefits for learning:

- No CORS setup on day one.
- Secrets and `API_BASE_URL` stay server-side.
- Form actions feel like classic HTML forms.

## Login cookie bridge

Go sets `Set-Cookie`, but a server-to-server `fetch` does not automatically become a browser cookie. Login returns `session_token` in JSON so the BFF can `cookies.set(...)` with `httpOnly`.

## Session, guards, and errors

- [`apps/web/src/lib/server/session.ts`](../../apps/web/src/lib/server/session.ts) owns the cookie name and max-age.
- [`apps/web/src/lib/server/guard.ts`](../../apps/web/src/lib/server/guard.ts) centralizes role redirects.
- [`apps/web/src/lib/server/errors.ts`](../../apps/web/src/lib/server/errors.ts) maps API `{error}` strings to friendlier copy.
- `hooks.server.ts` clears the cookie only on HTTP 401 from `/api/me`, so a temporary API outage does not log everyone out.

## Preload caveat

SvelteKit hover preload can hit `/inbox/[id]` and `/sim/[token]`, which record opens and clicks in the API. Those links set `data-sveltekit-preload-data="off"`.
