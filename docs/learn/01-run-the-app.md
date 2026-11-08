# 01 — Run the app

## Goal

Start the whole stack and click through both roles.

## Commands

```bash
# From the repo root (PowerShell or bash)
docker compose -f deploy/docker-compose.yml --env-file .env.example up --build -d
```

Or: `make up`

Open:

| URL | What |
|-----|------|
| http://localhost:5173 | Web UI |
| http://localhost:8080/healthz | API health |
| http://localhost:9090 | Prometheus |
| http://localhost:3002 | Grafana (anonymous viewer, or admin/admin) |

## Walkthrough

1. Log in as `admin@phishguard.local` / `phishguard-dev`.
2. Create a campaign with the password-expiry template and all employees.
3. Log out, log in as `alice@phishguard.local`.
4. Open the inbox message. Report it.
5. Log in as `bob@phishguard.local`, open the message, click the training link.
6. Back as admin, open the campaign detail and compare Reported vs Clicked.

## What to notice

- Subjects are prefixed with `[Training]`.
- The yellow banner never goes away — this is intentional.
- The simulation page has **no password field**.
