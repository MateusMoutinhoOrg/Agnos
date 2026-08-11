#!/usr/bin/env bash
# mac86.sh — build the CLI for macOS Intel (amd64) into release/mac86.bin.
#
# Run it from anywhere:
#   bash ./scripts/mac86.sh
set -euo pipefail

# The repository root, resolved from this script, so the build never depends on
# the directory it was called from.
root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
mkdir -p "$root/release"

# The Go runtime cross-compiles on its own: GOOS/GOARCH is the whole toolchain.
# CGO stays off so no macOS SDK is needed on the building machine.
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 \
	go build -C "$root" -o release/mac86.bin ./cmd/main

echo "built release/mac86.bin"
