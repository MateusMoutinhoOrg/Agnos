#!/usr/bin/env bash
# ManageCategories.sh — set up a budget: create the categories transactions are
# tracked under, list them back, and drop one that was a mistake.
#
# Run it from the project root:
#   bash ./examples/cliExamples/ManageCategories.sh
set -euo pipefail

# Build the CLI into a scratch directory and point it at a budget of its own,
# so the example never touches the records in your home directory.
workdir="$(mktemp -d)"
trap 'rm -rf "$workdir"' EXIT

go build -o "$workdir/agnos-cli" ./cmd/main
export AGNOS_DATA="$workdir/data"
agnos-cli() { "$workdir/agnos-cli" "$@"; }

echo "== create the categories"
agnos-cli category add groceries
agnos-cli category add salary
agnos-cli category add rent

echo
echo "== creating a category is idempotent: the stored one comes back"
agnos-cli category add groceries

echo
echo "== every category, with its balance and how many transactions it holds"
agnos-cli category list

echo
echo "== remove one, and its transactions with it"
agnos-cli category remove rent
agnos-cli category list
