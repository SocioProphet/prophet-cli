# prophet-cli (patch-only façade)

This repository is the **façade-only** command surface for ProphetCLI/ProfitCLI in phase 1.

## Role split
- `prophet-cli`: user-facing wrappers and interactive shell façade.
- `sourceos-sdk/cmd/sourceos-bootstrap`: bootstrap engine source home.
- `homebrew-sourceos`: distribution tap only.

This repository must not embed bootstrap business logic or own frozen object semantics.
