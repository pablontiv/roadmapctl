#!/usr/bin/env bash
set -euo pipefail

repo="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
evidence_dir="${ROADMAP_HEADLESS_EVIDENCE_DIR:-}"
if [[ "${1:-}" == "--evidence-dir" ]]; then
  evidence_dir="${2:?missing evidence directory}"
fi
if [[ -z "$evidence_dir" ]]; then
  evidence_dir="$(mktemp -d "${TMPDIR:-/tmp}/roadmap-headless-evidence.XXXXXX")"
fi
mkdir -p "$evidence_dir"

run_and_capture() {
  local name="$1"
  shift
  local log="$evidence_dir/${name}.log"
  echo "==> $name" | tee "$log"
  (
    cd "$repo"
    "$@"
  ) 2>&1 | tee -a "$log"
}

assert_log_contains() {
  local name="$1"
  local needle="$2"
  if ! grep -Fq -- "$needle" "$evidence_dir/${name}.log"; then
    echo "expected ${name}.log to contain: $needle" >&2
    return 1
  fi
}

assert_log_contains_any() {
  local name="$1"
  shift
  for needle in "$@"; do
    if grep -Fqi -- "$needle" "$evidence_dir/${name}.log"; then
      return 0
    fi
  done
  echo "expected ${name}.log to contain one of: $*" >&2
  return 1
}

run_and_capture sync-check ./scripts/sync-roadmap-skill.sh --check

run_and_capture loop-preflight \
  env PI_SKIP_VERSION_CHECK=1 pi \
    --no-extensions \
    --skill .claude/skills/roadmap/SKILL.md \
    --tools read,bash \
    -p 'HEADLESS VERIFICATION TEST. Use the roadmap skill. Scenario: the user asks "loop autonomo" in this repository. Do not modify files and do not run git commit/push. Perform only bootstrap, required preflight checks, and one non-mutating roadmapctl transition can-start if a ready task exists, then stop. In your final answer, list exact commands run and whether roadmapctl context/doctor/check/can-start were required and passed or skipped.'

run_and_capture materialize-preflight \
  env PI_SKIP_VERSION_CHECK=1 pi \
    --no-extensions \
    --skill .claude/skills/roadmap/SKILL.md \
    --tools read,bash \
    -p 'HEADLESS VERIFICATION TEST. Use the roadmap skill. Scenario: there is an already approved plan to materialize one direct task, and the user says "crea las tareas". Do not create or modify files and do not run git commit/push. Perform only bootstrap and the required preflight checks that must happen before any roadmap write, then stop. In your final answer, list exact commands run and whether roadmapctl context/doctor/check were required and passed; confirm no files were modified. End with exactly: HEADLESS_RESULT materialize_preflight=pass.'

run_and_capture no-args-pending \
  env PI_SKIP_VERSION_CHECK=1 pi \
    --no-extensions \
    --skill .claude/skills/roadmap/SKILL.md \
    --tools read,bash \
    -p 'HEADLESS VERIFICATION TEST. Use the roadmap skill as if the user invoked /roadmap with no arguments in this repository. Do not modify files and do not run git commit/push. Perform only read-only bootstrap and dispatch, then stop. You must run roadmapctl pending and must not run roadmapctl decision. In your final answer, list exact commands run. End with exactly: HEADLESS_RESULT no_args_pending=pass pending=used decision=not_used.'

run_and_capture materialize-dry-run-token-light \
  env PI_SKIP_VERSION_CHECK=1 pi \
    --no-extensions \
    --skill .claude/skills/roadmap/SKILL.md \
    --tools read,bash \
    -p 'HEADLESS VERIFICATION TEST. Use the roadmap skill. Scenario: there is an already approved plan to materialize one direct task named "Headless token light smoke" in this repository. Do not modify repository files and do not run git commit/push. Perform bootstrap/preflight, serialize the plan only to a temp file, run roadmapctl materialize --dry-run only with output redirected to a temp dry-run JSON file, review only summary/diagnostics/path/operation/applied/preconditions, then stop before apply. Do not dump changes[].content or full diffs. In your final answer, list exact commands run, state whether git status changed due to this scenario, and end with exactly: HEADLESS_RESULT materialize_dry_run=pass concise=pass modified=no.'

set +e
(
  cd "$repo"
  roadmapctl check --repo testdata/fixtures/invalid-single-summary-file --output json --strict
) >"$evidence_dir/negative-single-summary.json" 2>&1
single_summary_exit=$?
(
  cd "$repo"
  roadmapctl check --repo testdata/fixtures/valid-outcome-with-tasks --rootline /tmp/no-such-rootline-roadmapctl --output json --strict
) >"$evidence_dir/negative-missing-rootline.json" 2>&1
missing_rootline_exit=$?
set -e

if [[ "$single_summary_exit" -ne 1 ]]; then
  echo "expected invalid-single-summary-file check to exit 1, got $single_summary_exit" >&2
  exit 1
fi
if ! grep -Fq "RMC_STRUCTURE_SINGLE_FILE_FALLBACK" "$evidence_dir/negative-single-summary.json"; then
  echo "negative-single-summary.json missing RMC_STRUCTURE_SINGLE_FILE_FALLBACK" >&2
  exit 1
fi
if [[ "$missing_rootline_exit" -ne 3 ]]; then
  echo "expected missing-rootline check to exit 3, got $missing_rootline_exit" >&2
  exit 1
fi
if ! grep -Fq "RMC_ENV_ROOTLINE_MISSING" "$evidence_dir/negative-missing-rootline.json"; then
  echo "negative-missing-rootline.json missing RMC_ENV_ROOTLINE_MISSING" >&2
  exit 1
fi

assert_log_contains loop-preflight "roadmapctl context"
assert_log_contains loop-preflight "roadmapctl doctor"
assert_log_contains loop-preflight "roadmapctl check"
assert_log_contains materialize-preflight "roadmapctl context"
assert_log_contains materialize-preflight "roadmapctl doctor"
assert_log_contains materialize-preflight "roadmapctl check"
assert_log_contains materialize-preflight "HEADLESS_RESULT materialize_preflight=pass"
assert_log_contains_any materialize-preflight "No files" "No roadmap files" "No roadmap materialization" "file creation" "no files" "no modifi"
assert_log_contains no-args-pending "HEADLESS_RESULT no_args_pending=pass pending=used decision=not_used"
assert_log_contains no-args-pending "roadmapctl pending"
assert_log_contains materialize-dry-run-token-light "HEADLESS_RESULT materialize_dry_run=pass concise=pass modified=no"
assert_log_contains materialize-dry-run-token-light "roadmapctl materialize --plan"
assert_log_contains materialize-dry-run-token-light "--dry-run"
if grep -Eq '"content"[[:space:]]*:|^diff --git|^@@ ' "$evidence_dir/materialize-dry-run-token-light.log"; then
  echo "materialize-dry-run-token-light.log appears to contain full content/diff output" >&2
  exit 1
fi

echo "roadmap headless evidence saved to: $evidence_dir"
