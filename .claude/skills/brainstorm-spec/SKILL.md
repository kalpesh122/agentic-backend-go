---
name: brainstorm-spec
description: Turn a feature request into a spec, plan, and task list before any code is written. Use for every non-trivial change (new feature, new endpoint, schema change, refactor touching more than three files).
argument-hint: [feature-slug]
---

# Brainstorm → Spec → Plan → Tasks

Goal: leave `specs/NNN-<slug>/` with three files a different engineer (or agent) could execute without talking to you.

## 1. Classify

- **Bounded** (one file, existing flow, obvious change): skip the files, state the design in two sentences, get a nod, implement.
- **Architectural** (new capability, new data, new interface, cross-cutting): do all steps below.

## 2. Understand (read before asking)

- Read the code paths the change touches. Read the latest `specs/` and `docs/adr/` entries.
- Then ask the human **one question at a time**, only about things you cannot infer: purpose, constraints, success criteria, non-goals.

## 3. Propose 2–3 approaches

Each with: what it is, trade-offs, risk, effort. Lead with your recommendation and why. YAGNI: strip anything not needed for the stated success criteria.

## 4. Write the three files (copy `specs/000-template/`)

- `spec.md` — problem, goals, non-goals, acceptance criteria (testable sentences), data/API changes, risks, open questions (must be empty before implementation).
- `plan.md` — ordered steps, each small enough to finish in one sitting; files touched; test strategy per step; rollback.
- `tasks.md` — checklist derived from the plan. Each item: `- [ ] <verb> <object> — test: <how it's proven>`.

## 5. Self-review the spec

Placeholder scan (no TBD/TODO), internal consistency, scope (one plan, not three), ambiguity (each requirement has exactly one reading). Fix inline.

## 6. Hand off

Ask the human to read `spec.md`. Only after approval, start `tasks.md` top to bottom with the `tdd` skill. Update `tasks.md` as you go; it is the progress log.

Numbering: next free `NNN` in `specs/`. Slug: kebab-case, ≤4 words.
