#!/usr/bin/env bash
set -euo pipefail

threshold="${COVERAGE_THRESHOLD:-85.0}"
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
