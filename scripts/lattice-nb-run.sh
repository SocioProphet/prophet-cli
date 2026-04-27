#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="${PROPHET_PLATFORM_DIR:-$HOME/dev/prophet-platform}"
STUDIO_SRC="$ROOT_DIR/apps/lattice-studio/src"

if [[ ! -d "$STUDIO_SRC" ]]; then
  echo "prophet-cli: missing Lattice Studio source at $STUDIO_SRC" >&2
  echo "Set PROPHET_PLATFORM_DIR or clone SocioProphet/prophet-platform under ~/dev." >&2
  exit 2
fi

export PYTHONPATH="$STUDIO_SRC${PYTHONPATH:+:$PYTHONPATH}"
exec python3 -m lattice_studio.notebook_launch_cli "$@"
