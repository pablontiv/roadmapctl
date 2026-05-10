#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$script_dir/.." && pwd)"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

fake_bin="$tmp/bin"
mkdir -p "$fake_bin"
cat >"$fake_bin/go" <<'FAKEGO'
#!/usr/bin/env bash
set -euo pipefail
if [[ "$1" == "test" ]]; then
  exit 0
fi
if [[ "$1" == "tool" && "$2" == "cover" ]]; then
  printf 'total:\t(statements)\t90.0%%\n'
  exit 0
fi
echo "unexpected go invocation: $*" >&2
exit 2
FAKEGO
chmod +x "$fake_bin/go"

run_check() {
  local workdir="$1"
  shift
  (cd "$workdir" && PATH="$fake_bin:$PATH" "$@" "$repo_root/scripts/check-coverage.sh")
}

assert_contains() {
  local text="$1"
  local want="$2"
  if [[ "$text" != *"$want"* ]]; then
    printf 'missing %q in output:\n%s\n' "$want" "$text" >&2
    exit 1
  fi
}

case "${1:-all}" in
  env-override|all)
    workdir="$tmp/env-override"
    mkdir -p "$workdir/docs/roadmap"
    printf 'required_code_coverage = 99.0\n' >"$workdir/docs/roadmap/.roadmapctl.toml"
    output="$(run_check "$workdir" env COVERAGE_THRESHOLD=88.0)"
    assert_contains "$output" 'required 88.0%'
    ;;&
  toml|all)
    workdir="$tmp/toml"
    mkdir -p "$workdir/docs/roadmap"
    printf 'required_code_coverage = 89.5\n' >"$workdir/docs/roadmap/.roadmapctl.toml"
    output="$(run_check "$workdir" env -u COVERAGE_THRESHOLD)"
    assert_contains "$output" 'required 89.5%'
    ;;&
  fallback|all)
    workdir="$tmp/fallback"
    mkdir -p "$workdir"
    output="$(run_check "$workdir" env -u COVERAGE_THRESHOLD)"
    assert_contains "$output" 'required 85.0%'
    ;;
  *)
    echo "usage: $0 [all|env-override|toml|fallback]" >&2
    exit 2
    ;;
esac
