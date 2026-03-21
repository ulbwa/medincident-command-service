---
name: generate
description: Run all //go:generate directives across the project
---

Run `task generate` to execute all `//go:generate` directives (go-enum and other configured generators).

## When to use

- When a new enum type is defined using go-enum (`// ENUM(...)`).
- When an existing go-enum type is modified (values added, renamed, or removed).
- When any new `//go:generate` directive is added to the codebase.
- Before running `/format` and `/lint` in sessions that involve generated files, so that stale generated code does not cause lint errors.

## When NOT to use

- When no `//go:generate` directives were added or changed.
- When creating plain domain aggregates, services, handlers, or repositories that do not use code generation.
- When creating a new migration — migrations are SQL, not generated Go code.
