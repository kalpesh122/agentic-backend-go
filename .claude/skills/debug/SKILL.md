---
name: debug
description: Systematic debugging for any bug, failing test, flaky CI, or unexpected behaviour. Use before proposing a fix.
---

# Systematic debugging

Do not guess-and-patch. Follow the ladder; stop as soon as the root cause is proven.

1. **Reproduce.** Get a deterministic reproduction: a failing test, a curl, a script. If you cannot reproduce it, you cannot claim to have fixed it.
2. **Read the actual error.** Full stack trace, exact message, exit code. Search the codebase for the message.
3. **Bisect the surface.** Narrow to the smallest input / commit / config that still fails. `git bisect` when the regression is time-based; binary-search the input when it is data-based.
4. **Form one hypothesis.** State it in a sentence with a prediction ("if X, then logging Y will show Z"). Test that prediction. Discard and repeat if wrong.
5. **Fix at the root.** Not at the symptom. If the fix is in a dependency, pin/patch and record an ADR.
6. **Regression test.** Turn the reproduction from step 1 into a permanent test.
7. **Verify.** `just check` green. Paste the evidence.

## Anti-patterns

- Changing three things at once.
- Adding retries, sleeps, or `try/catch` swallowing to hide the failure.
- "It works on my machine" without a reproduction in CI's environment.
- Declaring victory after the symptom disappears once.

## Handy commands

`just test -- <pattern>` (single test), `git log -p -- <file>`, `git bisect start <bad> <good>`, `just dev` with the stack's debug env var (see `.claude/rules/`).
