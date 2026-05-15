#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Usage: scripts/sync-roadmap-skill.sh [--install|--check] [--skill NAME]

Synchronize a canonical skill from this repository to the user scope.

Modes:
  --install  Copy .claude/skills/NAME to ~/.claude/skills/NAME. (default)
  --check    Verify source and installed skill match without modifying files.

Options:
  --skill NAME  Skill directory name to sync (default: roadmap).

The script only reads .claude/skills/NAME and only writes that skill's
folder in ~/.claude/skills/NAME. It does not touch other user-scope skills.
USAGE
}

MODE="install"
SKILL_NAME="roadmap"
while [ "$#" -gt 0 ]; do
  case "${1:-}" in
    --install) MODE="install"; shift ;;
    --check)   MODE="check";   shift ;;
    --skill)
      if [ -z "${2:-}" ]; then
        echo "roadmapctl: --skill requires a name" >&2; usage >&2; exit 2
      fi
      SKILL_NAME="$2"; shift 2 ;;
    --help|-h) usage; exit 0 ;;
    *) echo "roadmapctl: unknown argument: $1" >&2; usage >&2; exit 2 ;;
  esac
done

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

SKILL_SRC="$REPO_ROOT/.claude/skills/$SKILL_NAME"
SKILLS_DEST="$HOME/.claude/skills"
SKILL_DEST="$SKILLS_DEST/$SKILL_NAME"
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
    echo "roadmapctl: $SKILL_NAME skill source not found: $SKILL_SRC" >&2
    exit 1
  fi
  if [ ! -f "$SKILL_SRC/SKILL.md" ]; then
    echo "roadmapctl: $SKILL_NAME skill source missing SKILL.md: $SKILL_SRC" >&2
    exit 1
  fi
}

check_sync() {
  if [ ! -d "$SKILL_DEST" ]; then
    echo "roadmapctl: installed $SKILL_NAME skill not found: $SKILL_DEST" >&2
    return 1
  fi
  diff -qr "$SKILL_SRC" "$SKILL_DEST" >/dev/null
}

require_source

case "$MODE" in
  check)
    if check_sync; then
      echo "roadmapctl: $SKILL_NAME skill source and installed copy match"
      echo "roadmapctl: source: $SKILL_SRC"
      echo "roadmapctl: installed: $SKILL_DEST"
    else
      echo "roadmapctl: $SKILL_NAME skill source and installed copy differ" >&2
      echo "roadmapctl: source: $SKILL_SRC" >&2
      echo "roadmapctl: installed: $SKILL_DEST" >&2
      diff -qr "$SKILL_SRC" "$SKILL_DEST" >&2 || true
      exit 1
    fi
    ;;
  install)
    mkdir -p "$SKILLS_DEST"
    TMP_DIR="$(mktemp -d "$SKILLS_DEST/.${SKILL_NAME}-sync.XXXXXX")"
    cp -R "$SKILL_SRC"/. "$TMP_DIR"/

    if [ -e "$SKILL_DEST" ]; then
      BACKUP_DIR="$(mktemp -d "$SKILLS_DEST/.${SKILL_NAME}-backup.XXXXXX")"
      mv "$SKILL_DEST" "$BACKUP_DIR/$SKILL_NAME"
    fi

    mv "$TMP_DIR" "$SKILL_DEST"
    TMP_DIR=""

    if ! check_sync; then
      echo "roadmapctl: installed $SKILL_NAME skill did not match source; restoring previous copy" >&2
      rm -rf "$SKILL_DEST"
      if [ -n "$BACKUP_DIR" ] && [ -d "$BACKUP_DIR/$SKILL_NAME" ]; then
        mv "$BACKUP_DIR/$SKILL_NAME" "$SKILL_DEST"
      fi
      exit 1
    fi

    if [ -n "$BACKUP_DIR" ] && [ -d "$BACKUP_DIR" ]; then
      rm -rf "$BACKUP_DIR"
      BACKUP_DIR=""
    fi

    echo "roadmapctl: installed $SKILL_NAME skill"
    echo "roadmapctl: source: $SKILL_SRC"
    echo "roadmapctl: installed: $SKILL_DEST"
    ;;
esac
