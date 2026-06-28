# Deploy (local Compose)

This folder holds the everyday local stack — not a cloud deploy.

## Files

- `docker-compose.yml` — API, worker, web, Postgres, Redis, ElasticMQ, Prometheus, Grafana
- `Dockerfile.api` / `Dockerfile.worker` / `Dockerfile.web` — multi-stage images
- `elasticmq.conf` — SQS-compatible queue + DLQ for deliveries

## Usage

From the repo root:

```bash
make up
# or
docker compose -f deploy/docker-compose.yml --env-file .env.example up --build -d
```

Published ports bind to `127.0.0.1` so the stack stays on the workstation.

## Notes

- Compose healthchecks hit `/healthz` on API and worker.
- Volumes persist Postgres/Redis across `make down`; use `make reset` to wipe.
- ElasticMQ speaks the SQS API; the Go client points `BaseEndpoint` at `http://elasticmq:9324`.
