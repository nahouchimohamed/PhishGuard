# 06 — GitHub Actions CI

## File

[`.github/workflows/ci.yml`](../../.github/workflows/ci.yml)

## What runs (and what does not)

| Job | Purpose |
|-----|---------|
| `go` | `go vet` + `go test` |
| `web` | `svelte-check` + Vitest |
| `docker` | Build api, worker, web images |
| `terraform` | `fmt -check` + `validate` |

There is **no deploy job**. That is intentional: you practice the safety net without paying for cloud hosting.

## Free minutes

GitHub-hosted runners are free for public repos (and have a monthly allotment for private ones). Builds run on GitHub’s machines, not yours.
