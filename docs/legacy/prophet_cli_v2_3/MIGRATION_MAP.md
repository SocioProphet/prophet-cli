# Legacy v2.3 to current layout migration map

| Legacy file | Keep / replace / adapt | New target | Notes |
|---|---|---|---|
| `cmd/prophet/main.go` | keep shape | `cmd/prophet/main.go` | keep a small entrypoint |
| `internal/cmd/root.go` | adapt | `internal/cmd/root.go` | expand to canonical verbs, profiles, spaces, agent overlay |
| `internal/cmd/a2a.go` | adapt | `internal/cmd/a2a.go` + `internal/a2a/` | preserve workflow phases, remove hardcoded demo semantics |
| `internal/mcp/tritrpc.go` | replace | `internal/config/` + `internal/transport/` | preserve header/socket/config-loading ideas; replace crypto and frame semantics |
| `internal/util/carrier.go` | replace | `internal/receipts/` + `internal/runtime/` | preserve local machine-readable emission; replace ephemeral key model |
| `internal/util/crc16.go` | optional compat | `internal/transport/crc16_compat.go` | only keep if needed for test fixtures or compatibility |
| `README.md` | replace | `README.md` + `docs/` | façade-oriented README; legacy README preserved as reference |

## Porting order

1. import legacy files as reference only
2. preserve root/main command shape
3. port A2A lifecycle into `internal/a2a/`
4. port MCP config-loading idea into `internal/config/`
5. redesign carrier/receipt emitter against frozen objects
6. redesign transport against the current canonical runtime and policy model

## Explicit non-goals

- do not compile the legacy transport/emitter as production code
- do not keep the ephemeral signing-key behavior
- do not treat the legacy `TritRPC` helper as the final protocol implementation
