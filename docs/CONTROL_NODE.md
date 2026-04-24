# control-node command surface

`prophet control-node` is the operator-facing façade for the local-first control-node lane.

## Current subcommands

- `prophet control-node status`
  - emits the current delegated repo surfaces for contracts, runtime, and execution
- `prophet control-node process <input.json> <outdir>`
  - delegates to the `agentplane` local-control-node processor and returns a structured probe envelope

## Current delegation split

- contracts: `SourceOS-Linux/sourceos-spec`
- runtime: `SocioProphet/prophet-platform`
- execution/evidence: `SocioProphet/agentplane`
- workstation/bootstrap: `SociOS-Linux/source-os`

## Why this lives here

`prophet-cli` is the façade command surface. It should expose the operator workflow and delegate into the underlying implementation repos instead of re-owning those semantics.
