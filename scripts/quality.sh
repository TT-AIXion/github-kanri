#!/bin/bash
set -euo pipefail

min_coverage=100

if [ "${1:-}" = "--min" ] && [ -n "${2:-}" ]; then
  min_coverage="$2"
fi

coverfile=$(mktemp)
cleanup() {
  rm -f "$coverfile"
}
trap cleanup EXIT

go vet ./...

go test ./... -coverprofile="$coverfile"

total=$(go tool cover -func="$coverfile" | awk '/^total:/{print $3}')
total=${total%%%}

awk -v total="$total" -v min="$min_coverage" 'BEGIN { if (total+0 < min+0) exit 1 }'

echo "OK coverage ${total}% (min ${min_coverage}%)"
