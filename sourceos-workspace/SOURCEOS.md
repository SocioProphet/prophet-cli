# SOURCEOS Workspace Context (Phase 1)

This file is workspace context for operators/agents.
It is **not** canonical policy.

## Frozen object set in phase 1
- ExperienceProfile.v1
- ReleaseSet.v1
- BootReleaseSet.v1
- Fingerprint.v1
- ConfigSource.v0
- TokenDoor.v0
- PolicyBundle.v0

## Boundary reminders
- `prophet-cli` is façade-only.
- `sourceos-sdk/cmd/sourceos-bootstrap` is bootstrap engine source home.
- `homebrew-sourceos` is distribution-only.
- Avoid boot enrollment and DeviceClaim/PAL-Mac final semantics in phase 1.
