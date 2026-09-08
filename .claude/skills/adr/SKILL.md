---
name: adr
description: Record an architecture decision in docs/adr/ using the MADR-style template. Use whenever a choice is made that a future engineer would ask "why?" about (library, pattern, data model, protocol, trade-off).
argument-hint: [short-title]
---

# Architecture Decision Record

1. Next number: `ls docs/adr | tail -1`.
2. Copy `docs/adr/0000-template.md` to `docs/adr/NNNN-<kebab-title>.md`.
3. Fill every section: Context (forces, constraints), Decision (one paragraph, active voice), Alternatives considered (each with why not), Consequences (good, bad, neutral), Status (`proposed` → `accepted` once merged).
4. Link the ADR from the spec or PR that motivated it. If it supersedes another ADR, mark that one `superseded by NNNN`.

Keep it under one screen. An ADR is a decision, not a design doc.
