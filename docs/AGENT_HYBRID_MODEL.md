# Agent hybrid model

`prophet-cli` follows a three-layer hybrid model that cleanly separates deterministic automation from AI-assisted workflows.

## Layers

### 1. Deterministic CLI

The deterministic CLI is the canonical, documented surface. These commands:

- Have stable, versioned semantics.
- Are safe to use in scripts and CI pipelines.
- Emit structured output (`--output json`) for downstream tooling.
- Are the **only** commands that may be used in automation without human review.

Examples: `prophet bootstrap doctor`, `prophet bootstrap validate`, `prophet a2a run`.

### 2. Agent assist

Read-mostly operations for planning and explanation. Agent-assist commands (`ask`, `plan`) operate over deterministic tools but do **not** mutate state. They are intended for:

- Explaining the current state of the system.
- Generating a proposed plan that a human or agent-execute layer can review before running.
- Answering questions about objects, configurations, or workflows.

Agent-assist commands must not introduce side effects.

### 3. Agent execute

Approval-gated execution that still resolves to deterministic tools. Agent-execute commands (`agent`, `mcp`) may trigger mutations, but only after explicit human or policy approval. They:

- Resolve to the same deterministic verbs as layer 1.
- Surface their intended actions before executing.
- Respect the same bootstrap delegation boundaries as the deterministic CLI.

## Guardrails

These rules apply to all three layers:

- No bootstrap business logic in this repo — delegate to `sourceos-bootstrap`.
- No silent plugin installation.
- No automatic project MCP loading in untrusted workspaces.
- No competing command language that bypasses deterministic verbs.
- All mutations must be traceable to a deterministic tool invocation.

## Phase 1 status

In phase 1, layers 2 and 3 (`ask`, `plan`, `agent`, `mcp`) are scaffolds. They emit structured stubs so the command surface and output shape are stable before the agent runtime is wired up.

See [CLI_SURFACE_POLICY.md](CLI_SURFACE_POLICY.md) for the canonical verb taxonomy.
