# Monitoring

Local Prometheus + Grafana for the training stack.

## Layout

- `prometheus/prometheus.yml` — scrape API `:8080` and worker `:8081` every 10s
- `prometheus/rules/` — recording/alert rules loaded by Prometheus
- `grafana/provisioning/` — datasource + dashboard providers
- `grafana/dashboards/phishguard.json` — training dashboard

## URLs (Compose)

| Service | URL |
|---------|-----|
| Prometheus | http://localhost:9090 |
| Grafana | http://localhost:3002 |

Grafana allows anonymous viewer access for classroom demos; `admin`/`admin` is also seeded for local learning only.

## Habit

When you add a counter in Go, update the dashboard (or at least the metrics contract script) in the same change so the lesson stays consistent.
