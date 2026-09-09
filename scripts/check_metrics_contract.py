#!/usr/bin/env python3
"""Keep Prometheus metric names aligned across Go, Grafana, and rules (Python 3.8+)."""

from __future__ import annotations

import json
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
METRICS_GO = ROOT / "backend" / "internal" / "metrics" / "metrics.go"
DASHBOARD = ROOT / "monitoring" / "grafana" / "dashboards" / "phishguard.json"
RULES = ROOT / "monitoring" / "prometheus" / "rules" / "phishguard.rules.yml"

NAME_RE = re.compile(r'Name:\s*"(phishguard_[a-z0-9_]+)"')
METRIC_TOKEN_RE = re.compile(r"phishguard_[a-z0-9_]+")


def metric_names_from_go():
    text = METRICS_GO.read_text(encoding="utf-8")
    names = NAME_RE.findall(text)
    if not names:
        raise SystemExit("check_metrics_contract: no metric Names found in metrics.go")
    return sorted(set(names))


def collect_promql_strings(obj, out):
    if isinstance(obj, dict):
        for key, value in obj.items():
            if key == "expr" and isinstance(value, str):
                out.append(value)
            else:
                collect_promql_strings(value, out)
    elif isinstance(obj, list):
        for item in obj:
            collect_promql_strings(item, out)


def main():
    names = metric_names_from_go()
    dashboard = json.loads(DASHBOARD.read_text(encoding="utf-8"))
    exprs = []
    collect_promql_strings(dashboard, exprs)
    dashboard_text = "\n".join(exprs)
    rules_text = RULES.read_text(encoding="utf-8") if RULES.is_file() else ""

    errors = []
    for name in names:
        # Counters/histograms appear as the base name or with _bucket/_count/_sum.
        base = name
        found_dash = base in dashboard_text or any(
            token.startswith(base) for token in METRIC_TOKEN_RE.findall(dashboard_text)
        )
        if not found_dash:
            errors.append("dashboard missing metric {!r}".format(name))

    # Every phishguard_* token in dashboard/rules should exist in Go (base name).
    go_bases = set(names)
    for source, text in (("dashboard", dashboard_text), ("rules", rules_text)):
        for token in sorted(set(METRIC_TOKEN_RE.findall(text))):
            base = token
            for suffix in ("_bucket", "_count", "_sum"):
                if base.endswith(suffix):
                    base = base[: -len(suffix)]
                    break
            if base not in go_bases:
                errors.append("{} references unknown metric {!r}".format(source, token))

    if "refresh" not in dashboard or dashboard.get("refresh") != "10s":
        errors.append("dashboard refresh should be '10s'")

    if errors:
        print("check_metrics_contract: FAILED")
        for err in errors:
            print("  - {}".format(err))
        return 1

    print(
        "check_metrics_contract: OK ({} metric(s), dashboard refresh={})".format(
            len(names), dashboard.get("refresh")
        )
    )
    return 0


if __name__ == "__main__":
    sys.exit(main())
