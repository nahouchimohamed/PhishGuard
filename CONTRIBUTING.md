# Contributing to PhishGuard

PhishGuard is a phishing simulation platform. Keep changes safe and runnable on a workstation with Docker Compose.

## Ground rules

- Do not add real email sending, credential-collection pages, or anything that could phish real people.
- Keep simulation subjects prefixed with `[Training]` and preserve the persistent simulation banner and awareness notice.
- Prefer small, reviewable PRs that touch one concern (backend, web, compose, docs, CI).

## Setup

1. Copy `.env.example` values as needed (Compose uses `--env-file .env.example` by default).
2. Start the stack: `make up` (or the `docker compose` command in the README).
3. Seed logins use password `phishguard-dev` (see README).

Published ports bind to `127.0.0.1` only. See [docs/operations.md](docs/operations.md) and the compose file.

## Checks

```bash
make api-check      # gofmt + vet + test
make web-check      # svelte-check + vitest
make ci-local       # both
make compose-config # validate compose YAML (needs Docker CLI)

# Optional static checks (stdlib Python 3.8+)
make check-static
# or: python scripts/check_links.py && python scripts/check_env_example.py && python scripts/check_metrics_contract.py
```

## Docs

Update [docs/architecture.md](docs/architecture.md) and [docs/operations.md](docs/operations.md) when ports, metrics, CI jobs, or the LocalStack flow change. Keep [SECURITY.md](SECURITY.md) accurate when safety invariants move.

## Pull requests

Use the PR template. Link related issues if any. Keep the safety checkbox honest.
