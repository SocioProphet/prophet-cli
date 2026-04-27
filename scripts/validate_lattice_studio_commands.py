#!/usr/bin/env python3
"""Validate Prophet CLI Lattice Studio command manifest."""

from __future__ import annotations

import json
import sys
from pathlib import Path

REQUIRED_COMMANDS = {
    "create-session",
    "emit-demo-catalog",
    "emit-platform-records",
    "emit-atlas-context",
    "emit-ontogenesis-context",
    "emit-paas-plan",
    "emit-local-dev",
    "emit-memory",
    "emit-lampstand-demo",
    "emit-notebook-plane",
    "emit-execution",
    "lattice-byoc",
    "lattice-m2-placement",
    "lattice-nb-run",
    "lattice-promote",
    "lattice-place",
    "lattice-safe-plan",
}
REQUIRED_SURFACES = {
    "byoc",
    "cloudshell-fog",
    "m2-topolvm",
    "notebook-launch",
    "notebook-promotion",
    "placement-decision",
    "safe-placement-execution",
}


def require(condition: bool, message: str) -> None:
    if not condition:
        raise ValueError(message)


def validate(path: Path) -> None:
    doc = json.loads(path.read_text(encoding="utf-8"))
    require(doc.get("apiVersion") == "cli.socioprophet.dev/v1", "apiVersion mismatch")
    require(doc.get("kind") == "ProphetCliCommandSurface", "kind mismatch")
    delegates = doc.get("delegatesTo")
    require(isinstance(delegates, dict), "delegatesTo must be object")
    require(delegates.get("repo") == "SocioProphet/prophet-platform", "delegate repo mismatch")
    commands = set(doc.get("commands", []))
    require(REQUIRED_COMMANDS.issubset(commands), f"missing commands: {sorted(REQUIRED_COMMANDS - commands)}")
    surfaces = set(doc.get("surfaces", []))
    require(REQUIRED_SURFACES.issubset(surfaces), f"missing surfaces: {sorted(REQUIRED_SURFACES - surfaces)}")


def main() -> int:
    try:
        validate(Path("manifests/lattice-studio-commands.json"))
        print("PASS manifests/lattice-studio-commands.json")
        return 0
    except Exception as exc:  # noqa: BLE001
        print(f"FAIL manifests/lattice-studio-commands.json: {exc}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    raise SystemExit(main())
