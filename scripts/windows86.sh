#!/usr/bin/env bash
# windows86.sh — build the CLI for Windows amd64 into release/windows86.exe.
#
# Run it from anywhere:
#   bash ./scripts/windows86.sh
set -euo pipefail

# The repository root, resolved from this script, so the build never depends on
# the directory it was called from.
root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
mkdir -p "$root/release"

# The Go runtime cross-compiles on its own: GOOS/GOARCH is the whole toolchain.
# CGO stays off so the binary carries no libc of the building machine.
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 \
	go build -C "$root" -o release/windows86.exe ./cmd/main

echo "built release/windows86.exe"
