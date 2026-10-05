# 06 — GitHub Actions CI

## Files

- [`.github/workflows/ci.yml`](../../.github/workflows/ci.yml)
- [`.github/dependabot.yml`](../../.github/dependabot.yml)
- [`.github/pull_request_template.md`](../../.github/pull_request_template.md)

## What runs (and what does not)

| Job | Purpose |
|-----|---------|
| `go` | `gofmt` check + `go vet` + `go test -race` |
| `web` | `npm ci` + `svelte-check` + Vitest + build |
| `docker` | Build api, worker, web images |
| `terraform` | `fmt -check` + `validate` (LocalStack-oriented module) |
| `contracts` | Python link, env, and metrics contract scripts |

There is **no deploy job**. That is intentional: you practice the safety net without paying for cloud hosting.

Dependabot watches GitHub Actions, Go modules, npm, Dockerfiles under `deploy/`, and Terraform weekly.

## Free minutes

GitHub-hosted runners are free for public repos (and have a monthly allotment for private ones). Builds run on GitHub’s machines, not yours.
