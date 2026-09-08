---
name: tdd
description: Red-green-refactor discipline for any behaviour change. Use whenever you add or change code that does something (not for pure renames or docs).
---

# TDD

The loop, every time:

1. **Red.** Write the smallest test that fails for the right reason. Run it. Paste the failure. If it passes immediately, the test is wrong or the feature already exists.
2. **Green.** Write the least code that makes it pass. No speculative generality.
3. **Refactor.** With tests green, remove duplication, improve names, extract only what is now used twice.
4. **Repeat** for the next behaviour in `tasks.md`.

## Rules

- One behaviour per test. Name tests as sentences: `returns 404 when the todo does not exist`.
- Test through the public surface (HTTP handler, exported function), not private internals.
- Prefer real dependencies in tests (in-memory Postgres, temp dirs) over mocks. Mock only the network edge you do not own.
- A bug fix starts with a failing regression test that reproduces the bug.
- Never weaken an assertion, skip a test, or widen a timeout to get green. If a test is wrong, delete it with a commit message that says why.
- Commit after each green step: `test: …`, `feat: …`, `refactor: …`.

## Where tests live

Colocated: `thing.ts` ↔ `thing.test.ts` (or the stack's convention in `.claude/rules/testing.md`). Integration tests under `tests/` or `test/` if the stack separates them.

## Done means

`just test` green, `just check` green, new behaviour covered, no `.only`/`.skip`/`xit` left behind.
