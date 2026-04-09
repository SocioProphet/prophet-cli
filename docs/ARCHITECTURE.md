# Architecture

This document describes the design of `prophet-cli`, its relationship to upstream systems, and the decisions that shape the codebase.

## Overview

```
┌──────────────────────────────────────────────────────────────┐
│                        prophet-cli                           │
│                                                              │
│  ┌─────────────────┐  ┌──────────────┐  ┌────────────────┐  │
│  │ Deterministic   │  │ Agent assist │  │ Agent execute  │  │
│  │ CLI             │  │ (ask, plan)  │  │ (agent, mcp)   │  │
│  │ bootstrap, a2a  │  │ read-only    │  │ approval-gated │  │
│  └────────┬────────┘  └──────┬───────┘  └───────┬────────┘  │
│           │                  │                  │            │
└───────────┼──────────────────┼──────────────────┼────────────┘
            │                  └──────────────────┘
            │                  resolves to deterministic tools
            ▼
┌───────────────────────┐
│  sourceos-bootstrap   │
│  (external engine)    │
└───────────────────────┘
```

## Repository layout

```
prophet-cli/
├── cmd/
│   └── prophet/
│       └── main.go          # Binary entry point
├── internal/
│   ├── cmd/
│   │   ├── root.go          # Root command, global flags, emit helper
│   │   ├── bootstrap.go     # prophet bootstrap subcommands
│   │   ├── workflow.go      # prophet a2a subcommand
│   │   └── root_test.go     # Top-level command surface tests
│   └── a2a/
│       ├── workflow.go      # A2A workflow model and Default factory
│       └── workflow_test.go # Workflow model tests
├── docs/
│   ├── ARCHITECTURE.md      # This document
│   ├── COMMANDS.md          # Full command reference
│   ├── AGENT_HYBRID_MODEL.md
│   ├── BOOTSTRAP_DELEGATION.md
│   └── CLI_SURFACE_POLICY.md
├── CONTRIBUTING.md
├── README.md
└── SECURITY.md
```

## Design decisions

### Façade-only in phase 1

`prophet-cli` does not implement bootstrap business logic. It owns the command surface (verbs, flags, output format) and delegates to `sourceos-bootstrap` for all bootstrap operations. This separation means:

- The CLI can be versioned and tested independently of the bootstrap engine.
- The command surface is stable even while backend logic is evolving.
- Sensitive logic (enrollment, signing, token management) never enters this repo.

### Structured output everywhere

Every command produces output through the shared `emit` helper in `internal/cmd/root.go`. This ensures that `--output json` always works, regardless of the command's implementation status. Scaffold commands emit a JSON stub with `status: "scaffold"` and `delegated_to` fields so callers can distinguish wired commands from placeholders.

### Scaffold-first development

Commands are added as scaffolds before they are fully wired. This approach:

- Locks in the command surface early so documentation and tests can be written.
- Makes the delegation target explicit in the output.
- Prevents the command surface from drifting as the implementation catches up.

### Global flags

All commands inherit a consistent set of global flags (`--profile`, `--space`, `--output`, `--query`, `--quiet`, `--debug`, `--no-pager`). These are defined once on the root command and available to every subcommand.

### A2A workflow model

The `internal/a2a` package defines the structured workflow model used by `prophet a2a run`. A workflow has a repository, a ticket identifier, a live/dry-run flag, and a fixed sequence of phases:

1. `propose` — author proposes a change
2. `test` — author validates the candidate change
3. `review` — reviewer evaluates and may block or request revision
4. `revise` — author revises after review
5. `merge` — arbiter or maintainer merges the approved change
6. `done` — workflow is complete

## Testing

Tests live next to the code they test (`_test.go` files in the same package). Run the full suite with:

```bash
go test ./...
go vet ./...
```

CI runs these checks on every push and pull request via `.github/workflows/ci.yml`.
