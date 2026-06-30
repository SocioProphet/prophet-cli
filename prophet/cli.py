import click
from rich.console import Console

from prophet.commands.infra import infra
from prophet.commands.workspace import workspace
from prophet.commands.kustomize import kustomize_cmd

console = Console()

@click.group()
@click.version_option("0.1.0", prog_name="prophet")
def main():
    """Prophet platform CLI — infra, workspace, and dev tooling."""
    pass

main.add_command(infra)
main.add_command(workspace)
main.add_command(kustomize_cmd, name="kustomize")
