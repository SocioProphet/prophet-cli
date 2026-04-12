#!/usr/bin/env bash
set -euo pipefail

MODE="${1:-compose}"
PLATFORM_ROOT="${PROPHET_PLATFORM_ROOT:-../prophet-platform}"

if [[ ! -d "${PLATFORM_ROOT}" ]]; then
  echo "prophet-platform repo not found at ${PLATFORM_ROOT}" >&2
  echo "Set PROPHET_PLATFORM_ROOT or use a sibling checkout layout." >&2
  exit 1
fi

cd "${PLATFORM_ROOT}"

case "${MODE}" in
  compose)
    make compose-up
    ;;
  kind)
    make kind-up
    ;;
  k3d)
    make k3d-up
    ;;
  *)
    echo "unsupported mode: ${MODE}" >&2
    echo "supported modes: compose | kind | k3d" >&2
    exit 1
    ;;
esac
