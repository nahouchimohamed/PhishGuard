## Summary

<!-- What changed and why (1–3 bullets). -->

-

## Test plan

- [ ] `make api-check` (or CI `go` job)
- [ ] `make web-check` (or CI `web` job)
- [ ] `make check-static` (or CI `contracts` job) when docs/env/metrics change
- [ ] Compose / Docker builds still make sense for touched images
- [ ] Docs / learn notes updated if behavior or ports changed

## Safety

- [ ] No real phishing / credential-harvesting flows introduced
- [ ] Training banners and simulation-only messaging preserved where relevant
