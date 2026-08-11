#!/usr/bin/env bash
# linuxarm64.sh — build the CLI for Linux arm64 into release/linuxarm64.out.
#
# Run it from anywhere:
#   bash ./scripts/linuxarm64.sh
set -euo pipefail

# The repository root, resolved from this script, so the build never depends on
# the directory it was called from.
root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
mkdir -p "$root/release"

# The Go runtime cross-compiles on its own: GOOS/GOARCH is the whole toolchain.
# CGO stays off so the binary carries no libc of the building machine.
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
	go build -C "$root" -o release/linuxarm64.out ./cmd/main

echo "built release/linuxarm64.out"
