# Operations

## Start and stop

```bash
docker compose -f deploy/docker-compose.yml --env-file .env.example up --build -d
```

Or: `make up`, `make down` (keeps volumes), `make stop`, `make reset` (deletes volumes).

| URL | What |
|-----|------|
| http://localhost:5173 | Web UI |
| http://localhost:8080/healthz | API health (Postgres and Redis) |
| http://localhost:9090 | Prometheus |
| http://localhost:3002 | Grafana (anonymous viewer, or admin/admin) |

## Campaign check

1. Sign in as `admin@phishguard.local` / `phishguard-dev`.
2. Create a campaign with the password-expiry template and all employees.
3. Sign out, then sign in as `alice@phishguard.local` and report the inbox message.
4. Sign in as `bob@phishguard.local`, open the message, and open the simulation link.
5. As admin, open the campaign and compare Reported vs Clicked.

Seed subjects are prefixed with `[Training]`. The banner stays visible. The simulation page has no password field.

## Troubleshooting

- Port clashes on `5432` or `6379`: stop local Postgres/Redis or change the published ports in `deploy/docker-compose.yml`.
- Schema changes in `backend/migrations/001_init.sql` apply only on a fresh volume. Use `make reset` to wipe and recreate.
- `make down` keeps volumes; `make reset` removes them.

Compose marks API and worker healthy when `GET /healthz` succeeds. If the web UI loads but login fails, open http://localhost:8080/healthz first. A failing healthcheck usually means Postgres or Redis is still starting.
