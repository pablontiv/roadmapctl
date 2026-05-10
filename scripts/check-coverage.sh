#!/usr/bin/env bash
set -euo pipefail

coverage_threshold_from_toml() {
    local config_path="$1"
    [[ -f "$config_path" ]] || return 0
    awk -F= '
        /^[[:space:]]*required_code_coverage[[:space:]]*=/ {
            value = $2
            sub(/#.*/, "", value)
            gsub(/[[:space:]]/, "", value)
            print value
            exit
        }
    ' "$config_path"
}

if [[ -n "${COVERAGE_THRESHOLD:-}" ]]; then
    threshold="$COVERAGE_THRESHOLD"
else
    roadmap_root="${ROADMAP_ROOT:-docs/roadmap}"
    threshold="$(coverage_threshold_from_toml "$roadmap_root/.roadmapctl.toml")"
    threshold="${threshold:-85.0}"
fi
profile="${COVERAGE_PROFILE:-${TMPDIR:-/tmp}/roadmapctl.cover}"

go test ./... -coverprofile="$profile"
total_line="$(go tool cover -func="$profile" | awk '/^total:/ {print $3}')"
total="${total_line%%%}"

python3 - "$total" "$threshold" <<'PY'
import sys
actual = float(sys.argv[1])
threshold = float(sys.argv[2])
if actual < threshold:
    print(f"coverage {actual:.1f}% is below required {threshold:.1f}%", file=sys.stderr)
    sys.exit(1)
print(f"coverage {actual:.1f}% meets required {threshold:.1f}%")
PY
