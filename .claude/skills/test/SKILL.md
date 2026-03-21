---
name: test
description: Run all Go unit tests via go test
---

Run `task test` to execute all Go tests.

## When to use

Use this skill **automatically** after any change to business logic — new features, bug fixes, refactors, or new/updated domain rules. Run it after `/format` and `/lint` pass, as the final gate before handing off.

Also use it explicitly when the user asks to verify tests, check coverage, or debug a failing test.

## When NOT to use

- When only documentation, configuration, or migration files were changed and no Go code was touched.
- When the user explicitly asks to skip tests (rare; they should acknowledge the risk).
