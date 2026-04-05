# Legacy Prophet v2.3 file manifest

The following unpublished files were reviewed and mapped into the new façade/engine design.

## Reviewed files

- `internal/cmd/a2a.go`
- `internal/cmd/root.go`
- `cmd/prophet/main.go`
- `internal/mcp/tritrpc.go`
- `internal/util/carrier.go`
- `internal/util/crc16.go`
- `README.md`

## What they taught us

### Keep
- small Cobra root and entrypoint structure
- A2A workflow phases: propose, test, review/block, revise, merge, done
- config-driven MCP server lookup
- machine-readable local carrier emission, even in dry-run mode

### Replace
- ephemeral signing keys generated at startup
- transport crypto and nonce handling from the prototype
- minimal command taxonomy (`prophet` + `a2a run` only)
- hardcoded demo-only tool IDs and method names

## New target layout

The old files should inform these targets:

- `cmd/prophet/main.go`
- `internal/cmd/root.go`
- `internal/cmd/a2a.go`
- `internal/a2a/`
- `internal/config/`
- `internal/runtime/`
- `internal/transport/`
- `internal/receipts/`

## Integration rule

Legacy files are **reference lineage**, not canonical production code. Port ideas intentionally and replace trust/protocol internals with the newer frozen design.
