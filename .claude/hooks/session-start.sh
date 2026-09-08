#!/usr/bin/env bash
# SessionStart hook: cheap orientation for the agent (branch, dirty files, recipes).
set -u
ROOT="${CLAUDE_PROJECT_DIR:-$(pwd)}"
cd "$ROOT" || exit 0
BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "no-git")
DIRTY=$(git status --porcelain 2>/dev/null | wc -l | tr -d ' ')
RECIPES=$(command -v just >/dev/null 2>&1 && just --summary 2>/dev/null || echo "just not installed")
SPECS=$(ls -d specs/[0-9]*/ 2>/dev/null | tail -n 3 | tr '\n' ' ')
CTX="Repo orientation: branch=$BRANCH, dirty files=$DIRTY. just recipes: $RECIPES. Recent specs: ${SPECS:-none}. Read AGENTS.md before making changes; 'just check' is the quality gate."
jq -n --arg c "$CTX" '{hookSpecificOutput:{hookEventName:"SessionStart",additionalContext:$c}}'
exit 0
