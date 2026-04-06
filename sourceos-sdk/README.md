# sourceos-sdk

SDK and bootstrap engine source home for SourceOS phase 1.

## Source home
Bootstrap engine source home is:
- `cmd/sourceos-bootstrap`

## Boundaries
- Deterministic CLI surface is canonical.
- Agent affordances are overlays that invoke deterministic tools.
- This repo does not finalize enrollment-token, DeviceClaim, boot provisioning protocol, or PAL-Mac semantics.
