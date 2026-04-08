# Security Policy

This repository must not contain:
- enrollment tokens
- private signing keys
- token-door secrets
- hidden bootstrap business logic

Phase 1 rules:
- keep sensitive semantics in their owning repositories and expose only façade wrappers here
- do not auto-enable project MCP servers in untrusted workspaces
- do not merge code that embeds bootstrap business logic into this repo
