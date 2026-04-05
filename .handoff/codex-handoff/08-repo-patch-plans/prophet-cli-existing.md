# Patch Plan: prophet-cli-existing — Revision 2.1

## Objective
Patch the existing ProphetCLI/ProfitCLI repository to act as the user-facing façade and interactive shell surface.

## Files to create or update
- README.md
- docs/BOOTSTRAP_DELEGATION.md
- docs/CLI_SURFACE_POLICY.md
- docs/AGENT_HYBRID_MODEL.md
- docs/COMMANDS.md
- src/ or cmd/ wrapper entrypoints for:
  - `prophet bootstrap ...`
  - `prophet ask ...`
  - `prophet plan ...`
  - `prophet agent run ...`
  - `prophet mcp serve`
  - `prophet mcp doctor`

## Rules
- Do not embed bootstrap business logic in this repo.
- All bootstrap actions delegate to the installed `sourceos-bootstrap` binary or SDK runtime.
- The deterministic command surface remains canonical.
- Agent surfaces must resolve to deterministic tools and approvals.
- Do not auto-install plugins or auto-enable project MCP in untrusted workspaces.

## Validation
- README explains façade/engine split
- bootstrap wrapper docs exist
- command docs distinguish deterministic vs agent overlay surfaces
- no file claims this repo owns BootReleaseSet, TokenDoor, or ConfigSource semantics
