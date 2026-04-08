# Bootstrap delegation

`prophet-cli` is a façade. Bootstrap actions delegate to the installed `sourceos-bootstrap` engine.

## Delegated commands
- `prophet bootstrap doctor`
- `prophet bootstrap login`
- `prophet bootstrap build`
- `prophet bootstrap fetch`
- `prophet bootstrap write`
- `prophet bootstrap info`
- `prophet bootstrap validate <kind> <path>`
- `prophet bootstrap verify <kind> <path>`

## Rules
- do not implement bootstrap business logic here
- do not claim ownership of ReleaseSet, BootReleaseSet, ConfigSource, or TokenDoor semantics here
- wrapper commands should remain transparent about delegation
