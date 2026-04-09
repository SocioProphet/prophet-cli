# prophet-cli

[![CI](https://github.com/SocioProphet/prophet-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/SocioProphet/prophet-cli/actions/workflows/ci.yml)
[![Go 1.22](https://img.shields.io/badge/Go-1.22-blue.svg)](https://go.dev/doc/go1.22)

`prophet` is a command-line façade for the Prophet / SourceOS platform. It exposes a deterministic set of verbs for bootstrap operations, agent-to-agent workflow orchestration, and agent-assist scaffolds — while keeping all sensitive bootstrap business logic delegated to the `sourceos-bootstrap` engine.

## Why prophet-cli?

- **Single entry point** — one binary, consistent global flags, predictable output formats.
- **Safe by design** — no bootstrap secrets, no hidden side-effects, no auto-enrollment. Wrapper commands are transparent about what they delegate.
- **Hybrid model** — deterministic CLI verbs for automation, plus read-only agent-assist and approval-gated agent-execute layers for AI-assisted workflows.
- **Discoverable** — every command surfaces its delegation target and scaffold status in its output so you always know what is wired up vs. what is coming.

## Installation

### From source

```bash
git clone https://github.com/SocioProphet/prophet-cli.git
cd prophet-cli
go install ./cmd/prophet
```

### Requirements

- Go 1.22 or later
- `sourceos-bootstrap` engine installed and on `$PATH` (required for `prophet bootstrap` subcommands)

## Quick start

```bash
# Check that host prerequisites are met
prophet bootstrap doctor

# Prepare an authenticated bootstrap session
prophet bootstrap login

# Submit or describe a bootstrap build request
prophet bootstrap build

# Fetch release artifacts
prophet bootstrap fetch

# Show bootstrap engine information
prophet bootstrap info

# Run an A2A workflow in dry-run mode
prophet a2a run --repo owner/repo --ticket TICKET-123

# Run an A2A workflow live
prophet a2a run --repo owner/repo --ticket TICKET-123 --live
```

## Commands

See [docs/COMMANDS.md](docs/COMMANDS.md) for the full command reference with flags, arguments, and examples.

### Top-level commands

| Command | Description |
|---|---|
| `prophet bootstrap` | Façade over the SourceOS bootstrap engine |
| `prophet a2a` | Agent-to-agent workflow orchestration |
| `prophet ask` | Agent assist: explain or inspect without mutating state *(scaffold)* |
| `prophet plan` | Agent assist: generate a plan over deterministic tools *(scaffold)* |
| `prophet agent` | Agent execute façade, approval-gated *(scaffold)* |
| `prophet mcp` | MCP boundary façade *(scaffold)* |

### Global flags

| Flag | Short | Default | Description |
|---|---|---|---|
| `--profile` | | `""` | Named profile to use |
| `--space` | | `""` | Execution space |
| `--output` | `-o` | `text` | Output format: `text`, `json`, `yaml`, `table`, `tsv`, `none` |
| `--query` | | `""` | Query expression applied to structured output |
| `--quiet` | `-q` | `false` | Suppress non-essential output |
| `--debug` | | `false` | Enable debug output |
| `--no-pager` | | `false` | Disable pager for long output |

## Architecture

`prophet-cli` follows a three-layer hybrid model:

1. **Deterministic CLI** — canonical, documented, testable commands that resolve to concrete tools. These are the only commands that may be used in automation.
2. **Agent assist** — read-only planning and explanation layered over deterministic tools (`ask`, `plan`).
3. **Agent execute** — approval-gated execution that still resolves to deterministic tools (`agent`, `mcp`).

Bootstrap business logic lives entirely in `sourceos-bootstrap`. This repo is an intentionally thin façade that delegates to that engine and owns only the command surface. See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for a deeper explanation.

## Documentation

| Document | Description |
|---|---|
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | System architecture and design decisions |
| [docs/COMMANDS.md](docs/COMMANDS.md) | Full command reference |
| [docs/AGENT_HYBRID_MODEL.md](docs/AGENT_HYBRID_MODEL.md) | Three-layer hybrid CLI/agent model |
| [docs/BOOTSTRAP_DELEGATION.md](docs/BOOTSTRAP_DELEGATION.md) | Bootstrap delegation rules and boundaries |
| [docs/CLI_SURFACE_POLICY.md](docs/CLI_SURFACE_POLICY.md) | Canonical verb taxonomy and surface policy |
| [CONTRIBUTING.md](CONTRIBUTING.md) | Contribution guidelines |
| [SECURITY.md](SECURITY.md) | Security policy and vulnerability reporting |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines.

## Security

See [SECURITY.md](SECURITY.md) for the security policy and how to report vulnerabilities.
