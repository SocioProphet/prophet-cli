# Contributing

This repository is a façade-only CLI in phase 1.

Please:
- keep changes small and reviewable
- prefer deterministic command semantics
- keep bootstrap logic delegated to `sourceos-bootstrap`
- preserve explicit docs for boundaries and non-goals

Do not:
- invent boot or enrollment semantics here
- add hidden side effects to wrapper commands
- auto-enable project MCP servers in untrusted workspaces
