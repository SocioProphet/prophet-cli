# Security Policy

This repository must not contain:
- enrollment tokens
- private signing keys
- token-door secrets
- hidden bootstrap business logic

Phase 1 rule: keep sensitive semantics in their owning repositories and expose only façade wrappers here.
