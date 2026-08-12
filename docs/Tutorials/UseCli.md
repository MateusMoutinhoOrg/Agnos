# Track a Budget from the Terminal

## Description
Covers running the `agnos-cli` CLI end to end: creating categories, recording what comes in and what goes out, and reading the result back. Installing the binary is covered by [InstallCli.md](/docs/Tutorials/InstallCli.md); every command and flag is listed in [Commands.md](/docs/References/Commands.md).

### Rules
- The commands below assume `agnos-cli` is on your `PATH`. From a checkout, replace `agnos-cli` with `go run ./cmd/main` throughout.
- Amounts are always positive: the command — `spend` or `received` — is what carries the direction.
- A category has to exist before money can be tracked under it.

---

## Workflow
1. Create the categories the budget is split into:
   ```bash
   agnos-cli category add groceries
   agnos-cli category add salary
   ```
2. Record money entering the budget:
   ```bash
   agnos-cli received salary "august paycheck" 2500.00
   ```
3. Record money leaving it. Descriptions with spaces are one argument, so quote them:
   ```bash
   agnos-cli spend groceries "weekly shopping" 84.50
   agnos-cli spend groceries "corner store" 12.9
   ```
4. Read the transactions back — everything, or one category's:
   ```bash
   agnos-cli transactions
   agnos-cli transactions groceries
   ```
5. Read the balances. Spending counts down, receiving counts up:
   ```bash
   agnos-cli balance groceries    # -97.40
   agnos-cli balance              # 2402.60
   ```
6. List the categories with their running balances:
   ```bash
   agnos-cli category list
   ```
7. Drop a category you no longer track, and its transactions with it:
   ```bash
   agnos-cli category remove groceries
   ```
8. Add `--quiet` when the confirmation lines are noise — in a loop, for instance:
   ```bash
   agnos-cli --quiet category add books
   ```
9. Check what a command did from a script by reading its exit code: `0` ran, `1` the command line was wrong, `2` the command failed. See [Commands.md](/docs/References/Commands.md#exit-codes):
   ```bash
   agnos-cli balance nosuchcategory || echo "failed with $?"
   ```
10. Read a worked transcript of all of this in [examples/cliExamples/](/examples/cliExamples/), following [HandleCliExamples.md](/docs/Tutorials/HandleCliExamples.md).
