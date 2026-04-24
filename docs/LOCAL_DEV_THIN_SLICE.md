# Local Dev Thin Slice via `prophet-cli`

This document records how `prophet-cli` should drive the local thin-slice lifecycle owned by `prophet-platform`.

## Repository split

- `prophet-platform` owns runtime truth, bootstrap files, services, charts, and smoke paths.
- `prophet-cli` owns the command façade that calls into those workflows.

The CLI should therefore **delegate** to the runtime repo instead of re-implementing the same orchestration logic.

## Expected local checkout shape

A simple sibling checkout layout is assumed for local development:

```text
~/dev/
  prophet-cli/
  prophet-platform/
```

The CLI can also support an explicit `--platform-root` flag later, but the sibling layout is the easiest first assumption.

## Thin-slice command surface

Initial command intents:

- `prophet dev up --mode compose`
- `prophet dev up --mode kind`
- `prophet dev up --mode k3d`
- `prophet dev e2e-local`
- `prophet dev destroy`
- `prophet train run --spec <path>`
- `prophet model register --from <metrics>`
- `prophet model promote --name <name> --stage prod`
- `prophet infer --values <csv>`

## Current delegation rule

Until `prophet-cli` has a richer implementation, these commands should delegate to scripts and Make targets in `prophet-platform`.

That means the CLI remains thin and the runtime repo remains canonical.

## Why this matters

The local thin slice is the first executable truth path.
If `prophet-cli` cannot drive it, then the command façade is not yet real.
