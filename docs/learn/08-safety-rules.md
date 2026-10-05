# 08 — Safety rules

PhishGuard is a **local training simulator**. Keep these invariants:

1. No real email is sent. Delivery is an inbox row written by the worker.
2. Simulation pages never collect passwords or credentials.
3. Templates are seeded and read-only.
4. Subjects stay prefixed with `[Training]` in seed data.
5. The yellow training banner stays visible on app pages.

The web safety test (`apps/web/src/safety.test.ts`) fails CI-style local runs if a password field appears outside login or if `{@html}` sneaks into templates.

## Reporting a concern

If you think a change weakens these rules, open a bug with the training-safety checkbox unchecked and point at the file. Prefer fixing the invariant in the same PR that introduced the risk.
