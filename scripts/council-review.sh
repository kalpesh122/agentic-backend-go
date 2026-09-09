#!/usr/bin/env bash
# council-review.sh — local multi-model code review.
#
# Runs whichever of Claude Code, OpenAI Codex, and Gemini CLI are installed over the same diff and
# rubric, then has Claude synthesize one deduplicated, agreement-weighted findings list.
#
# Usage:   scripts/council-review.sh [base-ref]      # default: origin/main (falls back to main)
# Env:     ANTHROPIC_API_KEY, OPENAI_API_KEY, GEMINI_API_KEY, DEEPSEEK_API_KEY (read from .env if present)
#          COUNCIL_CLAUDE_MODEL, COUNCIL_CODEX_MODEL, COUNCIL_GEMINI_MODEL, COUNCIL_DEEPSEEK_MODEL to override
#          COUNCIL_SKIP=codex,gemini,deepseek to skip reviewers
#          DeepSeek needs no CLI: it is called over its OpenAI-compatible API (scripts/council/openai-compatible-review.sh)
# Exit:    0 = no critical findings, 1 = critical findings, 2 = nothing to review / setup problem
# Output:  .council/report.md (gitignored) and stdout
set -uo pipefail

ROOT=$(git rev-parse --show-toplevel 2>/dev/null) || { echo "Not a git repo" >&2; exit 2; }
cd "$ROOT" || exit 2
[ -f .env ] && set -a && . ./.env && set +a

BASE="${1:-origin/main}"
git rev-parse --verify -q "$BASE" >/dev/null || BASE=main
git rev-parse --verify -q "$BASE" >/dev/null || { echo "Base ref not found: $BASE" >&2; exit 2; }

CLAUDE_MODEL="${COUNCIL_CLAUDE_MODEL:-claude-opus-5}"
CODEX_MODEL="${COUNCIL_CODEX_MODEL:-gpt-5.6-terra}"
GEMINI_MODEL="${COUNCIL_GEMINI_MODEL:-gemini-3.1-pro-preview}"
DEEPSEEK_MODEL="${COUNCIL_DEEPSEEK_MODEL:-deepseek-v4-flash}"
SKIP=",${COUNCIL_SKIP:-},"

OUT="$ROOT/.council"; rm -rf "$OUT"; mkdir -p "$OUT"
git diff "$BASE"...HEAD > "$OUT/diff.patch"
[ -s "$OUT/diff.patch" ] || { echo "No changes vs $BASE"; exit 2; }

RUBRIC=".github/review-rubric.md"; [ -f "$RUBRIC" ] || RUBRIC=/dev/null
SCHEMA=".github/review-schema.json"
[ -f "$SCHEMA" ] || { echo "Missing $SCHEMA" >&2; exit 2; }

{
  echo "You are a senior reviewer. Review ONLY the unified diff below."
  echo "Report a finding only if you are confident and it is on a changed line. No style nits."
  echo "Return a JSON object {\"findings\": [...]} matching this JSON schema exactly:"
  cat "$SCHEMA"
  echo; echo "## Rubric"; cat "$RUBRIC"
  echo; echo "## Diff"; echo '```diff'; cat "$OUT/diff.patch"; echo '```'
} > "$OUT/prompt.md"

have() { command -v "$1" >/dev/null 2>&1; }
skipped() { case "$SKIP" in *",$1,"*) return 0;; esac; return 1; }
PIDS=(); NAMES=()

if have claude && [ -n "${ANTHROPIC_API_KEY:-}" ] && ! skipped claude; then
  ( claude --bare -p "$(cat "$OUT/prompt.md")" --model "$CLAUDE_MODEL" \
      --output-format json --json-schema "$(cat "$SCHEMA")" \
      --permission-mode dontAsk --permission-prompts none --max-turns 4 2>"$OUT/claude.err" \
    | jq -c '[(.structured_output.findings // [])[] + {reviewer:"claude"}]' > "$OUT/claude.json" ) &
  PIDS+=($!); NAMES+=(claude)
else echo "skip: claude (missing CLI, ANTHROPIC_API_KEY, or COUNCIL_SKIP)"; fi

if have codex && [ -n "${OPENAI_API_KEY:-}" ] && ! skipped codex; then
  ( codex exec --model "$CODEX_MODEL" --sandbox read-only \
      --output-schema "$SCHEMA" -o "$OUT/codex.raw" --skip-git-repo-check - < "$OUT/prompt.md" >"$OUT/codex.err" 2>&1 \
    && jq -c '[(.findings // [])[] + {reviewer:"codex"}]' "$OUT/codex.raw" > "$OUT/codex.json" ) &
  PIDS+=($!); NAMES+=(codex)
else echo "skip: codex (missing CLI, OPENAI_API_KEY, or COUNCIL_SKIP)"; fi

if have gemini && [ -n "${GEMINI_API_KEY:-}" ] && ! skipped gemini; then
  ( gemini -m "$GEMINI_MODEL" --output-format json \
      -p "$(cat "$OUT/prompt.md") Return ONLY the JSON object, no prose." 2>"$OUT/gemini.err" \
    | jq -c '[((.response // "") | sub("^```json\\s*";"") | sub("```\\s*$";"") | fromjson? // {findings:[]}).findings[] + {reviewer:"gemini"}]' \
    > "$OUT/gemini.json" ) &
  PIDS+=($!); NAMES+=(gemini)
else echo "skip: gemini (missing CLI, GEMINI_API_KEY, or COUNCIL_SKIP)"; fi

if [ -n "${DEEPSEEK_API_KEY:-}" ] && ! skipped deepseek; then
  ( REVIEW_BASE_URL="${DEEPSEEK_BASE_URL:-https://api.deepseek.com}" REVIEW_API_KEY="$DEEPSEEK_API_KEY" REVIEW_MODEL="$DEEPSEEK_MODEL" \
      "$ROOT/scripts/council/openai-compatible-review.sh" deepseek < "$OUT/prompt.md" > "$OUT/deepseek.json" 2>"$OUT/deepseek.err" ) &
  PIDS+=($!); NAMES+=(deepseek)
else echo "skip: deepseek (missing DEEPSEEK_API_KEY or COUNCIL_SKIP)"; fi

[ "${#PIDS[@]}" -gt 0 ] || { echo "No reviewers available. Install claude/codex/gemini, or set DEEPSEEK_API_KEY, in .env." >&2; exit 2; }
echo "reviewers: ${NAMES[*]} (base: $BASE)"
for p in "${PIDS[@]}"; do wait "$p" || true; done

for n in "${NAMES[@]}"; do
  [ -s "$OUT/$n.json" ] && jq -e 'type=="array"' "$OUT/$n.json" >/dev/null 2>&1 || { echo "warn: $n produced no parsable findings (see $OUT/$n.err)"; echo '[]' > "$OUT/$n.json"; }
done
jq -s 'add // []' "$OUT"/*.json > "$OUT/all.json"
COUNT=$(jq 'length' "$OUT/all.json")
echo "raw findings: $COUNT"
if [ "$COUNT" -eq 0 ]; then echo "No blocking findings." | tee "$OUT/report.md"; exit 0; fi

if have claude && [ -n "${ANTHROPIC_API_KEY:-}" ]; then
  claude --bare -p "$(cat .claude/agents/council-synthesizer.md 2>/dev/null | sed '1,/^---$/d' | sed '1,/^---$/d')

FINDINGS (JSON):
$(cat "$OUT/all.json")

DIFF:
$(cat "$OUT/diff.patch")" \
    --model "$CLAUDE_MODEL" --output-format text --max-turns 2 \
    --permission-mode dontAsk --permission-prompts none 2>"$OUT/synth.err" | tee "$OUT/report.md"
else
  # No synthesizer: print raw findings grouped by severity.
  jq -r 'sort_by(.severity) | .[] | "### \(.severity) — \(.file):\(.line)\n\(.title)\n\(.detail)\n(flagged by: \(.reviewer))\n"' "$OUT/all.json" | tee "$OUT/report.md"
fi

grep -Eq '^### critical' "$OUT/report.md" && { echo "Blocking: critical findings."; exit 1; }
exit 0
