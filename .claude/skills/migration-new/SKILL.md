---
name: migration-new
description: Create a new database migration file via dbmate
---

Run `task migration-new -- <name>` where `<name>` is a short, snake_case description of the schema change (e.g. `add_users_table`, `add_status_to_orders`).

After the file is created, fill in both `-- migrate:up` and `-- migrate:down` sections completely before committing.

## When to use

Use this skill **only** when implementing or extending the **repository layer** in a way that requires a new or modified database schema:

- Adding a new table for a new aggregate root.
- Adding a column to an existing table (e.g. a new field required by the domain).
- Adding an index required for a new query in `read_repository.go`.
- Dropping a table or column that is no longer used (write both up and down).

## When NOT to use

- When creating a **domain model or aggregate** (`internal/model/`) — domain code has no DB schema.
- When creating a **service** (`internal/service/`) — services orchestrate, they do not own tables.
- When creating a **handler** (`internal/handler/`) — transport layer, no schema impact.
- When adding a new **read query** to an existing `read_repository.go` that uses already-existing columns.
- When refactoring code without changing the stored data shape.

> **Important:** `task migration-up` and `task migration-down` are intentionally blocked for the AI agent and must be run manually by the developer after verifying `DATABASE_URL` points to the correct environment.
