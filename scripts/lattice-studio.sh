#!/usr/bin/env bash
set -euo pipefail

# Prophet CLI facade for the Lattice Studio command surface.
#
# Expected local checkout layout:
#   ~/dev/prophet-cli
#   ~/dev/prophet-platform
#
# This script delegates to prophet-platform/apps/lattice-studio without
# duplicating implementation logic in the CLI facade repo.

ROOT_DIR="${PROPHET_PLATFORM_DIR:-$HOME/dev/prophet-platform}"
STUDIO_SRC="$ROOT_DIR/apps/lattice-studio/src"

if [[ ! -d "$STUDIO_SRC" ]]; then
  echo "prophet-cli: missing Lattice Studio source at $STUDIO_SRC" >&2
  echo "Set PROPHET_PLATFORM_DIR or clone SocioProphet/prophet-platform under ~/dev." >&2
  exit 2
fi

export PYTHONPATH="$STUDIO_SRC${PYTHONPATH:+:$PYTHONPATH}"
exec python3 -m lattice_studio.cli "$@"
