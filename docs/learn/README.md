# Learn PhishGuard

Follow in order. Each note points at the real files.

1. [Run the app](01-run-the-app.md)
2. [Campaign service and SQL](02-campaign-and-sql.md)
3. [Worker and SQS](03-worker-and-sqs.md)
4. [SvelteKit BFF](04-sveltekit-bff.md)
5. [Prometheus + Grafana](05-prometheus-grafana.md)
6. [GitHub Actions](06-github-actions.md)
7. [LocalStack + Terraform](07-localstack-terraform.md) (optional)

## Architecture (everyday Compose)

```mermaid
flowchart LR
  browser[Browser on localhost]
  web[web :5173]
  api[api :8080]
  worker[worker :8081]
  pg[(postgres)]
  redis[(redis)]
  mq[elasticmq SQS]
  prom[prometheus]
  graf[grafana]

  browser --> web
  web -->|BFF| api
  api --> pg
  api --> redis
  api -->|enqueue| mq
  worker -->|long poll| mq
  worker --> pg
  worker --> redis
  prom -->|scrape| api
  prom -->|scrape| worker
  graf --> prom
```

All published ports bind to `127.0.0.1`. Training only — no real phishing.
