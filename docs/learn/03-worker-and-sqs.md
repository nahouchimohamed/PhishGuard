# 03 — Worker loop and ElasticMQ

## Files

- [`backend/cmd/worker/main.go`](../../backend/cmd/worker/main.go)
- [`backend/internal/queue/queue.go`](../../backend/internal/queue/queue.go)
- [`deploy/elasticmq.conf`](../../deploy/elasticmq.conf)

## Why ElasticMQ?

Everyday Compose needs an SQS-shaped queue without an AWS account. ElasticMQ speaks the SQS API. The Go client uses `aws-sdk-go-v2` with `BaseEndpoint` pointed at `http://elasticmq:9324`.

## Loop

1. Long-poll `ReceiveMessage` (wait up to **10s**, matching `receiveMessageWait` in `elasticmq.conf`).
2. Call `campaign.Deliver` (idempotent).
3. `DeleteMessage` on success.
4. On deliver failure, leave the message — visibility timeout returns it. After 3 receives, ElasticMQ moves it to the DLQ.

## Health

The worker exposes `GET /healthz` and `/metrics` on `:8081` (Compose healthcheck uses `/healthz`).

## Lambda vs worker

[`backend/cmd/lambda`](../../backend/cmd/lambda) wraps the **same** `Deliver` function. Do not attach both consumers to one queue in Compose — they would race and split traffic, which hides the lesson. Use Lambda only in the LocalStack chapter.
