---
name: lint
description: Run golangci-lint across the entire project
---

Run `task lint` to lint all Go source files using golangci-lint.

## When to use

Use this skill **automatically** at the end of every session in which one or more `.go` files were created or modified — right after `/format`, before presenting the result to the user for review.

Fix every linter error before handing off. Do not leave the output unread; if lint fails, address the issues and re-run.

## When NOT to use

- When only non-Go files were changed.
- When no files were modified at all.
