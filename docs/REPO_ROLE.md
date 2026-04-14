# Repo Role

`prophet-cli` is the facade repository for the Prophet command surface.

## It owns
- command grammar
- wrapper command definitions
- interactive and hybrid CLI documentation
- tests that protect command-surface shape

## It does not own
- bootstrap engine business logic
- final transport internals
- final receipt or enrollment semantics
- ReleaseSet, BootReleaseSet, ConfigSource, or TokenDoor ownership

## Engine boundary

The bootstrap engine source home for phase 1 remains:
- `sourceos-sdk/cmd/sourceos-bootstrap`
