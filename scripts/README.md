# Scripts

Stdlib Python 3.8+ contract checks used by `make check-docs` and `make check-static`.

| Script | Purpose |
|--------|---------|
| `check_links.py` | Verify relative Markdown links under `docs/` resolve |
| `check_env_example.py` | Ensure compose/env docs mention keys present in `.env.example` |
| `check_metrics_contract.py` | Ensure Grafana PromQL metric names exist in Go metrics |

Run from the repo root:

```bash
make check-docs
make check-static
# or
python scripts/check_links.py
python scripts/check_env_example.py
python scripts/check_metrics_contract.py
```

These scripts intentionally avoid third-party packages so CI and Windows laptops can run them the same way.
