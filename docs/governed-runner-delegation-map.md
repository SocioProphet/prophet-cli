# Governed Runner Delegation Map

Status: v0.1 facade boundary

Owning repo for this document: `SocioProphet/prophet-cli`

## Purpose

This document records the cross-repo delegation boundary for the governed-runner product surface.

`prophet-cli` owns the stable operator-facing command names. It does not own governed-runner implementation, safety semantics, authority state, evidence contracts, or install formulae.

## Operator command path

The product-facing command is:

```bash
prophet governed-runner ...
```

The implementation delegate is:

```bash
sp-run ...
```

Delegation rule:

```text
prophet governed-runner <args> -> sp-run <args>
```

Plane-name alias:

```text
prophet agentplane <args> -> sp-run <args>
```

## Cross-repo ownership

| Concern | Owning repo | Notes |
|---|---|---|
| Product facade | `SocioProphet/prophet-cli` | Owns `prophet governed-runner ...` and `prophet agentplane ...` command names only |
| Governed-runner implementation | `SocioProphet/agentplane` | Owns `sp-run`, contracts, receipts, smoke, inspection, dossiers, local JSON tool adapter |
| Safety/preflight semantics | `SocioProphet/guardrail-fabric` | Owns TrustOps safety policy and runtime action mapping |
| Authority state | `SocioProphet/agent-registry` | Owns authority decisions, current authority state, and restoration decisions |
| Install path | `SocioProphet/homebrew-prophet` | Owns `prophet-cli` and `agentplane` formulae |

## Current read-only surface

`prophet-cli` forwards all of these to `sp-run`:

```bash
prophet governed-runner doctor
prophet governed-runner smoke --output-dir ./.socioprophet/smoke/governed-runner
prophet governed-runner list --runs-root ./.socioprophet/smoke/governed-runner
prophet governed-runner status ./.socioprophet/smoke/governed-runner/run
prophet governed-runner inspect ./.socioprophet/smoke/governed-runner/run
prophet governed-runner tool list-tools
prophet governed-runner tool call governed_runner.doctor --args-json '{}'
prophet governed-runner preflight ./governed-run-contract.json
prophet governed-runner admit ./governed-run-contract.json --preflight ./preflight-receipt.json --authority-state ./agent-authority-current-state.json --projected-cost-usd 0.25
prophet governed-runner dossier ./.socioprophet/smoke/governed-runner/run
prophet governed-runner validate-dossier ./run-dossier.json
```

## Install contract

The expected install path is:

```bash
brew tap SocioProphet/prophet
brew install prophet-cli agentplane
```

`prophet-cli` installs `prophet`.

`agentplane` installs `sp-run`.

If `sp-run` is missing, `prophet governed-runner ...` must fail clearly as a missing delegate. The facade must not emulate or reimplement AgentPlane behavior.

## Non-ownership list for prophet-cli

`prophet-cli` does not own:

- `GovernedRunContract` schema;
- `PreflightReceipt` schema;
- `AttemptAdmissionReceipt` schema;
- rollback evidence schemas;
- `RunDossier` schema;
- `sp-run` implementation;
- governed-runner smoke evidence generation;
- run-store inspection logic;
- local JSON tool adapter logic;
- TrustOps safety semantics;
- Agent Registry authority state;
- Homebrew formulae.

## Hard non-goals for the facade

`prophet-cli` must not add:

- live agent execution;
- verifier execution;
- governed workspace mutation;
- rollback restoration;
- Agent Registry authority update;
- budget settlement;
- provider invocation;
- network activity for a governed run;
- durable memory writeback;
- policy adjudication.

If any future `prophet governed-runner ...` command appears to need one of those capabilities, the implementation belongs in the owning repo after a separate policy-gated design tranche. The facade may only delegate.

## Regression rule

Any change to `prophet governed-runner ...` should preserve this invariant:

```text
prophet-cli parses command names and delegates argv; AgentPlane performs governed-runner logic.
```

Tests should assert forwarded argv, not duplicate AgentPlane behavior.
