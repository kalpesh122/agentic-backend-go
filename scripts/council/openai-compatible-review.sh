#!/usr/bin/env bash
# openai-compatible-review.sh — one council reviewer over any OpenAI-compatible chat API (DeepSeek, OpenAI,
# Groq, Mistral, a local vLLM/Ollama, …). Reads the prompt (rubric + schema + diff) on stdin, prints a JSON
# array of findings tagged with `reviewer`.
#
#   REVIEW_BASE_URL=https://api.deepseek.com REVIEW_API_KEY=$DEEPSEEK_API_KEY REVIEW_MODEL=deepseek-v4-flash \
#     scripts/council/openai-compatible-review.sh deepseek < prompt.md > findings.json
#
# Exit 0 with `[]` on any API problem (the council degrades, it never blocks), non-zero only on misuse.
set -uo pipefail
REVIEWER="${1:?reviewer name, e.g. deepseek}"
: "${REVIEW_BASE_URL:?REVIEW_BASE_URL is required}"; : "${REVIEW_API_KEY:?REVIEW_API_KEY is required}"; : "${REVIEW_MODEL:?REVIEW_MODEL is required}"
PROMPT=$(cat)
[ -n "$PROMPT" ] || { echo "no prompt on stdin" >&2; exit 2; }

BODY=$(jq -n --arg model "$REVIEW_MODEL" --arg prompt "$PROMPT" '{
  model: $model,
  temperature: 0,
  max_tokens: 4000,
  response_format: { type: "json_object" },
  messages: [
    { role: "system", content: "You are a senior code reviewer. Respond with a single JSON object {\"findings\": [...]} that matches the schema in the user message. No prose." },
    { role: "user", content: $prompt }
  ]
}')
RESP=$(curl -sS --max-time "${REVIEW_TIMEOUT:-180}" "$REVIEW_BASE_URL/chat/completions" \
  -H "Authorization: Bearer $REVIEW_API_KEY" -H 'content-type: application/json' -d "$BODY" 2>/dev/null) || { echo "warn: $REVIEWER request failed" >&2; echo '[]'; exit 0; }
if printf '%s' "$RESP" | jq -e '.error' >/dev/null 2>&1; then
  echo "warn: $REVIEWER API error: $(printf '%s' "$RESP" | jq -r '.error.message // .error')" >&2; echo '[]'; exit 0
fi
printf '%s' "$RESP" | jq -c --arg r "$REVIEWER" '
  (.choices[0].message.content // "") | sub("^```json\\s*";"") | sub("```\\s*$";"") | (try fromjson catch {findings: []})
  | [ (.findings // [])[] | select(type=="object") | . + {reviewer: $r} ]' 2>/dev/null || echo '[]'
