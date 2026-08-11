#!/usr/bin/env bash
# windowsi32.sh — build the CLI for Windows 386 into release/windowsi32.exe.
#
# Run it from anywhere:
#   bash ./scripts/windowsi32.sh
set -euo pipefail

# The repository root, resolved from this script, so the build never depends on
# the directory it was called from.
root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
mkdir -p "$root/release"

# The Go runtime cross-compiles on its own: GOOS/GOARCH is the whole toolchain.
# CGO stays off so the binary carries no libc of the building machine.
CGO_ENABLED=0 GOOS=windows GOARCH=386 \
	go build -C "$root" -o release/windowsi32.exe ./cmd/main

echo "built release/windowsi32.exe"
