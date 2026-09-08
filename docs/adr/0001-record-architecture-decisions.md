# 0001. Record architecture decisions

Date: 2026-09-09
Status: accepted

## Context

Agents and humans both need to know why the codebase looks the way it does. Chat history and PR threads are not durable or discoverable. Decisions that are not written down get re-litigated or silently reversed.

## Decision

We will record every architecturally significant decision as a short ADR in `docs/adr/`, numbered sequentially, using `0000-template.md`. The `adr` skill guides the format. ADRs are immutable once accepted; a change is a new ADR that supersedes the old one.

## Alternatives considered

- **Wiki / Notion pages** — not versioned with the code, invisible to agents.
- **Comments in code** — too local; decisions span files.

## Consequences

- Good: agents can read the rationale before proposing changes; onboarding is faster.
- Bad: small friction per decision.
- Neutral: ADRs must be linked from specs and PRs to stay discoverable.
