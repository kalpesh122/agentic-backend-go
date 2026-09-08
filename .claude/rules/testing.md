---
paths:
  - "**/*.test.*"
  - "**/*_test.*"
  - "**/tests/**"
  - "**/test/**"
  - "**/__tests__/**"
---

# Testing rules

- Tests are specifications: name them as sentences describing behaviour.
- No sleeps, no real network, no shared mutable state between tests. Use the stack's in-process DB or containers.
- Each test arranges its own data; do not depend on ordering.
- Assert on outcomes (response body, DB row, emitted event), not on implementation calls.
- A flaky test is a bug: fix it or delete it with a reason. Never retry to green.
- Coverage is a signal, not a target. New behaviour needs at least one test on the happy path and one on the failure path.
