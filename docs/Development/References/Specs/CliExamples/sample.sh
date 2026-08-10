#!/usr/bin/env bash
# This file is an illustrative sample, not part of the build.
#
# <Name>.sh — one sentence naming the goal this script demonstrates.
#
# Run it from the project root:
#   bash ./examples/cliExamples/<Name>.sh
set -euo pipefail

# Build the CLI into a scratch directory and point it at a budget of its own,
# so the example never touches the records in your home directory.
workdir="$(mktemp -d)"
trap 'rm -rf "$workdir"' EXIT

go build -o "$workdir/agnos-cli" ./cmd/main
export AGNOS_DATA="$workdir/data"
agnos-cli() { "$workdir/agnos-cli" "$@"; }

echo "== what the commands below show"
agnos-cli category add examplecategory
agnos-cli spend examplecategory "an example transaction" 12.34

echo
echo "== what the next commands show"
agnos-cli transactions examplecategory
agnos-cli balance examplecategory
