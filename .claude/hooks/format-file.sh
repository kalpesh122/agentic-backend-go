#!/usr/bin/env bash
# PostToolUse hook (Edit|Write): format the file that was just touched.
# Delegates to `just fmt-file <path>` so the kit stays stack-agnostic.
# Never blocks: formatting failures are reported as context, not errors.
set -u
INPUT=$(cat)
FILE=$(printf '%s' "$INPUT" | jq -r '.tool_input.file_path // empty')
[ -n "$FILE" ] && [ -f "$FILE" ] || exit 0
ROOT="${CLAUDE_PROJECT_DIR:-$(pwd)}"
cd "$ROOT" || exit 0
command -v just >/dev/null 2>&1 || exit 0
just --summary 2>/dev/null | tr ' ' '\n' | grep -qx 'fmt-file' || exit 0
if OUT=$(just fmt-file "$FILE" 2>&1); then
  exit 0
fi
jq -n --arg msg "Formatter reported a problem for $FILE: $(printf '%s' "$OUT" | tail -n 5)" \
  '{systemMessage: $msg}'
exit 0
