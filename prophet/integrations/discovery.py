"""Zero-configuration discovery for the Prophet developer continuum.

The point of the unification: a developer runs `prophet up` / `prophet ai ...`
and the CLI finds the SourceOS Continuum control plane and the TritFabric (Atlas)
AI runtime **without any manual configuration** — env override → well-known
config → sibling checkout → running-endpoint default, in that order. The first
hit wins; nothing found yields a clear, actionable message (never a silent guess).

This module is pure + offline-testable (no network, no subprocess); the command
groups consume its result. See tests/test_discovery.py for the both-ways proof.
"""
from __future__ import annotations

import os
import shutil
from dataclasses import dataclass
from pathlib import Path
from typing import Optional


@dataclass(frozen=True)
class Endpoint:
    """A discovered integration target."""
    name: str                 # "continuum" | "tritfabric"
    kind: str                 # how it resolves: "cli" | "repo" | "endpoint" | "config"
    location: str             # a CLI name, a repo path, or a URL
    source: str               # which discovery rule matched (for `prophet doctor`)

    @property
    def ok(self) -> bool:
        return bool(self.location)


# Sibling repos are conventionally checked out next to prophet-cli under ~/dev.
_DEV_ROOT = Path(__file__).resolve().parents[3]


def _first(*candidates: Optional[Endpoint]) -> Optional[Endpoint]:
    for c in candidates:
        if c is not None:
            return c
    return None


def _env(name: str, target: str, kind: str) -> Optional[Endpoint]:
    val = os.environ.get(name)
    return Endpoint(target, kind, val, f"env:{name}") if val else None


def _cli(bin_name: str, target: str) -> Optional[Endpoint]:
    path = shutil.which(bin_name)
    return Endpoint(target, "cli", path, f"PATH:{bin_name}") if path else None


def _config(rel: str, target: str) -> Optional[Endpoint]:
    p = Path(os.environ.get("XDG_CONFIG_HOME", Path.home() / ".config")) / rel
    return Endpoint(target, "config", str(p), f"config:{rel}") if p.exists() else None


def _repo(name: str, target: str) -> Optional[Endpoint]:
    p = _DEV_ROOT / name
    return Endpoint(target, "repo", str(p), f"repo:{name}") if (p / ".git").exists() else None


def discover_continuum(dev_root: Path = _DEV_ROOT) -> Optional[Endpoint]:
    """SourceOS Continuum (onboard→dev→test→rollout). Operator surface is
    `sourceosctl` (from sourceos-devtools); the control plane lives in
    sourceos-continuum."""
    return _first(
        _env("SOURCEOS_CONTINUUM", "continuum", "endpoint"),
        _cli("sourceosctl", "continuum"),
        _config("sourceos/continuum.toml", "continuum"),
        _repo_at(dev_root, "sourceos-continuum", "continuum"),
    )


def discover_tritfabric(dev_root: Path = _DEV_ROOT) -> Optional[Endpoint]:
    """TritFabric / Atlas OS Service (AI runtime + promotion gates). gRPC-first
    (api/atlas/v1/atlas.proto), REST-mirror friendly."""
    return _first(
        _env("TRITFABRIC_ENDPOINT", "tritfabric", "endpoint"),
        _env("ATLAS_ENDPOINT", "tritfabric", "endpoint"),
        _cli("atlas", "tritfabric"),
        _config("tritfabric/atlas.toml", "tritfabric"),
        _repo_at(dev_root, "tritfabric", "tritfabric"),
    )


def _repo_at(dev_root: Path, name: str, target: str) -> Optional[Endpoint]:
    p = dev_root / name
    return Endpoint(target, "repo", str(p), f"repo:{name}") if (p / ".git").exists() else None


def discover_all(dev_root: Path = _DEV_ROOT) -> dict[str, Optional[Endpoint]]:
    return {
        "continuum": discover_continuum(dev_root),
        "tritfabric": discover_tritfabric(dev_root),
    }
