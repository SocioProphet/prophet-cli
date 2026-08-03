"""prophet continuum — the SourceOS Continuum developer continuum, zero-config.

Onboard → develop → cloud-native test → rollout, the same workload from n=1 local
up. prophet discovers the Continuum control plane (via prophet.integrations.
discovery) and delegates to its operator surface (`sourceosctl` / the
sourceos-continuum tools) — no manual configuration.

  prophet continuum status
  prophet continuum up        # bring up the local sovereign forge + cluster
  prophet continuum dev       # develop against the local forge
  prophet continuum test      # cloud-native test on the local cluster (kind/k3s)
  prophet continuum rollout   # roll the same workload out
"""
import subprocess
import sys

import click
from rich.console import Console

from prophet.integrations.discovery import discover_continuum

console = Console()


def _resolve():
    ep = discover_continuum()
    if ep is None or not ep.ok:
        console.print(
            "[red]No SourceOS Continuum found.[/red] Zero-config discovery looked for: "
            "$SOURCEOS_CONTINUUM, `sourceosctl` on PATH, ~/.config/sourceos/continuum.toml, "
            "and a sibling sourceos-continuum checkout. Install `sourceosctl` or set "
            "$SOURCEOS_CONTINUUM."
        )
        sys.exit(1)
    console.print(f"[dim]  continuum: {ep.location} (via {ep.source})[/dim]")
    return ep


def _delegate(ep, args: list[str]):
    if ep.kind == "cli":
        cmd = [ep.location, *args]
    elif ep.kind == "repo":
        # Prefer the repo's own operator entrypoint.
        cmd = ["sourceosctl", *args] if _has("sourceosctl") else ["make", "-C", ep.location, *args]
    else:  # endpoint/config — sourceosctl talks to it via its own resolution
        cmd = ["sourceosctl", *args]
    console.print(f"[dim]  {' '.join(cmd)}[/dim]")
    r = subprocess.run(cmd)
    if r.returncode != 0:
        sys.exit(r.returncode)


def _has(binary: str) -> bool:
    import shutil
    return shutil.which(binary) is not None


@click.group()
def continuum():
    """SourceOS Continuum: onboard → dev → test → rollout (zero-config)."""
    pass


@continuum.command()
def status():
    """Show the discovered Continuum and its state."""
    _delegate(_resolve(), ["status"])


@continuum.command()
def up():
    """Bring up the local sovereign forge + cluster."""
    _delegate(_resolve(), ["up"])


@continuum.command()
def dev():
    """Develop against the local forge."""
    _delegate(_resolve(), ["dev"])


@continuum.command("test")
def test_cmd():
    """Cloud-native test the workload on the local cluster (kind/k3s)."""
    _delegate(_resolve(), ["test"])


@continuum.command()
def rollout():
    """Roll the same workload out (n=1 local → composable cluster)."""
    _delegate(_resolve(), ["rollout"])
