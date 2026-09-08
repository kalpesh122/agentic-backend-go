---
name: git-hygiene
description: Branching, committing, and PR conventions. Use when committing, pushing, or opening a pull request.
---

# Git hygiene

- **Branches**: `feat/<slug>`, `fix/<slug>`, `chore/<slug>`, `docs/<slug>`. Never commit directly to `main`.
- **Commits**: Conventional Commits. `type(scope): imperative summary ≤72 chars`, blank line, why-not-what body. Types: `feat fix refactor test docs chore perf build ci`.
- **Size**: one logical change per commit; a PR should be reviewable in ten minutes (~≤400 changed lines). Split otherwise.
- **Never**: force-push shared branches, rewrite `main`, commit secrets, commit generated files that the build produces, edit lockfiles by hand.
- **PR description**: use `.github/PULL_REQUEST_TEMPLATE.md` — link the spec, list the acceptance criteria met, paste `just check` output, justify any new dependency in one line.
- **Before push**: `just check`, then `just council` if the change is non-trivial.
- Merge strategy: squash unless the commit history is itself the documentation.
