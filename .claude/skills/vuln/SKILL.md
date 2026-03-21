---
name: vuln
description: Scan all Go packages for known vulnerabilities using govulncheck
---

Run `task vuln` to check all packages for known CVEs using govulncheck.

## When to use

- After adding or upgrading a dependency in `go.mod` / `go.sum`.
- When the user explicitly asks for a security or vulnerability audit.
- Periodically as part of a release checklist.

## When NOT to use

- During regular feature development when no dependencies were changed.
- After every single edit — this is not a routine step like format or lint.
