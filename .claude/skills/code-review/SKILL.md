---
name: code-review
description: Review a diff against the shared rubric in .github/review-rubric.md and report findings in the shared JSON/markdown shape. Use for self-review before opening a PR and when asked to review someone else's change.
argument-hint: [base-ref]
---

# Code review

Rubric: `.github/review-rubric.md` (read it first). Output shape: `.github/review-schema.json`.

## Procedure

1. `git diff ${1:-main}...HEAD` — review only changed lines and what they directly affect.
2. Read each changed file in full once; understand intent from the spec/PR description.
3. Walk the rubric categories in order: correctness → security → data/migrations → performance → tests → maintainability.
4. For each finding: file, line, severity (critical/high/medium/low), category, one-sentence claim, the concrete failure scenario, and a suggested fix.
5. Report only findings you are confident in and that sit on changed lines. No style nits (the formatter owns style).
6. End with a verdict: `approve`, `approve-with-nits`, or `request-changes`, and the single most important thing to fix.

## Calibration

- **critical**: data loss, security hole, crash on common path, silent wrong result.
- **high**: wrong behaviour on a realistic path, missing auth check, unbounded resource use.
- **medium**: fragile code, missing test for new behaviour, unclear error handling.
- **low**: naming, minor duplication, docs.

Praise is noise; skip it. If nothing survives, say "No blocking findings." and why you believe the change is safe.
