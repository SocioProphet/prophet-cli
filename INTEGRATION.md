# prophet — the one developer entrypoint (zero-config E2E)

`prophet` unifies the estate's dev/CLI/SDK tools and wires them **end-to-end with
zero configuration** to the two systems a developer needs:

- **SourceOS Continuum** — onboard → dev → cloud-native test → rollout (the same
  workload from n=1 local up). Operator surface: `sourceosctl`.
- **TritFabric / Atlas** — the AI runtime + **fail-closed promotion gates**
  (SHACL + ONNX round-trip + eval deltas), gRPC-first.

## Zero-config contract
`prophet` never asks you to configure endpoints. `prophet.integrations.discovery`
resolves each target in order, first hit wins, and a miss is **reported clearly,
never guessed**:

| order | continuum | tritfabric |
|---|---|---|
| 1 | `$SOURCEOS_CONTINUUM` | `$TRITFABRIC_ENDPOINT` / `$ATLAS_ENDPOINT` |
| 2 | `sourceosctl` on PATH | `atlas` on PATH |
| 3 | `~/.config/sourceos/continuum.toml` | `~/.config/tritfabric/atlas.toml` |
| 4 | sibling `sourceos-continuum` checkout | sibling `tritfabric` checkout |

## Commands
```
prophet doctor                 # is the whole E2E wired? (zero-config health)
prophet tools list             # the unified dev/CLI/SDK tool surface
prophet continuum up|dev|test|rollout|status
prophet ai serve|gate|promote|status MODEL
prophet infra|workspace|kustomize ...        # existing surfaces
```
`prophet doctor` is the E2E gate: green ⇒ onboard→dev→test→rollout→AI-gates works
with no setup; otherwise it names the missing piece and how to resolve it.

## Extending
Add a tool by registering it in `prophet/commands/tools.py::TOOL_REGISTRY` and, if
it's an integration target, a resolver in `prophet/integrations/discovery.py`
(keep it pure + add a both-ways test in `tests/test_discovery.py`).
