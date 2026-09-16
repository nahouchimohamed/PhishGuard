# 02 — Campaign service and SQL

## Files

- [`backend/internal/campaign/service.go`](../../backend/internal/campaign/service.go) — launch + deliver rules
- [`backend/internal/store/store.go`](../../backend/internal/store/store.go) — SQL next to each method
- [`backend/migrations/001_init.sql`](../../backend/migrations/001_init.sql) — schema

## Commit-then-publish

`Launch` does this in order:

1. Insert campaign with status `launching` and all recipients (one transaction).
2. Publish one SQS message per recipient.
3. Only if publish succeeds for the batch, mark the campaign `active`.

If step 2 fails partway, the campaign stays `launching`. The admin UI shows **Retry queue publish**, which re-publishes undelivered recipients. That is the failure mode this pattern exists for.

## Idempotent delivery

```sql
UPDATE campaign_recipients
SET delivered_at = now()
WHERE id = $1 AND delivered_at IS NULL
```

A duplicate SQS message returns `RowsAffected = 0`. The worker still deletes the message. No second inbox row appears because there is only one recipient row (`UNIQUE (campaign_id, user_id)`).

## Timestamps, not a single status

`opened_at`, `clicked_at`, and `reported_at` are separate columns. An employee can open, click, **and** report. The UI label prefers Reported — see `RecipientLabel`.

## Stats cache

Campaign totals are expensive to recompute on every admin refresh, so `GetStats` checks Redis first and falls back to SQL on a miss. After a delivery, open, click, or report, the service invalidates that campaign's cache key so the next read rebuilds fresh numbers. When Redis is down, the SQL path still works — the cache is an accelerator, not a source of truth.
