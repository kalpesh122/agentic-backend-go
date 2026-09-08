# Review rubric

Shared by human reviewers, the `reviewer` subagent, `just council`, and the CI council. Review **only changed lines** and what they directly affect. Report confident findings only. The formatter owns style; do not comment on it.

## Severity

| Level | Meaning |
|-------|---------|
| critical | Data loss, security hole, crash on a common path, silently wrong result, secret exposure |
| high | Wrong behaviour on a realistic path, missing auth/authz, unbounded resource use, broken migration |
| medium | Fragile logic, missing test for new behaviour, unclear or swallowed errors, N+1 on a hot path |
| low | Naming, minor duplication, docs, dead code |

## Categories (check in this order)

1. **correctness** — logic errors, off-by-one, wrong null handling, race conditions, wrong async/await, unhandled rejections, incorrect type narrowing.
2. **security** — injection, missing auth, secrets, SSRF, path traversal, unsafe deserialisation, permissive CORS, prompt-injection paths in LLM code.
3. **data** — migrations reversible and safe under load, constraints and indexes match access patterns, no destructive change without a plan, transactions where needed.
4. **performance** — N+1, unbounded queries, missing pagination, sync IO on hot paths, unnecessary allocation in loops, missing caching where the spec demands it.
5. **tests** — new behaviour without tests, tests that cannot fail, weakened assertions, skipped tests, flaky patterns (sleeps, ordering, shared state).
6. **maintainability** — duplication of an existing utility, function doing three things, misleading names, dead code, placeholders (TODO/FIXME/stub), undocumented magic numbers.

## Output

One entry per finding: `file`, `line`, `severity`, `category`, `title` (one-sentence claim), `detail` (concrete failure scenario: input → wrong outcome), `suggested_fix`. Schema: `.github/review-schema.json`.

Verdict: `approve` / `approve-with-nits` / `request-changes` plus the single most important fix. If nothing survives: "No blocking findings." and why the change is safe.
