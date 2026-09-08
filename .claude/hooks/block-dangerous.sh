#!/usr/bin/env bash
# PreToolUse hook (Bash): deny destructive or exfiltrating shell commands.
# Exit 0 with a deny decision so the reason reaches the agent verbatim.
set -u
INPUT=$(cat)
CMD=$(printf '%s' "$INPUT" | jq -r '.tool_input.command // empty')
[ -n "$CMD" ] || exit 0

deny() {
  jq -n --arg r "$1" '{hookSpecificOutput:{hookEventName:"PreToolUse",permissionDecision:"deny",permissionDecisionReason:$r}}'
  exit 0
}

# Normalise whitespace for matching.
C=$(printf '%s' "$CMD" | tr -s '[:space:]' ' ')

# rm -rf outside of temp/cache/build dirs
if printf '%s' "$C" | grep -Eq '(^|[;&| ])rm +(-[a-zA-Z]*r[a-zA-Z]*f|-[a-zA-Z]*f[a-zA-Z]*r)[a-zA-Z]* '; then
  TARGETS=$(printf '%s' "$C" | sed -E 's/.*rm +-[a-zA-Z]+ +//')
  if ! printf '%s' "$TARGETS" | grep -Eq '^(/tmp/|/private/tmp/|\$TMPDIR|\./?(node_modules|dist|build|\.next|\.turbo|coverage|\.venv|target|bin|out)(/| |$))'; then
    deny "Blocked: 'rm -rf' outside temp/build directories. Delete files explicitly or ask the user."
  fi
fi

printf '%s' "$C" | grep -Eq 'git +push +.*(--force|-f)( |$)' && deny "Blocked: force push. Use --force-with-lease only after the user explicitly asks."
printf '%s' "$C" | grep -Eq 'git +reset +--hard' && deny "Blocked: git reset --hard discards work. Use git stash or ask the user."
printf '%s' "$C" | grep -Eq 'git +clean +-[a-zA-Z]*f' && deny "Blocked: git clean -f deletes untracked files. Ask the user."
printf '%s' "$C" | grep -Eq 'git +checkout +(--|\.) *$|git +restore +\. *$' && deny "Blocked: discarding all working-tree changes. Ask the user."
printf '%s' "$C" | grep -Eq 'git +branch +-D ' && deny "Blocked: force-deleting a branch. Ask the user."
printf '%s' "$C" | grep -Eq '(curl|wget) [^|]*\| *(sudo +)?(ba|z)?sh' && deny "Blocked: piping a download into a shell. Download, inspect, then run."
printf '%s' "$C" | grep -Eq 'chmod +(-R +)?777' && deny "Blocked: chmod 777."
printf '%s' "$C" | grep -Eiq '(drop +(table|database|schema)|truncate +table)' && deny "Blocked: destructive SQL. Write a migration and ask the user."
printf '%s' "$C" | grep -Eq '(^|[;&| ])(cat|less|more|head|tail|bat) +[^|;&]*\.env($|[^.a-zA-Z])' && deny "Blocked: reading .env. Use .env.example for the variable names."
printf '%s' "$C" | grep -Eq '>+ *\.env($| )' && deny "Blocked: writing .env directly. Tell the user which variable to set."
printf '%s' "$C" | grep -Eq '(^|[;&| ])sudo ' && deny "Blocked: sudo. Ask the user to run privileged commands."

exit 0
