# 07 — LocalStack + Terraform (optional)

Everyday Compose already teaches SQS via ElasticMQ. This chapter adds **Lambda** and **Terraform** without a paid AWS machine.

## Free tier

[LocalStack Hobby](https://www.localstack.cloud/pricing) is free for non-commercial use. As of 2026 you need a free account and `LOCALSTACK_AUTH_TOKEN`.

Do **not** apply this Terraform to real AWS. The provider endpoints are hard-wired to LocalStack defaults.

## Steps

1. Create a LocalStack account and copy your auth token into `.env` as `LOCALSTACK_AUTH_TOKEN`.
2. Start LocalStack (Docker). Example:

```bash
docker run --rm -it -p 4566:4566 -e LOCALSTACK_AUTH_TOKEN -v /var/run/docker.sock:/var/run/docker.sock localstack/localstack
```

3. Build the Lambda binary (Linux):

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

5. Stop the Compose **worker** if you point the same queue at LocalStack, otherwise two consumers fight.

## What you learn

- IAM role + Lambda + SQS event source mapping as code
- `provided.al2023` custom runtime (the `bootstrap` binary)
- Why endpoint overrides exist in the AWS provider
