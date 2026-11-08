# PhishGuard

Local phishing-**awareness** drill for learning the HENNGE Tadrill-shaped stack:

**SvelteKit + TypeScript · Go · PostgreSQL · Redis · ElasticMQ (SQS) · Docker Compose · GitHub Actions · Prometheus/Grafana · Terraform/LocalStack (optional)**

No real email. No password-collection pages. Templates are seeded and read-only.

## Quick start

Requirements: Docker Desktop (or Docker Engine + Compose).

```bash
docker compose -f deploy/docker-compose.yml --env-file .env.example up --build -d
```

| Service | URL |
|---------|-----|
| Web | http://localhost:5173 |
| API | http://localhost:8080/healthz |
| Prometheus | http://localhost:9090 |
| Grafana | http://localhost:3002 (admin/admin or anonymous viewer) |

### Seed logins

Password for every seeded user: `phishguard-dev`

| Email | Role |
|-------|------|
| admin@phishguard.local | admin |
| alice@phishguard.local | employee |
| bob@phishguard.local | employee |
| cara@phishguard.local | employee |

## Learn path

Read these in order under [`docs/learn/`](docs/learn/):

1. Run the app
2. Campaign service and SQL
3. Worker and SQS
4. SvelteKit BFF
5. Prometheus + Grafana
6. GitHub Actions
7. LocalStack + Terraform (optional, free Hobby plan)

## Free cloud alternatives (no paid EC2)

| Need | Free approach used here |
|------|-------------------------|
| Everyday stack | Docker Compose on your laptop |
| SQS | ElasticMQ in Compose |
| Lambda + Terraform practice | [LocalStack Hobby](https://www.localstack.cloud/pricing) (free auth token) |
| CI | GitHub Actions (no deploy job) |
| Real AWS later | Lambda/SQS are pay-per-request and usually free-tier; this repo never applies Terraform to real AWS |

## Layout

```
backend/          Go API, worker, lambda entry
apps/web/         SvelteKit UI (BFF to the API)
deploy/           Compose + Dockerfiles + ElasticMQ config
infra/terraform/  LocalStack-oriented Lambda/SQS
monitoring/       Prometheus + Grafana
docs/learn/       Short lessons tied to real files
```

## Make targets

```bash
make up               # compose up --build -d
make down             # compose down -v
make logs             # follow logs
make api-test         # go test ./...
make web-test         # vitest
```

## Safety

PhishGuard is a **training simulator**. It must not be used to phish real people. Simulation links open a lesson page with no credential fields.
