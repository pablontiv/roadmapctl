#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Usage: scripts/install-user.sh

Install the user-scope roadmap skill and rebuild/install the roadmapctl binary.

Environment:
  ROADMAPCTL_BIN  Destination binary path (default: /usr/local/bin/roadmapctl)
USAGE
}

if [ "${1:-}" = "--help" ] || [ "${1:-}" = "-h" ]; then
  usage
  exit 0
elif [ "$#" -gt 0 ]; then
  echo "roadmapctl: unknown argument: $1" >&2
  usage >&2
  exit 2
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

"$REPO_ROOT/scripts/sync-roadmap-skill.sh" --install
"$REPO_ROOT/scripts/sync-roadmap-skill.sh" --install --skill retrospective

resolve_roadmapctl_bin() {
  if [ -n "${ROADMAPCTL_BIN:-}" ]; then
    echo "$ROADMAPCTL_BIN"
    return
  fi
  echo "/usr/local/bin/roadmapctl"
}

ROADMAPCTL_BIN="$(resolve_roadmapctl_bin)"
mkdir -p "$(dirname "$ROADMAPCTL_BIN")"
LEGACY_ROADMAPCTL_BIN="$HOME/.local/bin/roadmapctl"

install_roadmapctl_binary() {
  local tmp
  tmp="$(mktemp)"

  if ! (cd "$REPO_ROOT" && go build -o "$tmp" ./cmd/roadmapctl) 2>&1; then
    rm -f "$tmp"
    return 1
  fi

  if [ -w "$ROADMAPCTL_BIN" ] || { [ ! -e "$ROADMAPCTL_BIN" ] && [ -w "$(dirname "$ROADMAPCTL_BIN")" ]; }; then
    install -m 0755 "$tmp" "$ROADMAPCTL_BIN"
  elif command -v sudo >/dev/null 2>&1 && sudo -n true 2>/dev/null; then
    sudo install -m 0755 "$tmp" "$ROADMAPCTL_BIN"
  else
    echo "Warning: $ROADMAPCTL_BIN is not writable and passwordless sudo is unavailable" >&2
    rm -f "$tmp"
    return 1
  fi

  rm -f "$tmp"
}

echo "Rebuilding roadmapctl..."
if install_roadmapctl_binary; then
  if [ "$ROADMAPCTL_BIN" = "/usr/local/bin/roadmapctl" ] && [ -e "$LEGACY_ROADMAPCTL_BIN" ]; then
    rm -f "$LEGACY_ROADMAPCTL_BIN"
  fi
  echo "roadmapctl rebuilt: $ROADMAPCTL_BIN"
  if command -v roadmapctl >/dev/null 2>&1 && [ "$(command -v roadmapctl)" != "$ROADMAPCTL_BIN" ]; then
    echo "Warning: $(command -v roadmapctl) shadows $ROADMAPCTL_BIN in PATH" >&2
  fi
else
  echo "Warning: roadmapctl rebuild failed (continuing)" >&2
fi
