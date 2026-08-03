"""prophet tools / prophet doctor — the unified dev-tool surface + E2E health.

`prophet` is the ONE developer entrypoint: the estate's dev/CLI/SDK tools are
surfaced here, and `prophet doctor` proves the zero-config end-to-end wiring
(Continuum + TritFabric) resolves without any manual setup.
"""
import shutil

import click
from rich.console import Console
from rich.table import Table

from prophet.integrations.discovery import discover_all

console = Console()

# The estate dev/CLI/SDK tools prophet unifies. Each is surfaced through prophet
# and/or delegated to; discovery makes them work zero-config.
TOOL_REGISTRY = [
    ("continuum", "sourceosctl", "SourceOS Continuum operator (onboard→dev→test→rollout)"),
    ("ai", "atlas", "TritFabric / Atlas AI runtime + promotion gates"),
    ("kustomize", "kustomize", "manifest build/diff/apply"),
    ("infra", "tofu", "infra plan/apply"),
    ("workspace", "docker", "workspace service lifecycle"),
]


@click.group()
def tools():
    """The unified Prophet dev-tool surface."""
    pass


@tools.command("list")
def list_cmd():
    """List the dev/CLI/SDK tools prophet unifies + whether each is available."""
    t = Table(title="prophet — unified dev tools")
    t.add_column("prophet cmd"); t.add_column("backing tool"); t.add_column("on PATH"); t.add_column("purpose")
    for cmd, backing, purpose in TOOL_REGISTRY:
        avail = "[green]yes[/green]" if shutil.which(backing) else "[yellow]no[/yellow]"
        t.add_row(f"prophet {cmd}", backing, avail, purpose)
    console.print(t)


@click.command()
def doctor():
    """Zero-config E2E health: does prophet resolve the whole continuum + AI plane?"""
    found = discover_all()
    t = Table(title="prophet doctor — zero-config discovery")
    t.add_column("integration"); t.add_column("status"); t.add_column("resolved to"); t.add_column("via")
    ok = True
    for name, ep in found.items():
        if ep and ep.ok:
            t.add_row(name, "[green]OK[/green]", ep.location, ep.source)
        else:
            t.add_row(name, "[red]MISSING[/red]", "-", "-")
            ok = False
    console.print(t)
    if ok:
        console.print("[green]E2E ready:[/green] onboard → dev → test → rollout → AI gates, zero-config.")
    else:
        console.print("[yellow]Some integrations are not discoverable yet[/yellow] — see the rows above; "
                      "install the backing tool or set its env override. prophet fails clearly, never silently.")
        raise SystemExit(1)
