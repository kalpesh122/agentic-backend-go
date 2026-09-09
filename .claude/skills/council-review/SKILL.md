---
name: council-review
description: Run or interpret the multi-model code review (Claude + Codex + Gemini reviewers, Claude synthesizer). Use before pushing a branch or when a PR has a council comment to act on.
argument-hint: [base-ref]
---

# Multi-model council review

Locally: `just council` (wraps `scripts/council-review.sh [base-ref]`). Needs whichever of `claude`, `codex`, `gemini` CLIs are installed and the matching API keys in `.env`; DeepSeek joins with just `DEEPSEEK_API_KEY` (called over its OpenAI-compatible API, no CLI). Missing reviewers are skipped, not fatal.

In CI: `.github/workflows/ai-council-review.yml` runs the same rubric with three read-only reviewer jobs and one synthesizer that posts a single sticky PR comment.

## How findings are weighted

- Findings are clustered when they name the same file, within 5 lines, same root cause.
- Flagged by 2+ models → severity raised one level. Flagged by exactly one → lowered one level.
- Anything below **medium** after weighting is dropped. Anything not on a changed line is dropped.

## Acting on a council comment

1. Fix every `critical`/`high` finding, or reply on the PR with a precise reason it is a false positive.
2. `medium`: fix or open a follow-up issue; say which.
3. Re-run `just council` before pushing the fix commit.

Do not argue with the council in code comments; argue in the PR thread.
