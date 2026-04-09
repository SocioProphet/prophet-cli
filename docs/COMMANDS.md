# Commands

This document is the canonical reference for all `prophet` commands. Commands marked **scaffold** emit a structured JSON stub and indicate their `delegated_to` target; they are placeholders for the full implementation.

---

## Global flags

These flags apply to every command:

| Flag | Short | Default | Description |
|---|---|---|---|
| `--profile` | | `""` | Named profile to use |
| `--space` | | `""` | Execution space |
| `--output` | `-o` | `text` | Output format: `text`, `json`, `yaml`, `table`, `tsv`, `none` |
| `--query` | | `""` | Query expression applied to structured output |
| `--quiet` | `-q` | `false` | Suppress non-essential output |
| `--debug` | | `false` | Enable debug output |
| `--no-pager` | | `false` | Disable pager |

---

## Deterministic façade commands

### `prophet bootstrap`

Façade over the installed `sourceos-bootstrap` engine. All subcommands delegate to that engine; this CLI owns only the command surface and output formatting.

#### `prophet bootstrap doctor`

Diagnose host prerequisites. Checks that required tools and environment configuration are in place before a bootstrap session.

```bash
prophet bootstrap doctor
prophet bootstrap doctor --output json
```

#### `prophet bootstrap login`

Prepare an authenticated bootstrap session. Establishes credentials that subsequent bootstrap commands will use.

```bash
prophet bootstrap login
prophet bootstrap login --profile myprofile
```

#### `prophet bootstrap build`

Submit or describe a bootstrap build request. Delegates build orchestration to the `sourceos-bootstrap` engine.

```bash
prophet bootstrap build
prophet bootstrap build --output json
```

#### `prophet bootstrap fetch`

Fetch release artifacts from the bootstrap engine's artifact store.

```bash
prophet bootstrap fetch
prophet bootstrap fetch --output json
```

#### `prophet bootstrap write`

Prepare install or recovery media using fetched artifacts.

```bash
prophet bootstrap write
```

#### `prophet bootstrap info`

Show version and configuration information for the installed bootstrap engine.

```bash
prophet bootstrap info
prophet bootstrap info --output json
```

#### `prophet bootstrap validate <kind> <path>`

Validate a frozen object through the bootstrap engine. `<kind>` is the object type (e.g. `ReleaseSet`, `ConfigSource`) and `<path>` is the file system path to the object.

```bash
prophet bootstrap validate ReleaseSet ./releaseset.json
prophet bootstrap validate ConfigSource ./config.yaml --output json
```

#### `prophet bootstrap verify <kind> <path>`

Cryptographically verify a frozen object or artifact through the bootstrap engine.

```bash
prophet bootstrap verify ReleaseSet ./releaseset.json
prophet bootstrap verify BootReleaseSet ./boot.json --output json
```

---

### `prophet a2a`

Agent-to-agent workflow orchestration. Manages the lifecycle of a structured change workflow across propose → test → review → revise → merge → done phases.

#### `prophet a2a run` *(scaffold)*

Start or describe an A2A workflow run for a repository and ticket.

**Flags:**

| Flag | Default | Description |
|---|---|---|
| `--repo` | required | Repository in `owner/name` format |
| `--ticket` | required | Ticket or issue identifier |
| `--live` | `false` | Execute live; default is dry-run |

```bash
# Dry-run: print workflow plan without executing
prophet a2a run --repo owner/repo --ticket TICKET-123

# Live run
prophet a2a run --repo owner/repo --ticket TICKET-123 --live

# JSON output
prophet a2a run --repo owner/repo --ticket TICKET-123 --output json
```

---

## Hybrid overlay placeholders

These commands are scaffolds. They emit a structured stub and will be wired to agent runtimes in a future phase.

### `prophet ask`

Agent assist: explain or inspect without mutating state. Intended for read-only questions about the system, objects, or plans.

```bash
prophet ask
prophet ask --output json
```

### `prophet plan`

Agent assist: generate a plan over deterministic tools. Produces a structured plan without executing any mutations.

```bash
prophet plan
prophet plan --output json
```

### `prophet agent`

Agent execute façade. Approval-gated execution that resolves to deterministic tools. Mutations require explicit approval.

```bash
prophet agent
prophet agent --output json
```

### `prophet mcp`

MCP boundary façade. Interface to the Model Context Protocol boundary for agent integrations.

```bash
prophet mcp
prophet mcp --output json
```

---

## Notes

- All commands support `--output json` for structured, machine-readable output.
- Scaffold commands include a `status: scaffold` field and a `delegated_to` field in their output to make delegation visible.
- The deterministic verbs (`bootstrap`, `a2a`) are the canonical surface for automation. Agent overlays (`ask`, `plan`, `agent`, `mcp`) must not define competing semantics.
