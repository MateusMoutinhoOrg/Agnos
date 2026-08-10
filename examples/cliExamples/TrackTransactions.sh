#!/usr/bin/env bash
# TrackTransactions.sh — track money: record spend and received transactions, list
# them, and read the balance of one category and of the whole budget.
#
# Run it from the project root:
#   bash ./examples/cliExamples/TrackTransactions.sh
set -euo pipefail

# Build the CLI into a scratch directory and point it at a budget of its own,
# so the example never touches the records in your home directory.
workdir="$(mktemp -d)"
trap 'rm -rf "$workdir"' EXIT

go build -o "$workdir/agnos-cli" ./cmd/main
export AGNOS_DATA="$workdir/data"
agnos-cli() { "$workdir/agnos-cli" "$@"; }

echo "== a category has to exist before money can be tracked under it"
agnos-cli --quiet category add groceries
agnos-cli --quiet category add salary

echo
echo "== record what came in and what went out"
agnos-cli received salary "august paycheck" 2500.00
agnos-cli spend groceries "weekly shopping" 84.50
agnos-cli spend groceries "corner store" 12.9

echo
echo "== every transaction, oldest first, grouped by category"
agnos-cli transactions

echo
echo "== only the ones under groceries"
agnos-cli transactions groceries

echo
echo "== groceries balance, then the whole budget's"
agnos-cli balance groceries
agnos-cli balance
