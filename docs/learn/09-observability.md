# 09 — Observability habits

## Files

- [`backend/internal/metrics/metrics.go`](../../backend/internal/metrics/metrics.go)
- [`monitoring/prometheus/prometheus.yml`](../../monitoring/prometheus/prometheus.yml)
- [`monitoring/grafana/dashboards/phishguard.json`](../../monitoring/grafana/dashboards/phishguard.json)
- [`scripts/check_metrics_contract.py`](../../scripts/check_metrics_contract.py)

## What to watch while training

1. Launch a campaign and confirm `phishguard_campaign_launches_total` moves on the API.
2. Watch `phishguard_deliveries_total` on the worker after the queue drains.
3. Report and click from employee accounts; totals and rates update on Grafana within about 10s.

## Contract check

`make check-static` runs a small Python script that fails if dashboard PromQL references a metric name that is not registered in `metrics.go`. That keeps lesson dashboards honest when you rename a counter.

## Habit

Prefer one counter next to the business event over a dashboard-only query. Latency histograms and the HTTP 5xx ratio help you notice regressions during local drills without a production APM stack.
