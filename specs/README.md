# specs/

Spec-driven development artifacts. One folder per feature: `NNN-<slug>/` with `spec.md`, `plan.md`, `tasks.md`.
Create them with the `brainstorm-spec` skill; execute `tasks.md` with the `tdd` skill; close with `verify-before-done`.

- `spec.md` — what and why; acceptance criteria are testable sentences.
- `plan.md` — ordered steps, files touched, test per step, rollback.
- `tasks.md` — the checklist and progress log.

Numbering starts at `001`. `000-template/` is the template; copy it.
