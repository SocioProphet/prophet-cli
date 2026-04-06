# Agent Hybrid Model

## Surfaces
1. Deterministic CLI (canonical)
2. Agent assist (`ask`, `plan`)
3. Agent execute (`agent run`, approval-gated)

## Wrapper parity rules
- Overlay commands must resolve to deterministic tools.
- No semantic drift from deterministic core.
- Do not auto-load project MCP servers in untrusted workspaces.
