#!/usr/bin/env bash
# Stop hook: if tracked source changed since the last green `just check`, run it now.
# Exit 2 blocks the stop and hands the failure output back to the agent.
# Disable for a session with AGENTIC_SKIP_VERIFY=1 (e.g. in settings.local.json env).
set -u
INPUT=$(cat)
[ "${AGENTIC_SKIP_VERIFY:-0}" = "1" ] && exit 0
# Never re-trigger while the agent is already responding to a Stop-hook failure.
printf '%s' "$INPUT" | jq -e '.stop_hook_active == true' >/dev/null 2>&1 && exit 0

ROOT="${CLAUDE_PROJECT_DIR:-$(pwd)}"
cd "$ROOT" || exit 0
command -v just >/dev/null 2>&1 || exit 0
git rev-parse --is-inside-work-tree >/dev/null 2>&1 || exit 0
just --summary 2>/dev/null | tr ' ' '\n' | grep -qx 'check' || exit 0

STAMP="$ROOT/.claude/.last-check"
# Fingerprint = HEAD + working-tree diff + untracked files. Cheap and deterministic.
FP=$( { git rev-parse HEAD 2>/dev/null; git diff HEAD --stat 2>/dev/null; git ls-files --others --exclude-standard 2>/dev/null; } | shasum | cut -d' ' -f1)
if [ -f "$STAMP" ] && [ "$(cat "$STAMP")" = "$FP" ]; then
  exit 0
fi
# Nothing changed at all? Skip.
if git diff --quiet HEAD 2>/dev/null && [ -z "$(git ls-files --others --exclude-standard)" ]; then
  exit 0
fi

if OUT=$(just check 2>&1); then
  printf '%s' "$FP" > "$STAMP"
  exit 0
fi
{
  echo "just check FAILED. Fix the problems below before finishing (do not disable tests)."
  printf '%s\n' "$OUT" | tail -n 60
} >&2
exit 2
