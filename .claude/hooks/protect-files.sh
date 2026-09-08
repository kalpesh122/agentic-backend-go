#!/usr/bin/env bash
# PreToolUse hook (Edit|Write|MultiEdit|NotebookEdit): deny edits to files agents must not touch by hand.
set -u
INPUT=$(cat)
FILE=$(printf '%s' "$INPUT" | jq -r '.tool_input.file_path // .tool_input.notebook_path // empty')
[ -n "$FILE" ] || exit 0
BASE=$(basename "$FILE")

deny() {
  jq -n --arg r "$1" '{hookSpecificOutput:{hookEventName:"PreToolUse",permissionDecision:"deny",permissionDecisionReason:$r}}'
  exit 0
}

case "$BASE" in
  .env|.env.local|.env.development|.env.production|.env.test|.env.staging)
    deny "Blocked: $BASE holds secrets. Edit .env.example instead and tell the user which value to set." ;;
  pnpm-lock.yaml|package-lock.json|yarn.lock|bun.lockb|bun.lock|uv.lock|poetry.lock|Cargo.lock|go.sum|Gemfile.lock|composer.lock)
    deny "Blocked: $BASE is generated. Run the package manager instead of editing the lockfile." ;;
esac

case "$FILE" in
  */.git/*) deny "Blocked: never edit .git internals directly." ;;
  */node_modules/*|*/.venv/*|*/vendor/*) deny "Blocked: $FILE is a dependency, not project code." ;;
esac

exit 0
