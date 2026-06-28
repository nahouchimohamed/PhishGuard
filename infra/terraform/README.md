# Terraform (LocalStack only)

These modules provision Lambda and SQS against [LocalStack](https://localstack.cloud/), not real AWS.

## Files

- `versions.tf` / `providers.tf` — Terraform and AWS provider pinned to LocalStack endpoints
- `variables.tf` — auth token, queue names, Lambda env (Postgres/Redis via `host.docker.internal`)
- `main.tf` — queue, DLQ, Lambda function, event source mapping
- `outputs.tf` — queue URL and function name for smoke checks

## Safety

Do **not** `terraform apply` this configuration to a production AWS account. Provider endpoints default to LocalStack (`http://localhost:4566`).

## Quick path

1. Export `LOCALSTACK_AUTH_TOKEN` and start LocalStack.
2. `make lambda-zip` to build `backend/bin/bootstrap`.
3. `terraform init` then `terraform apply` in this directory.

