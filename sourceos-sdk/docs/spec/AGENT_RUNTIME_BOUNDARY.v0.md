# AGENT_RUNTIME_BOUNDARY.v0

## Purpose
Define phase-1 runtime boundaries for agent assist/execute surfaces.

## Rules
- Deterministic command/tool invocations are canonical.
- `ask` and `plan` are read-mostly overlays.
- `agent run` requires explicit approval gates.
- Planner components must not directly mutate boot/system state.
- Agent overlays must not redefine frozen object semantics.
