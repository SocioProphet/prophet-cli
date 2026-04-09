# Bootstrap delegation

`prophet-cli` is a façade. All `prophet bootstrap` subcommands delegate to the installed `sourceos-bootstrap` engine. This repo owns the command surface, flag grammar, and output format — it does **not** own bootstrap business logic.

## How delegation works

Every `prophet bootstrap` command emits a response that includes a `delegated_to` field:

```json
{
  "command": "prophet bootstrap doctor",
  "summary": "diagnose host prerequisites",
  "delegated_to": "sourceos-bootstrap",
  "status": "scaffold"
}
```

In phase 1 this is a scaffold response. In later phases the command will invoke the `sourceos-bootstrap` binary and stream its output back through the selected `--output` format.

## Delegated commands

| Command | Delegated action |
|---|---|
| `prophet bootstrap doctor` | Diagnose host prerequisites |
| `prophet bootstrap login` | Prepare authenticated bootstrap session |
| `prophet bootstrap build` | Submit or describe a bootstrap build request |
| `prophet bootstrap fetch` | Fetch release artifacts |
| `prophet bootstrap write` | Prepare install or recovery media |
| `prophet bootstrap info` | Show bootstrap engine information |
| `prophet bootstrap validate <kind> <path>` | Validate a frozen object |
| `prophet bootstrap verify <kind> <path>` | Cryptographically verify a frozen object or artifact |

## Frozen object types

The `validate` and `verify` commands accept a `<kind>` argument. Recognised kinds are defined by `sourceos-bootstrap` and include (but are not limited to):

- `ReleaseSet`
- `BootReleaseSet`
- `ConfigSource`
- `TokenDoor`

`prophet-cli` does not claim ownership of these type definitions. It passes the kind argument transparently to the bootstrap engine.

## Rules

- Do not implement bootstrap business logic in this repo.
- Do not claim ownership of `ReleaseSet`, `BootReleaseSet`, `ConfigSource`, or `TokenDoor` semantics here.
- Wrapper commands must remain transparent about delegation; always include `delegated_to` in structured output.
- Do not add enrollment or auto-trust logic to bootstrap wrappers.
