#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Usage: scripts/sync-roadmap-skill.sh [--install|--check] [--skill NAME|--all]

Synchronize canonical skills from this repository to the user scope.

Modes:
  --install  Copy skill(s) to ~/.claude/skills/. (default)
  --check    Verify source and installed skill(s) match without modifying files.

Options:
  --skill NAME  Skill directory name to sync (default: roadmap).
  --all         Sync all skill directories under .claude/skills/ that contain SKILL.md.

The script only reads .claude/skills/ and only writes to ~/.claude/skills/.
It does not touch other user-scope skills.
USAGE
}

MODE="install"
SKILL_NAME="roadmap"
ALL_SKILLS=false
while [ "$#" -gt 0 ]; do
  case "${1:-}" in
    --install) MODE="install"; shift ;;
    --check)   MODE="check";   shift ;;
    --all)     ALL_SKILLS=true; shift ;;
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
SKILLS_DEST="$HOME/.claude/skills"
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

sync_one_skill() {
  local skill_name="$1"
  local skill_src="$REPO_ROOT/.claude/skills/$skill_name"
  local skill_dest="$SKILLS_DEST/$skill_name"

  if [ ! -d "$skill_src" ]; then
    echo "roadmapctl: $skill_name skill source not found: $skill_src" >&2
    return 1
  fi
  if [ ! -f "$skill_src/SKILL.md" ]; then
    echo "roadmapctl: $skill_name skill source missing SKILL.md: $skill_src" >&2
    return 1
  fi

  check_one() {
    if [ ! -d "$skill_dest" ]; then
      echo "roadmapctl: installed $skill_name skill not found: $skill_dest" >&2
      return 1
    fi
    diff -qr "$skill_src" "$skill_dest" >/dev/null
  }

  case "$MODE" in
    check)
      if check_one; then
        echo "roadmapctl: $skill_name skill source and installed copy match"
        echo "roadmapctl: source: $skill_src"
        echo "roadmapctl: installed: $skill_dest"
      else
        echo "roadmapctl: $skill_name skill source and installed copy differ" >&2
        echo "roadmapctl: source: $skill_src" >&2
        echo "roadmapctl: installed: $skill_dest" >&2
        diff -qr "$skill_src" "$skill_dest" >&2 || true
        return 1
      fi
      ;;
    install)
      mkdir -p "$SKILLS_DEST"
      TMP_DIR="$(mktemp -d "$SKILLS_DEST/.${skill_name}-sync.XXXXXX")"
      cp -R "$skill_src"/. "$TMP_DIR"/

      if [ -e "$skill_dest" ]; then
        BACKUP_DIR="$(mktemp -d "$SKILLS_DEST/.${skill_name}-backup.XXXXXX")"
        mv "$skill_dest" "$BACKUP_DIR/$skill_name"
      fi

      mv "$TMP_DIR" "$skill_dest"
      TMP_DIR=""

      if ! check_one; then
        echo "roadmapctl: installed $skill_name skill did not match source; restoring previous copy" >&2
        rm -rf "$skill_dest"
        if [ -n "$BACKUP_DIR" ] && [ -d "$BACKUP_DIR/$skill_name" ]; then
          mv "$BACKUP_DIR/$skill_name" "$skill_dest"
        fi
        return 1
      fi

      if [ -n "$BACKUP_DIR" ] && [ -d "$BACKUP_DIR" ]; then
        rm -rf "$BACKUP_DIR"
        BACKUP_DIR=""
      fi

      echo "roadmapctl: installed $skill_name skill"
      echo "roadmapctl: source: $skill_src"
      echo "roadmapctl: installed: $skill_dest"
      ;;
  esac
}

if [ "$ALL_SKILLS" = true ]; then
  failed=0
  for skill_dir in "$REPO_ROOT/.claude/skills"/*/; do
    [ -f "${skill_dir}SKILL.md" ] || continue
    sync_one_skill "$(basename "$skill_dir")" || failed=1
  done
  exit "$failed"
else
  sync_one_skill "$SKILL_NAME"
fi
