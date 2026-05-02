"""Prophet facade CLI.

The `prophet` command is a stable operator-facing facade. It delegates local
SourceOS implementation work to the owning CLIs instead of duplicating engine
logic here.
"""

from __future__ import annotations

import argparse
import shutil
import subprocess
import sys
from typing import Sequence

from prophet_cli import __version__


DELEGATES = {
    "sourceosctl": "Install sourceos-devtools from SocioProphet/homebrew-prophet or run it from SourceOS-Linux/sourceos-devtools.",
    "agent-term": "Install agent-term from SocioProphet/homebrew-prophet or run it from SourceOS-Linux/agent-term.",
}


def _delegate(binary: str, args: Sequence[str]) -> int:
    resolved = shutil.which(binary)
    if not resolved:
        print(f"error: required delegate not found: {binary}", file=sys.stderr)
        print(DELEGATES.get(binary, "Install the missing delegate and retry."), file=sys.stderr)
        return 127
    completed = subprocess.run([resolved, *args], check=False)
    return int(completed.returncode)


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        prog="prophet",
        description="Prophet facade command surface for SocioProphet and SourceOS tools.",
    )
    parser.add_argument("--version", action="version", version=f"prophet {__version__}")

    sub = parser.add_subparsers(dest="command", required=True)

    sourceos = sub.add_parser("sourceos", help="Delegate SourceOS local workflows")
    sourceos_sub = sourceos.add_subparsers(dest="sourceos_command", required=True)

    agent_machine = sourceos_sub.add_parser("agent-machine", help="Delegate SourceOS Agent Machine commands")
    agent_machine.add_argument("args", nargs=argparse.REMAINDER, help="Arguments passed to sourceosctl agent-machine")

    office = sourceos_sub.add_parser("office", help="Delegate SourceOS Office Plane commands")
    office.add_argument("args", nargs=argparse.REMAINDER, help="Arguments passed to sourceosctl office")

    agent_term = sourceos_sub.add_parser("agent-term", help="Delegate AgentTerm commands")
    agent_term.add_argument("args", nargs=argparse.REMAINDER, help="Arguments passed to agent-term")

    return parser


def main(argv: list[str] | None = None) -> int:
    parser = build_parser()
    args = parser.parse_args(argv)

    if args.command == "sourceos":
        if args.sourceos_command == "agent-machine":
            return _delegate("sourceosctl", ["agent-machine", *args.args])
        if args.sourceos_command == "office":
            return _delegate("sourceosctl", ["office", *args.args])
        if args.sourceos_command == "agent-term":
            return _delegate("agent-term", list(args.args))

    parser.error("unknown command")
    return 2


if __name__ == "__main__":
    raise SystemExit(main(sys.argv[1:]))
