# 07 — LocalStack + Terraform (optional)

Everyday Compose already teaches SQS via ElasticMQ. This chapter adds **Lambda** and **Terraform** without a paid AWS machine.

## Free tier

[LocalStack Hobby](https://www.localstack.cloud/pricing) is free for non-commercial use. As of 2026 you need a free account and `LOCALSTACK_AUTH_TOKEN`.

Do **not** apply this Terraform to real AWS. The provider endpoints are hard-wired to LocalStack defaults.

## Local ports (Compose)

Compose publishes web/api/worker/prometheus/grafana/elasticmq/postgres/redis on **`127.0.0.1` only**. That keeps the training stack off the LAN while you still open the usual `http://localhost:...` URLs from your machine. LocalStack Lambda reaches Postgres/Redis via `host.docker.internal` (see `infra/terraform` variables).

## Steps

1. Create a LocalStack account and copy your auth token into `.env` as `LOCALSTACK_AUTH_TOKEN`.
2. Start LocalStack (Docker). Example:

```bash
docker run --rm -it -p 4566:4566 -e LOCALSTACK_AUTH_TOKEN -v /var/run/docker.sock:/var/run/docker.sock localstack/localstack
```

3. Build the Lambda binary (Linux):

```bash
make lambda-zip
```

Or:

```bash
cd backend
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -o bin/bootstrap ./cmd/lambda
```

4. Apply Terraform against LocalStack only:

```bash
cd infra/terraform
terraform init
terraform apply
```

The module attaches `AWSLambdaSQSQueueExecutionRole`, zips `backend/bin/bootstrap` with mode `0755`, and sets Lambda timeout to **5s**. URLs come from Terraform variables (`database_url`, `redis_url`, `public_app_url`).

5. Stop the Compose **worker** if you point the same queue at LocalStack, otherwise two consumers fight.

## What you learn

- IAM role + managed SQS execution policy + Lambda + SQS event source mapping as code
- `provided.al2023` custom runtime (the `bootstrap` binary)
- Why endpoint overrides and `default_tags` exist on the AWS provider

## Cleanup

Destroy LocalStack resources with `terraform destroy` in `infra/terraform` when you are done. Leaving a Hobby LocalStack container running is fine for a day of practice, but tear down Terraform state so the next apply starts clean.
