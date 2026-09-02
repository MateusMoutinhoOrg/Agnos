#!/usr/bin/env bash
# macarm64.sh — build the CLI for macOS Apple Silicon (arm64) into
# release/macarm64.bin.
#
# Run it from anywhere:
#   bash ./scripts/macarm64.sh
set -euo pipefail

# The repository root, resolved from this script, so the build never depends on
# the directory it was called from.
root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
mkdir -p "$root/release"

# The Go runtime cross-compiles on its own: GOOS/GOARCH is the whole toolchain.
# CGO stays off so no macOS SDK is needed on the building machine.
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 \
	go build -C "$root" -o release/macarm64.bin ./cmd/main

echo "built release/macarm64.bin"
