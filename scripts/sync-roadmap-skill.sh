#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Usage: scripts/sync-roadmap-skill.sh [--install|--check]

Synchronize the canonical roadmap skill from this repository to the user scope.

Modes:
  --install  Copy .claude/skills/roadmap to ~/.claude/skills/roadmap.
  --check    Verify source and installed skill match without modifying files.

The script only reads .claude/skills/roadmap and only writes the roadmap skill
folder in ~/.claude/skills/roadmap. It does not touch other user-scope skills.
USAGE
}

MODE="install"
if [ "${1:-}" = "--install" ]; then
  MODE="install"
elif [ "${1:-}" = "--check" ]; then
  MODE="check"
elif [ "${1:-}" = "--help" ] || [ "${1:-}" = "-h" ]; then
  usage
  exit 0
elif [ "$#" -gt 0 ]; then
  echo "roadmapctl: unknown argument: $1" >&2
  usage >&2
  exit 2
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

SKILL_SRC="$REPO_ROOT/.claude/skills/roadmap"
SKILLS_DEST="$HOME/.claude/skills"
SKILL_DEST="$SKILLS_DEST/roadmap"
TMP_DIR=""
BACKUP_DIR=""

cleanup() {
  if [ -n "$TMP_DIR" ] && [ -d "$TMP_DIR" ]; then
    rm -rf "$TMP_DIR"
  fi
  if [ -n "$BACKUP_DIR" ] && [ -d "$BACKUP_DIR" ]; then
    rm -rf "$BACKUP_DIR"
  fi
}
trap cleanup EXIT

require_source() {
  if [ ! -d "$SKILL_SRC" ]; then
    echo "roadmapctl: roadmap skill source not found: $SKILL_SRC" >&2
    exit 1
  fi
  if [ ! -f "$SKILL_SRC/SKILL.md" ]; then
    echo "roadmapctl: roadmap skill source missing SKILL.md: $SKILL_SRC" >&2
    exit 1
  fi
}

check_sync() {
  if [ ! -d "$SKILL_DEST" ]; then
    echo "roadmapctl: installed roadmap skill not found: $SKILL_DEST" >&2
    return 1
  fi
  diff -qr "$SKILL_SRC" "$SKILL_DEST" >/dev/null
}

require_source

case "$MODE" in
  check)
    if check_sync; then
      echo "roadmapctl: roadmap skill source and installed copy match"
      echo "roadmapctl: source: $SKILL_SRC"
      echo "roadmapctl: installed: $SKILL_DEST"
    else
      echo "roadmapctl: roadmap skill source and installed copy differ" >&2
      echo "roadmapctl: source: $SKILL_SRC" >&2
      echo "roadmapctl: installed: $SKILL_DEST" >&2
      diff -qr "$SKILL_SRC" "$SKILL_DEST" >&2 || true
      exit 1
    fi
    ;;
  install)
    mkdir -p "$SKILLS_DEST"
    TMP_DIR="$(mktemp -d "$SKILLS_DEST/.roadmap-sync.XXXXXX")"
    cp -R "$SKILL_SRC"/. "$TMP_DIR"/

    if [ -e "$SKILL_DEST" ]; then
      BACKUP_DIR="$(mktemp -d "$SKILLS_DEST/.roadmap-backup.XXXXXX")"
      mv "$SKILL_DEST" "$BACKUP_DIR/roadmap"
    fi

    mv "$TMP_DIR" "$SKILL_DEST"
    TMP_DIR=""

    if ! check_sync; then
      echo "roadmapctl: installed roadmap skill did not match source; restoring previous copy" >&2
      rm -rf "$SKILL_DEST"
      if [ -n "$BACKUP_DIR" ] && [ -d "$BACKUP_DIR/roadmap" ]; then
        mv "$BACKUP_DIR/roadmap" "$SKILL_DEST"
      fi
      exit 1
    fi

    if [ -n "$BACKUP_DIR" ] && [ -d "$BACKUP_DIR" ]; then
      rm -rf "$BACKUP_DIR"
      BACKUP_DIR=""
    fi

    echo "roadmapctl: installed roadmap skill"
    echo "roadmapctl: source: $SKILL_SRC"
    echo "roadmapctl: installed: $SKILL_DEST"
    ;;
esac
