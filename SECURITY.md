# Security Policy

This repository is a façade layer and should not hold long-lived secrets, token-door material, enrollment tokens, or private keys.

Phase-1 rules:
- keep sensitive semantics in the owning repos
- do not auto-enable project MCP servers in untrusted workspaces
- do not merge code that embeds bootstrap business logic into this repo
