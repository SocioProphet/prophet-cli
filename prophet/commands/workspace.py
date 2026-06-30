"""
prophet workspace — workspace service lifecycle.

Commands:
  prophet workspace validate
  prophet workspace test
  prophet workspace build  [--env local]
  prophet workspace up     [--env local]
  prophet workspace down   [--env local]
  prophet workspace logs   [--env local] [--service SVC]
  prophet workspace smoke  [--env local]
  prophet workspace status [--env local]
"""
import subprocess
import sys
from pathlib import Path

import click
from rich.console import Console
from rich.table import Table

console = Console()

PLATFORM_ROOT = Path(__file__).parent.parent.parent.parent / "prophet-platform"
COMPOSE_MAP = {
    "local": PLATFORM_ROOT / "infra" / "local" / "docker-compose.workspace.yml",
}


def _check_docker():
    r = subprocess.run(["docker", "info"], capture_output=True)
    if r.returncode != 0:
        console.print("[red]Docker daemon is not running. Start Docker Desktop and retry.[/red]")
        sys.exit(1)


def _compose(env: str, args: list[str]):
    compose_file = COMPOSE_MAP.get(env)
    if not compose_file:
        console.print(f"[red]Unknown workspace env: {env}. Available: {list(COMPOSE_MAP)}[/red]")
        sys.exit(1)
    cmd = ["docker", "compose", "-f", str(compose_file)] + args
    console.print(f"[dim]  {' '.join(cmd)}[/dim]")
    result = subprocess.run(cmd)
    if result.returncode != 0:
        sys.exit(result.returncode)


def _make(target: str):
    cmd = ["make", "-C", str(PLATFORM_ROOT), target]
    console.print(f"[dim]  {' '.join(cmd)}[/dim]")
    result = subprocess.run(cmd)
    if result.returncode != 0:
        sys.exit(result.returncode)


@click.group()
def workspace():
    """Workspace service lifecycle (mail, CalDAV, MinIO, etc.)."""
    pass


@workspace.command()
def validate():
    """Validate workspace configs and manifests (no Docker required)."""
    console.print("[cyan]Validating workspace service configs...[/cyan]")
    _make("validate-workspace-services")


@workspace.command("test")
def test_cmd():
    """Run pytest workspace infra test suite (no Docker required)."""
    console.print("[cyan]Running workspace infra tests...[/cyan]")
    _make("test-workspace")


@workspace.command()
@click.option("--env", "-e", default="local", show_default=True)
def build(env):
    """Build workspace service Docker images."""
    _check_docker()
    console.print(f"[cyan]Building workspace images (env={env})...[/cyan]")
    _compose(env, ["build"])
    console.print("[green]Build complete.[/green]")


@workspace.command()
@click.option("--env", "-e", default="local", show_default=True)
def up(env):
    """Start workspace services."""
    _check_docker()
    console.print(f"[cyan]Starting workspace services (env={env})...[/cyan]")
    _compose(env, ["up", "-d"])
    console.print("[green]Services started.[/green]")
    console.print("  IMAP  : localhost:143")
    console.print("  SMTP  : localhost:25   (submission: 587)")
    console.print("  CalDAV: http://localhost:5232")
    console.print("  MinIO : http://localhost:9000  (console: 9001)")
    console.print("  PG    : localhost:5432")
    console.print("  Redis : localhost:6379")
    console.print("\n[dim]prophet workspace logs --env local[/dim]")


@workspace.command()
@click.option("--env", "-e", default="local", show_default=True)
def down(env):
    """Stop and remove workspace services."""
    _check_docker()
    _compose(env, ["down"])
    console.print("[green]Workspace stopped.[/green]")


@workspace.command()
@click.option("--env", "-e", default="local", show_default=True)
@click.option("--service", "-s", default="", help="Follow a specific service (blank = all)")
def logs(env, service):
    """Tail workspace service logs."""
    _check_docker()
    args = ["logs", "-f"]
    if service:
        args.append(service)
    _compose(env, args)


@workspace.command()
@click.option("--env", "-e", default="local", show_default=True)
def smoke(env):
    """Run full smoke test (Docker required — builds, starts, probes, tears down)."""
    _check_docker()
    console.print("[cyan]Running workspace smoke test...[/cyan]")
    _make("smoke-workspace")


@workspace.command()
@click.option("--env", "-e", default="local", show_default=True)
def status(env):
    """Show running status of workspace services."""
    _check_docker()
    console.print(f"[cyan]Workspace status (env={env})[/cyan]")
    _compose(env, ["ps"])
