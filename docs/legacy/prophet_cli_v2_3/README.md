# Legacy Prophet CLI v2.3 reference

This directory preserves the unpublished `prophet_cli_v2_3_final` attempt as **reference lineage**, not as canonical production code.

## Why this is here

The older CLI has several good ideas we want to preserve while rebuilding the façade and engine split:

- a minimal Cobra root and `cmd/prophet/main.go` entrypoint
- an early A2A workflow that already models author/test/review/revise/merge phases
- a transport framing experiment with structured headers, socket IPC, and authenticated encryption
- local carrier emission with signed JSON records written to disk

## What we keep

- command/workflow shape
- A2A lifecycle phases
- MCP config loading pattern
- the instinct that dry-run should still emit machine-readable carrier-like records

## What we replace

- ephemeral signing keys generated at process start
- the legacy TritRPC crypto/frame details
- zero-nonce / prototype transport assumptions
- the too-small command taxonomy (`prophet` + `a2a run` only)
- hardcoded demo-only tool IDs and method names

## Integration rule

Nothing in this directory should be compiled directly into the new façade or engine. These files are imported as a legacy reference set. Codex should port ideas out of them intentionally into:

- `internal/cmd/`
- `internal/a2a/`
- `internal/config/`
- `internal/runtime/`
- `internal/transport/`
- `internal/receipts/`

## Canonical direction

The new design remains:

- `prophet-cli` = façade and interaction surface
- `sourceos-sdk/cmd/sourceos-bootstrap` = bootstrap engine source home
- deterministic CLI verbs are canonical
- agent flows are overlays on deterministic tools

See also:
- `.handoff/codex-handoff/08-repo-patch-plans/prophet-cli-existing.md`
- `.handoff/codex-handoff/06-frozen-specs/CLI_Ergonomics_and_Agent_Hybrid.v1.md`
