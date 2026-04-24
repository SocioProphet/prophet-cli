# Contributing

This repository is a façade-only CLI in phase 1.

Please:
- keep changes small and reviewable
- prefer deterministic command semantics
- keep bootstrap logic delegated to `sourceos-bootstrap`
- preserve explicit docs for boundaries and non-goals
- improve command grammar
- add wrapper docs
- add compileable scaffolds
- improve tests and examples

Do not:
- embed bootstrap business logic here
- invent boot or enrollment semantics here
- add hidden side effects to wrapper commands
- auto-enable project MCP servers in untrusted workspaces
- bypass the frozen object set from the SourceOS handoff materials
