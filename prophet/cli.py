import click
from rich.console import Console

from prophet.commands.infra import infra
from prophet.commands.workspace import workspace
from prophet.commands.kustomize import kustomize_cmd
from prophet.commands.continuum import continuum
from prophet.commands.ai import ai
from prophet.commands.tools import tools, doctor

console = Console()

@click.group()
@click.version_option("0.1.0", prog_name="prophet")
def main():
    """Prophet CLI — the one developer entrypoint.

    Unifies the estate's dev/CLI/SDK tools with zero-config end-to-end wiring to
    SourceOS Continuum (onboard→dev→test→rollout) and TritFabric/Atlas (AI
    runtime + fail-closed promotion gates). Run `prophet doctor` to see the wiring.
    """
    pass

main.add_command(infra)
main.add_command(workspace)
main.add_command(kustomize_cmd, name="kustomize")
main.add_command(continuum)
main.add_command(ai)
main.add_command(tools)
main.add_command(doctor)
