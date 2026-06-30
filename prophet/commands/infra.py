"""
prophet infra — Terraform lifecycle for Prophet environments.

Commands:
  prophet infra init   [--env ENV]
  prophet infra plan   [--env ENV]
  prophet infra apply  [--env ENV] [--auto-approve]
  prophet infra destroy [--env ENV] [--auto-approve]
  prophet infra output [--env ENV] [--json]
  prophet infra envs
"""
import os
import subprocess
import sys
from pathlib import Path

import click
from rich.console import Console
from rich.table import Table

console = Console()

PLATFORM_ROOT = Path(__file__).parent.parent.parent.parent / "prophet-platform"
TERRAFORM_ROOT = PLATFORM_ROOT / "infra" / "terraform" / "environments"

ENVS = {
    "p0-lab":         "Local k3d cluster (Docker required)",
    "p1-single-site": "Hetzner Cloud single-site k3s cluster",
    "p2-private-mesh": "Multi-site private mesh (stub)",
    "p3-regional":    "Regional multi-cluster (stub)",
}


def _tf_dir(env: str) -> Path:
    d = TERRAFORM_ROOT / env
    if not d.exists():
        console.print(f"[red]No Terraform environment found for '{env}' at {d}[/red]")
        sys.exit(1)
    return d


def _run_tf(env: str, args: list[str], capture: bool = False):
    tf_dir = _tf_dir(env)
    cmd = ["terraform"] + args
    console.print(f"[dim]  cwd: {tf_dir}[/dim]")
    console.print(f"[dim]  cmd: {' '.join(cmd)}[/dim]")
    result = subprocess.run(cmd, cwd=tf_dir, capture_output=capture, text=True)
    if result.returncode != 0:
        if capture and result.stderr:
            console.print(f"[red]{result.stderr}[/red]")
        sys.exit(result.returncode)
    return result


def _check_terraform():
    if subprocess.run(["terraform", "version"], capture_output=True).returncode != 0:
        console.print("[red]terraform not found in PATH. Install from https://developer.hashicorp.com/terraform/downloads[/red]")
        sys.exit(1)


@click.group()
def infra():
    """Terraform IaC lifecycle for Prophet environments."""
    pass


@infra.command()
@click.option("--env", "-e", default="p0-lab", show_default=True, help="Environment name")
def init(env):
    """Initialize Terraform for an environment (downloads providers)."""
    _check_terraform()
    console.print(f"[cyan]Initializing Terraform for env=[bold]{env}[/bold][/cyan]")
    _run_tf(env, ["init", "-upgrade"])
    console.print(f"[green]Init complete.[/green]")


@infra.command()
@click.option("--env", "-e", default="p0-lab", show_default=True)
@click.option("--out", default="", help="Write plan to file (e.g. plan.out)")
def plan(env, out):
    """Show Terraform execution plan."""
    _check_terraform()
    console.print(f"[cyan]Planning env=[bold]{env}[/bold][/cyan]")
    args = ["plan"]
    if out:
        args += [f"-out={out}"]
    _run_tf(env, args)


@infra.command()
@click.option("--env", "-e", default="p0-lab", show_default=True)
@click.option("--auto-approve", is_flag=True, default=False)
@click.option("--plan-file", default="", help="Apply a saved plan file")
def apply(env, auto_approve, plan_file):
    """Apply Terraform changes."""
    _check_terraform()
    if env in ("p1-single-site", "p2-private-mesh", "p3-regional") and not auto_approve:
        console.print(f"[yellow]WARNING: This will provision real cloud resources for env={env}.[/yellow]")
        click.confirm("Continue?", abort=True)
    console.print(f"[cyan]Applying env=[bold]{env}[/bold][/cyan]")
    args = ["apply"]
    if auto_approve:
        args.append("-auto-approve")
    if plan_file:
        args.append(plan_file)
    _run_tf(env, args)
    console.print(f"[green]Apply complete.[/green]")


@infra.command()
@click.option("--env", "-e", default="p0-lab", show_default=True)
@click.option("--auto-approve", is_flag=True, default=False)
def destroy(env, auto_approve):
    """Destroy all Terraform-managed resources for an environment."""
    _check_terraform()
    console.print(f"[red]DESTROY: This will destroy all resources for env={env}.[/red]")
    if not auto_approve:
        click.confirm("Are you sure?", abort=True)
    args = ["destroy"]
    if auto_approve:
        args.append("-auto-approve")
    _run_tf(env, args)
    console.print(f"[green]Destroy complete.[/green]")


@infra.command("output")
@click.option("--env", "-e", default="p0-lab", show_default=True)
@click.option("--json", "as_json", is_flag=True, default=False)
def output_cmd(env, as_json):
    """Show Terraform outputs for an environment."""
    _check_terraform()
    args = ["output"]
    if as_json:
        args.append("-json")
    result = _run_tf(env, args, capture=not as_json)
    if as_json:
        pass  # already printed by subprocess
    else:
        console.print(result.stdout)


@infra.command("envs")
def list_envs():
    """List available Terraform environments."""
    table = Table(title="Prophet Environments", show_header=True)
    table.add_column("Env", style="bold cyan")
    table.add_column("Description")
    table.add_column("TF Dir Exists", justify="center")

    for env, desc in ENVS.items():
        exists = "✓" if (TERRAFORM_ROOT / env).exists() else "✗"
        table.add_row(env, desc, exists)

    console.print(table)
