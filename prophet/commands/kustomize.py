"""
prophet kustomize — Kustomize build, diff, and apply helpers.

Commands:
  prophet kustomize build  <service> [--overlay OVERLAY]
  prophet kustomize diff   <service> [--overlay OVERLAY]
  prophet kustomize apply  <service> [--overlay OVERLAY] [--dry-run]
  prophet kustomize list
"""
import subprocess
import sys
from pathlib import Path

import click
from rich.console import Console
from rich.table import Table

console = Console()

PLATFORM_ROOT = Path(__file__).parent.parent.parent.parent / "prophet-platform"
K8S_ROOT = PLATFORM_ROOT / "infra" / "k8s"


def _kust_path(service: str, overlay: str) -> Path:
    p = K8S_ROOT / service / "overlays" / overlay
    if not p.exists():
        # Fall back to base
        p = K8S_ROOT / service / "base"
    if not p.exists():
        console.print(f"[red]No kustomize dir found for service={service} overlay={overlay}[/red]")
        sys.exit(1)
    return p


def _run(cmd: list[str]):
    console.print(f"[dim]  {' '.join(cmd)}[/dim]")
    result = subprocess.run(cmd, text=True)
    if result.returncode != 0:
        sys.exit(result.returncode)


@click.group("kustomize")
def kustomize_cmd():
    """Kustomize build / diff / apply helpers."""
    pass


@kustomize_cmd.command("build")
@click.argument("service")
@click.option("--overlay", "-o", default="p0-lab", show_default=True)
def build(service, overlay):
    """Render Kustomize manifests to stdout."""
    path = _kust_path(service, overlay)
    console.print(f"[cyan]Building {service}/{overlay}[/cyan]")
    _run(["kubectl", "kustomize", str(path)])


@kustomize_cmd.command("diff")
@click.argument("service")
@click.option("--overlay", "-o", default="p0-lab", show_default=True)
def diff(service, overlay):
    """Diff rendered manifests against the live cluster state."""
    path = _kust_path(service, overlay)
    console.print(f"[cyan]Diffing {service}/{overlay} against live cluster[/cyan]")
    _run(["kubectl", "diff", "-k", str(path)])


@kustomize_cmd.command("apply")
@click.argument("service")
@click.option("--overlay", "-o", default="p0-lab", show_default=True)
@click.option("--dry-run", is_flag=True, default=False)
def apply(service, overlay, dry_run):
    """Apply Kustomize manifests to the current kubeconfig cluster."""
    path = _kust_path(service, overlay)
    mode = "client" if dry_run else None
    console.print(f"[cyan]Applying {service}/{overlay}{'  (dry-run)' if dry_run else ''}[/cyan]")
    cmd = ["kubectl", "apply", "-k", str(path)]
    if dry_run:
        cmd += ["--dry-run=client"]
    _run(cmd)


@kustomize_cmd.command("list")
def list_services():
    """List all kustomize services with available overlays."""
    table = Table(title="Kustomize Services", show_header=True)
    table.add_column("Service", style="bold cyan")
    table.add_column("Overlays")
    table.add_column("Base")

    for service_dir in sorted(K8S_ROOT.iterdir()):
        if not service_dir.is_dir() or service_dir.name in ("argo-cd", "namespaces"):
            continue
        overlays_dir = service_dir / "overlays"
        overlays = sorted(p.name for p in overlays_dir.iterdir()) if overlays_dir.exists() else []
        base_exists = "✓" if (service_dir / "base").exists() else "✗"
        table.add_row(service_dir.name, ", ".join(overlays) or "—", base_exists)

    console.print(table)
