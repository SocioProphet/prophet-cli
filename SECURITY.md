# Security Policy

## Supported versions

| Version | Supported |
|---|---|
| `main` branch | ✅ Active development |
| Phase 1 scaffolds | ✅ Maintained |

## Reporting a vulnerability

If you discover a security vulnerability in `prophet-cli`, please **do not** open a public GitHub issue. Instead:

1. Email the maintainers at the address listed in the repository's GitHub profile, or
2. Use [GitHub private vulnerability reporting](https://docs.github.com/en/code-security/security-advisories/guidance-on-reporting-and-writing/privately-reporting-a-security-vulnerability) if enabled for this repository.

Please include:
- A clear description of the vulnerability and its potential impact
- Steps to reproduce or a minimal proof-of-concept
- Any suggested mitigations

You can expect an acknowledgement within 72 hours and a resolution timeline within 14 days for critical issues.

## What this repository must not contain

- Enrollment tokens
- Private signing keys
- Token-door secrets
- Hidden bootstrap business logic
- Auto-enrollment or auto-trust logic

## Phase 1 rules

- Keep sensitive semantics in their owning repositories; expose only façade wrappers here.
- Do not auto-enable project MCP servers in untrusted workspaces.
- Do not merge code that embeds bootstrap business logic into this repo.
- Wrapper commands must be transparent about delegation — every scaffold response includes a `delegated_to` field so callers know where logic actually runs.

## Dependency policy

Dependencies are kept minimal (currently only `github.com/spf13/cobra` and its transitive requirements). New dependencies require explicit justification in the pull request.
