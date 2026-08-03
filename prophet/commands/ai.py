"""prophet ai — the TritFabric / Atlas AI runtime + promotion gates, zero-config.

prophet discovers the Atlas endpoint (via prophet.integrations.discovery) and
delegates to TritFabric's Atlas surface (gRPC-first, REST-mirror). Promotion is
fail-closed by default (SHACL + ONNX round-trip + eval deltas).

  prophet ai status
  prophet ai serve MODEL      # serve a model via the Atlas router/autoscaler
  prophet ai gate  MODEL      # run the fail-closed promotion gates
  prophet ai promote MODEL    # promote a model that passed its gates
"""
import subprocess
import sys

import click
from rich.console import Console

from prophet.integrations.discovery import discover_tritfabric

console = Console()


def _resolve():
    ep = discover_tritfabric()
    if ep is None or not ep.ok:
        console.print(
            "[red]No TritFabric / Atlas found.[/red] Zero-config discovery looked for: "
            "$TRITFABRIC_ENDPOINT / $ATLAS_ENDPOINT, `atlas` on PATH, "
            "~/.config/tritfabric/atlas.toml, and a sibling tritfabric checkout. "
            "Install `atlas` or set $ATLAS_ENDPOINT."
        )
        sys.exit(1)
    console.print(f"[dim]  tritfabric: {ep.location} (via {ep.source})[/dim]")
    return ep


def _delegate(ep, args: list[str]):
    import shutil
    if ep.kind == "endpoint":
        cmd = ["atlas", "--endpoint", ep.location, *args] if shutil.which("atlas") else None
        if cmd is None:
            console.print(f"[yellow]Atlas endpoint {ep.location} discovered but no `atlas` client on PATH.[/yellow]")
            sys.exit(1)
    elif ep.kind == "cli":
        cmd = [ep.location, *args]
    else:  # repo checkout — run its Atlas CLI module from the repo root
        if shutil.which("atlas"):
            cmd = ["atlas", *args]
        else:
            cmd = ["python3", "-m", "cli.atlas", *args]  # tritfabric's cli/ package
    console.print(f"[dim]  {' '.join(cmd)}[/dim]")
    r = subprocess.run(cmd, cwd=(ep.location if ep.kind == "repo" else None))
    if r.returncode != 0:
        sys.exit(r.returncode)


@click.group()
def ai():
    """TritFabric / Atlas: AI runtime + fail-closed promotion gates (zero-config)."""
    pass


@ai.command()
def status():
    """Show the discovered Atlas endpoint and its state."""
    _delegate(_resolve(), ["status"])


@ai.command()
@click.argument("model")
def serve(model):
    """Serve MODEL via the Atlas router/autoscaler."""
    _delegate(_resolve(), ["serve", model])


@ai.command()
@click.argument("model")
def gate(model):
    """Run the fail-closed promotion gates on MODEL (SHACL + ONNX + eval)."""
    _delegate(_resolve(), ["gate", model])


@ai.command()
@click.argument("model")
def promote(model):
    """Promote MODEL once its gates pass."""
    _delegate(_resolve(), ["promote", model])
