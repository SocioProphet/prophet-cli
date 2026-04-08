# Canonical mainline and boundary

This note makes the current intended repository boundary explicit.

## Boundary

`prophet-cli` is the operator façade.
It should provide:

- command grammar
- wrapper behavior
- machine-readable local output
- delegation to owning runtime or bootstrap surfaces

It should not become the owner of:

- bootstrap engine logic
- protocol canon
- transport crypto or frame semantics
- long-lived secret material

## Mainline expectation

Feature work should converge into a canonical mainline rather than living indefinitely in split façade/reference branches.

## Related repos

- `TriTRPC` for protocol truth
- `agentplane` for execution control
- `sociosphere` for workspace control
- the owning bootstrap engine repo for bootstrap implementation
