# Commands

## Deterministic façade commands

### Bootstrap
- `prophet bootstrap doctor`
- `prophet bootstrap login`
- `prophet bootstrap build`
- `prophet bootstrap fetch`
- `prophet bootstrap write`
- `prophet bootstrap info`
- `prophet bootstrap validate <kind> <path>`
- `prophet bootstrap verify <kind> <path>`

### Workflow
- `prophet a2a run`

### Diagnostics
- `prophet doctor` — probe delegate engine readiness.
- `prophet status` — façade boundary legibility: enumerates every top-level surface
  as `real` / `delegating` / `scaffold`, names each delegate engine, and probes which
  engines are installed on PATH. Read-only; the runtime twin of
  [IMPLEMENTATION_STATUS.md](IMPLEMENTATION_STATUS.md).

## Delegated-action receipts
Pass `--receipt <path>` on any delegating command to write a small, command-scoped
JSON receipt (schema `prophet-cli/receipt/v0`) recording command, delegate, status,
args, timing, and any error. If `<path>` is a directory it receives a timestamped
file; otherwise it is treated as a file path. Receipts are **façade-local** convenience
artifacts — aligned conceptually with the estate ProofArtifact idea but **not** that
spine and **not** a runtime ledger. A receipt-write failure never masks the delegated
result.

## Hybrid overlay placeholders
- `prophet ask`
- `prophet plan`
- `prophet agent`
- `prophet mcp`

## Notes
These commands are scaffolds in phase 1. The deterministic verbs and delegation boundaries are more important than implementing broad business logic too early.
