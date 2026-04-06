# Bootstrap delegation

`prophet-cli` is a façade. Bootstrap actions delegate to the installed `sourceos-bootstrap` engine.

Phase 1 commands:
- `prophet bootstrap doctor`
- `prophet bootstrap login`
- `prophet bootstrap build`
- `prophet bootstrap validate <kind> <path>`
- `prophet bootstrap verify <kind> <path>`

Non-goals:
- no bootstrap business logic in this repo
- no BootReleaseSet, ConfigSource, or TokenDoor ownership here
