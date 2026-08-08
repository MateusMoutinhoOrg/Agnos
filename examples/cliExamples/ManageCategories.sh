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

go build -o "$workdir/agnos" ./cmd/main
export AGNOS_DATA="$workdir/data"
agnos() { "$workdir/agnos" "$@"; }

echo "== create the categories"
agnos category add groceries
agnos category add salary
agnos category add rent

echo
echo "== creating a category is idempotent: the stored one comes back"
agnos category add groceries

echo
echo "== every category, with its balance and how many transactions it holds"
agnos category list

echo
echo "== remove one, and its transactions with it"
agnos category remove rent
agnos category list
