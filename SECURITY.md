# Security Policy

## What this project is

PhishGuard is a phishing simulation platform. It must not be used to phish real people, collect credentials, or send unsolicited email.

## Reporting a vulnerability

If you find a security issue in this repository (for example auth bypass, injection, or unsafe defaults in the local stack), please open a private report via GitHub Security Advisories for the repo, or contact the maintainer listed in `CODEOWNERS`.

Please include:

- Affected component (API, web, worker, compose, CI)
- Steps to reproduce on a local Compose stack
- Impact assessment (local-only vs anything that could escape the workstation)

## Safe local defaults

- Compose publishes services on `127.0.0.1` only.
- Seed passwords and Grafana `admin/admin` are for **local development**, not production.
- Terraform defaults target **LocalStack**, not real AWS. Do not apply this configuration to a production account.
- `apps/web/src/safety.test.ts` guards against password fields outside login and `{@html}` usage.

## Out of scope

- Using the simulator against people who did not consent
- Turning templates into real credential-harvesting pages
