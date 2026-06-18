# PhishGuard

PhishGuard is a phishing simulation platform for security teams:

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
| API metrics | http://localhost:8080/metrics |
| Worker metrics | http://localhost:8081/metrics |
| ElasticMQ UI | http://localhost:9325 |
| Postgres | localhost:5432 |
| Redis | localhost:6379 |
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
```

## Make targets

```bash
make help             # list targets
make up               # compose up --build -d
make down             # compose down (keeps volumes)
make stop             # stop containers
make reset            # compose down -v (wipes DB)
make logs             # follow logs
make api-test         # go test ./...
make api-check        # gofmt + vet + test
make web-test         # vitest
make web-check        # svelte-check + vitest
make ci-local         # api-check + web-check
make compose-config   # validate compose file
make lambda-zip       # build Linux bootstrap for LocalStack
make check-docs       # Markdown links + .env.example refs
make check-static     # check-docs + metrics contract
```

## Safety

PhishGuard is a **simulation platform**. It must not be used to phish real people. Simulation links open a lesson page with no credential fields.
