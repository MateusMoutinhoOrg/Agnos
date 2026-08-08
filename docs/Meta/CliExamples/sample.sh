#!/usr/bin/env bash
# This file is an illustrative sample, not part of the build.
#
# <Name>.sh — one sentence naming the goal this script demonstrates.
#
# Run it from the project root:
#   bash ./cliExamples/<Name>.sh
set -euo pipefail

# Build the CLI into a scratch directory and point it at a budget of its own,
# so the example never touches the records in your home directory.
workdir="$(mktemp -d)"
trap 'rm -rf "$workdir"' EXIT

go build -o "$workdir/agnos" ./cmd/main
export AGNOS_DATA="$workdir/data"
agnos() { "$workdir/agnos" "$@"; }

echo "== what the commands below show"
agnos category add examplecategory
agnos spend examplecategory "an example transaction" 12.34

echo
echo "== what the next commands show"
agnos transactions examplecategory
agnos balance examplecategory
