#!/usr/bin/env bash
set -euo pipefail

PLATFORM_ROOT="${PROPHET_PLATFORM_ROOT:-../prophet-platform}"

if [[ ! -d "${PLATFORM_ROOT}" ]]; then
  echo "prophet-platform repo not found at ${PLATFORM_ROOT}" >&2
  echo "Set PROPHET_PLATFORM_ROOT or use a sibling checkout layout." >&2
  exit 1
fi

cd "${PLATFORM_ROOT}"
make nuke
