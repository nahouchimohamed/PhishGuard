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
