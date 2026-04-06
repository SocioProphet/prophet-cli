# CLI Ergonomics and Agent Hybrid — Revision 2.1

## Purpose
Freeze the command-surface philosophy for phase 1 so Codex does not invent divergent semantics.

## Model
The system exposes three coordinated surfaces:
1. **Deterministic CLI** — canonical, documented, scriptable.
2. **Agent assist** — read-mostly planning/explanation over the same tools.
3. **Agent execute** — approval-gated execution of deterministic tools.

## Deterministic core
The canonical verbs are:
- read: `list`, `show`, `describe`, `explain`, `find`, `search`
- mutate: `create`, `update`, `delete`, `apply`, `wait`, `validate`, `verify`
- runtime: `status`, `logs`, `start`, `stop`, `restart`, `enable`, `disable`
- bootstrap/host: `doctor`, `login`, `build`, `fetch`, `write`, `info`, `audit`

Global flags:
- `--profile`
- `--space`
- `--output` / `-o`
- `--query`
- `--quiet`
- `--debug`
- `--no-pager`

## Hybrid overlay
Agent affordances may be added, but they must resolve to deterministic tools:
- `ask`
- `plan`
- `agent run`
- `skill list`
- `skill run`
- `session resume`
- `mcp serve`
- `mcp doctor`

## Repo boundary
- Existing ProphetCLI/ProfitCLI repo: façade, wrappers, docs, command delegation.
- `sourceos-sdk/cmd/sourceos-bootstrap`: bootstrap engine source home.
- `homebrew-sourceos`: distribution only.

## Workspace trust
Project hooks, MCP servers, local memory, and custom commands must not auto-load in untrusted workspaces.

## Non-goals in this revision
- final boot enrollment semantics
- final MCP server auto-discovery policy
- final subagent topology
