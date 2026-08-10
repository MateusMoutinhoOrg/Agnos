#!/usr/bin/env bash
# ScriptTheCli.sh — script the CLI: drive it from a shell, read its exit codes,
# and feed a listing into other tools.
#
# Run it from the project root:
#   bash ./examples/cliExamples/ScriptTheCli.sh
set -uo pipefail   # no -e here: this example inspects failing exit codes

# Build the CLI into a scratch directory and point it at a budget of its own,
# so the example never touches the records in your home directory.
workdir="$(mktemp -d)"
trap 'rm -rf "$workdir"' EXIT

go build -o "$workdir/agnos-cli" ./cmd/main
export AGNOS_DATA="$workdir/data"
agnos-cli() { "$workdir/agnos-cli" "$@"; }

echo "== --quiet keeps a mutating command silent, so a loop stays readable"
agnos-cli --quiet category add groceries
for entry in "monday market:21.40" "tuesday bakery:6.25" "friday market:33.10"; do
    agnos-cli --quiet spend groceries "${entry%%:*}" "${entry##*:}"
done
agnos-cli transactions groceries

echo
echo "== exit codes: 0 ran, 1 the command line was wrong, 2 the command failed"
agnos-cli balance groceries > /dev/null
echo "balance groceries          -> $?"

agnos-cli spend groceries "bad amount" not-a-number > /dev/null
echo "spend with a bad amount    -> $?"

agnos-cli balance nosuchcategory > /dev/null
echo "balance of a missing one   -> $?"

echo
echo "== a listing is one record per line, so the usual text tools apply"
echo "transactions recorded: $(agnos-cli transactions groceries | wc -l | tr -d ' ')"
echo "most recent one:       $(agnos-cli transactions groceries | tail -1)"
echo "everything at a market:"
agnos-cli transactions groceries | grep market
