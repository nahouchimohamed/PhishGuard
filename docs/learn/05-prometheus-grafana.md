# 05 — Metrics you can see

## Files

- [`backend/internal/metrics/metrics.go`](../../backend/internal/metrics/metrics.go)
- [`monitoring/prometheus/prometheus.yml`](../../monitoring/prometheus/prometheus.yml)
- [`monitoring/prometheus/rules/phishguard.rules.yml`](../../monitoring/prometheus/rules/phishguard.rules.yml)
- [`monitoring/grafana/dashboards/phishguard.json`](../../monitoring/grafana/dashboards/phishguard.json)

## Path

1. Launch a campaign → `phishguard_campaign_launches_total` increments in the API.
2. Worker delivers → `phishguard_deliveries_total` on the worker (`:8081/metrics`).
3. Report / click → counters on the API.

Prometheus scrapes `api:8080` and `worker:8081` every **10s**, evaluates alert rules from `rules/phishguard.rules.yml`, and Grafana refreshes on the same interval.

Dashboard panels include totals, report/click rates, HTTP request rate, 5xx ratio, p95 latency, scrape `up`, and `go_goroutines`.

Open http://localhost:3002 → folder **PhishGuard** → dashboard **PhishGuard**.

## Habit

When you add a feature, ask: “Which counter proves it happened?” Add the counter next to the business logic, not only in a dashboard later. Keep dashboard queries aligned with names in `metrics.go` (see `scripts/check_metrics_contract.py`).

## Alert rules

`monitoring/prometheus/rules/phishguard.rules.yml` evaluates alongside scrapes. Treat rules as classroom signals (for example sustained 5xx) rather than production paging — the point is to see a rule fire locally after you break something on purpose.
