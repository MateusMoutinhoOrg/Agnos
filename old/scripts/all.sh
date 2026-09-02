#!/usr/bin/env bash
# all.sh — build every target into release/, one after the other.
#
# Run it from anywhere:
#   bash ./scripts/all.sh
set -euo pipefail

# The repository root, resolved from this script, so the build never depends on
# the directory it was called from.
root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# One entry per target script beside this one. Adding a target means adding its
# script here.
targets=(
	linux86
	linuxarm64
	linuxi32
	windows86
	windowsi32
	mac86
	macarm64
)

for target in "${targets[@]}"; do
	bash "$root/scripts/$target.sh"
done

echo
echo "every target is in $root/release"
