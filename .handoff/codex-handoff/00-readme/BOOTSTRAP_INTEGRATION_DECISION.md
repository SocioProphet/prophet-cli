# Bootstrap Integration Decision — Revision 2.1

## Decision
Use a façade + engine model.

- Existing ProphetCLI/ProfitCLI repo = façade
- `sourceos-sdk/cmd/sourceos-bootstrap` = bootstrap engine source home
- `homebrew-sourceos` = distribution tap

## Rationale
This keeps:
- user-facing CLI continuity
- a narrow host bootstrap engine
- independent distribution packaging
- clear separation between wrapper commands and business logic

## Consequences
- ProphetCLI gets wrapper commands only.
- Bootstrap implementation logic does not live in ProphetCLI.
- Homebrew docs must explain façade/engine split.
