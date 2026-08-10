# CLI Commands

## Description
Every command, flag, and exit code of the `agnos` command-line interface. The interface itself is `api.Lib.Sandboxmain`, a field of the library like any other; the binary in [cmd/main](/cmd/main/) only wires an adapter into it and exits with what it returns. To install it, follow [InstallCli.md](/docs/CliUsage/Protocols/InstallCli.md); to walk through a first budget, follow [UseCli.md](/docs/CliUsage/Protocols/UseCli.md).

```bash
agnos <command> [arguments] [flags]
```

---

## Commands

| Command | Description |
|---------|-------------|
| `category add <name>` | Creates the category. Already-taken names are not an error: the stored category comes back instead. |
| `category list` | Prints every category, oldest first, with its balance and transaction count. |
| `category remove <name>` | Deletes the category and every transaction stored under it. |
| `spend <category> <description> <amount>` | Records money leaving the budget under an existing category. |
| `received <category> <description> <amount>` | Records money entering the budget under an existing category. |
| `transactions [category]` | Prints every transaction, or only the given category's when a name follows. |
| `balance [category]` | Prints the total balance, or the given category's when a name follows. |
| `help` | Prints the usage screen. |
| `version` | Prints the interface version. |

---

## Flags

| Flag | Description |
|------|-------------|
| `-h`, `--help` | Prints the usage screen and exits, whatever else is on the command line. |
| `-v`, `--version` | Prints the interface version and exits. |
| `-q`, `--quiet` | Suppresses the confirmation line a mutating command prints, leaving listings and errors. |

Flags are read by the injected [Verb](https://github.com/MateusMoutinhoOrg/Verb) parser before the command words are drained, so a flag may appear anywhere on the command line: `agnos spend groceries lunch 12.00 --quiet` and `agnos --quiet spend groceries lunch 12.00` do the same thing.

---

## Amounts

Amounts are written the way a person types money — `84.50`, `84.5`, and `84` are all accepted — and stored in the smallest currency unit, which is what the library's [`AddSpend`](/docs/LibUsage/References/PublicApi/api.AddSpend.md) and [`AddReceived`](/docs/LibUsage/References/PublicApi/api.AddReceived.md) take.

| Written | Stored | |
|---------|--------|--|
| `84.50` | `8450` | Two decimal places. |
| `84.5` | `8450` | One place is padded. |
| `84` | `8400` | No places at all. |
| `-84.50` | — | Rejected: the command chooses the direction, so an amount is always positive. |
| `84.505` | — | Rejected: more than two places. |

---

## Exit Codes

| Code | Constant | Meaning |
|------|----------|---------|
| `0` | `api.ExitOk` | The command ran to completion. |
| `1` | `api.ExitUsage` | The command line was wrong — unknown command, missing operand, unparsable amount — and the usage screen was printed. |
| `2` | `api.ExitFailure` | The command was well-formed but could not be carried out: a record was missing or could not be written. |

---

## Environment

| Variable | Description |
|----------|-------------|
| `AGNOS_DATA` | Directory the records are kept in. Defaults to `.agnos` in the user's home directory. Read by [cmd/main](/cmd/main/), never by the sandbox — where state lives is an OS-bound choice, so it is made outside the library. |
