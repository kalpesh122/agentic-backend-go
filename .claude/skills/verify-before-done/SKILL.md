---
name: verify-before-done
description: Evidence before claims. Use before saying a task is complete, before committing, and before opening a PR.
---

# Verify before done

"Done" is a claim; claims need evidence. Run these and paste the relevant output (last lines are enough):

1. `just check` — must be green. If any step is skipped by the recipe, run it explicitly.
2. `git status --short` — every intended file is staged; nothing unintended is.
3. `git diff --cached --stat` — the change is the size you expected.
4. Grep the diff for placeholders: `TODO|FIXME|XXX|HACK|coming soon|not implemented|\.only\(|\.skip\(|xit\(`. Zero hits, or each one justified in the PR.
5. Re-read the acceptance criteria in `specs/NNN-*/spec.md`. Tick each one in `tasks.md` with a note on how it is proven.
6. If behaviour changed: docs/README/ADR updated in the same commit.

Then say "done" with the evidence. If any step fails, you are not done; say what is left.
