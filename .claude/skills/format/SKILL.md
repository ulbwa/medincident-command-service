---
name: format
description: Format all Go source files using gofumpt and gci
---

Run `task format` to format all Go source files.

## When to use

Use this skill **automatically** at the end of every session in which one or more `.go` files were created or modified — before presenting the result to the user for review.

Do **not** run it mid-task between intermediate edits; run it once, as the final step before handing off.

## When NOT to use

- When only non-Go files were changed (markdown, yaml, sql, etc.).
- When no files were modified at all.
