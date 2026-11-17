# PhishGuard

PhishGuard is a phishing simulation platform for security teams. Operators run campaigns against consented participants, track delivery and reporting, and review results in Grafana.

Messages never leave the platform. There is no outbound email and no credential collection. Templates are seeded and read-only.

**SvelteKit + TypeScript · Go · PostgreSQL · Redis · ElasticMQ (SQS) · Docker Compose · GitHub Actions · Prometheus/Grafana · Terraform/LocalStack (optional)**

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

## Architecture

See [docs/architecture.md](docs/architecture.md) for the runtime diagram and component boundaries. Local operations, health checks, and volume resets are in [docs/operations.md](docs/operations.md).

## Runtime

| Need | Approach |
|------|----------|
| Full stack | Docker Compose |
| Queue | ElasticMQ in Compose |
| Lambda runtime | LocalStack (optional) |
| CI | GitHub Actions |

Terraform in this repo targets LocalStack. Do not apply it to a production AWS account.

## Layout

```
backend/          Go API, worker, lambda entry
apps/web/         SvelteKit UI (BFF to the API)
deploy/           Compose + Dockerfiles + ElasticMQ config
infra/terraform/  LocalStack Lambda/SQS
monitoring/       Prometheus + Grafana
scripts/          Stdlib Python contract checks
docs/             Architecture and operations
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

PhishGuard must not be used to phish real people. Simulation links open an awareness notice with no credential fields. Subjects in seed data stay prefixed with `[Training]`.

## License

Released under the [MIT License](LICENSE).
