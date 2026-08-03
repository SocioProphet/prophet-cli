# Implementation Status

Current-state snapshot of the `prophet-cli` façade. This file is repo-local working
memory: it records what the command surface actually does today — real vs delegating
vs scaffold — so an operator or next agent can orient without re-deriving it from source.

- **Repo role:** façade only. See [REPO_ROLE.md](REPO_ROLE.md). We own command grammar,
  wrapper definitions, docs, and surface-shape tests. We do **not** own bootstrap engine
  logic, transport internals, receipt/enrollment semantics, or ReleaseSet/ConfigSource/TokenDoor.
- **Language/toolchain:** Go 1.22, Cobra. Build: `go build ./...`. Test: `go test ./...`.
- **Last updated:** 2026-08-03

## Boundary legend

| kind | meaning |
| --- | --- |
| `real` | executes locally in the façade with no external engine required |
| `delegating` | shells out to a named engine binary (or a local-dev repo fallback); reports `not-yet-wired`/`not-yet-installed` when the engine is absent |
| `scaffold` | placeholder surface; emits `status: scaffold` and mutates nothing |

The live boundary is queryable at runtime: `prophet status` enumerates the same surfaces
and probes which delegate engines are installed. This document is the human-readable twin.

## Top-level surface

| command | kind | delegate engine | notes |
| --- | --- | --- | --- |
| `version` / `doctor` / `self-test` / `emit-evidence` | real | — | suite diagnostics, façade-local |
| `status` | real | — | boundary legibility + delegate presence probe |
| `bootstrap` | delegating | `sourceos-bootstrap` | engine home: `sourceos-sdk/cmd/sourceos-bootstrap` |
| `vocab` | delegating | ontogenesis vocab surface | fetch/gate/promote/sr |
| `bindings` | delegating | atomic bindings engine | validate |
| `k8s` | delegating | k8s policy engine | scheduling checks |
| `control-node` | delegating | local control-node | status/local-first |
| `a2a` | scaffold | — | workflow façade skeleton |
| `devtools` / `lab` | delegating | `sourceos-devtools` | profile/lab management; `not-yet-wired` |
| `sourceos install` | delegating | `sourceos-installer` | `not-yet-wired` |
| `sourceos carry` | delegating | `sourceos-ai` | list/validate/doctor/emit-evidence |
| `holmes` | delegating | `holmes` | analyze/search/graph/govern |
| `model route` | delegating | `model-router` | local-dev python fallback |
| `guardrail test` | delegating | `guardrail-fabric` | local-dev python fallback |
| `ledger` | delegating | `model-governance-ledger` | validate/records; local-dev fallback |
| `agent registry` | delegating | `agent-registry` | list; local-dev record fallback |
| `spine` | delegating | per-repo spine gates | validate gates (`--repo`) |
| `enrichment` | delegating | enrichment twin | corpus/lifecycle/gate |
| `ask` / `plan` / `mcp` | scaffold | — | agent-assist / MCP boundary placeholders |

## Delegate engines (external, not owned here)

`sourceos-bootstrap`, `sourceos-ai`, `sourceos-devtools`, `sourceos-installer`, `holmes`,
`model-router`, `guardrail-fabric`, `model-governance-ledger`, `agent-registry`.

Delegating commands resolve engines via `exec.LookPath`, then a local-dev repo fallback
rooted at `$PROPHET_DEV_ROOT` (default `~/dev/<repo>`). Absence is reported, never faked.

## Known gaps / next weakest link

- `devtools`, `lab`, and `sourceos install` are declared but return `not-yet-wired` — the
  façade grammar is ahead of engine wiring for these surfaces.
- `a2a`, `ask`, `plan`, `mcp` are scaffolds only.
- Receipts (`--receipt`) are façade-local and command-scoped; they are **not** the estate
  ProofArtifact spine and must not be treated as a runtime ledger.
