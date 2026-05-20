#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

BIN="${REPO_ROOT}/bin/roadmapctl-smoke"
mkdir -p "${REPO_ROOT}/bin"
go build -o "${BIN}" "${REPO_ROOT}/cmd/roadmapctl"

start_ns=$(date +%s%N)
ROADMAPCTL_NO_UPDATE=1 "${BIN}" --version >/dev/null
end_ns=$(date +%s%N)
elapsed_ms=$(( (end_ns - start_ns) / 1000000 ))

rm -f "${BIN}"

if [ "${elapsed_ms}" -gt 100 ]; then
    echo "FAIL: roadmapctl --version took ${elapsed_ms}ms (>100ms)" >&2
    exit 1
fi
echo "OK: roadmapctl --version completed in ${elapsed_ms}ms"
