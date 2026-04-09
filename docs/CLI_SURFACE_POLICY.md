# CLI surface policy

The deterministic CLI surface is canonical. This document defines the approved verb taxonomy for `prophet` commands.

## Guiding principles

- **Stability over breadth** — a small set of well-defined verbs is better than a large set of loosely defined ones.
- **No competing command language** — agent overlays (`ask`, `plan`, `agent`, `mcp`) must not define verbs that shadow or replace the deterministic surface.
- **Transparency** — every command must be honest about what it does, what it delegates, and what its current phase/status is.

## Approved verbs by category

### Read verbs

Used for non-mutating queries. Safe in automation and agent-assist contexts.

| Verb | Example |
|---|---|
| `list` | `prophet bootstrap list` |
| `show` | `prophet bootstrap show` |
| `describe` | `prophet bootstrap describe` |
| `explain` | `prophet bootstrap explain` |
| `find` | `prophet bootstrap find` |
| `search` | `prophet bootstrap search` |

### Mutation verbs

Used for state changes. Require explicit invocation; agent-execute layer must surface these to the user before running.

| Verb | Example |
|---|---|
| `create` | `prophet bootstrap create` |
| `update` | `prophet bootstrap update` |
| `delete` | `prophet bootstrap delete` |
| `apply` | `prophet bootstrap apply` |
| `wait` | `prophet bootstrap wait` |
| `validate` | `prophet bootstrap validate <kind> <path>` |
| `verify` | `prophet bootstrap verify <kind> <path>` |

### Runtime verbs

Used for managing running processes or services.

| Verb | Example |
|---|---|
| `status` | `prophet a2a status` |
| `logs` | `prophet a2a logs` |
| `start` | `prophet a2a start` |
| `stop` | `prophet a2a stop` |
| `restart` | `prophet a2a restart` |

### Bootstrap verbs

Specific to the bootstrap domain; delegate to `sourceos-bootstrap`.

| Verb | Example |
|---|---|
| `doctor` | `prophet bootstrap doctor` |
| `login` | `prophet bootstrap login` |
| `build` | `prophet bootstrap build` |
| `fetch` | `prophet bootstrap fetch` |
| `write` | `prophet bootstrap write` |
| `info` | `prophet bootstrap info` |

## Adding new commands

Before adding a new command:

1. Check that the verb fits one of the above categories, or propose a new category in a pull request.
2. Ensure the command emits structured output via the shared `emit` helper.
3. Include `delegated_to` in the output if the command delegates to an external engine.
4. Document the command in [COMMANDS.md](COMMANDS.md).

Agent affordances may exist, but they must resolve to deterministic tools and must not define a competing command language.
