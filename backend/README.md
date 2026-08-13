# PhishGuard backend

Go API, worker, and Lambda entrypoints for the local phishing-awareness drill.

## Offline checks

```bash
gofmt -l .
go vet ./...
go test ./...
```

Unit tests are pure (no Postgres, Redis, or SQS required).

## Layout

- `cmd/api` — HTTP API + `/healthz` + `/metrics`
- `cmd/worker` — SQS consumer that "delivers" inbox messages
- `cmd/lambda` — LocalStack-oriented SQS handler
- `internal/` — auth, campaign, store, queue, cache, httpapi, metrics, config
- `migrations/` — SQL applied by Postgres init on a fresh volume
