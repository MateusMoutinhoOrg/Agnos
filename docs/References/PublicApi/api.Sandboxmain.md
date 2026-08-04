# `api.Lib.Sandboxmain`

**Type:** Field

## Signature

```go
Sandboxmain func(args []string) int
```

## Description

The command-line interface, run inside the sandbox. It reads the commands and flags of the given argument vector through the injected [`deps.Deps.VerbLib`](/docs/References/PublicApi/verbdeps.Lib.md) parser, calls the other fields of [`api.Lib`](/docs/References/PublicApi/api.Lib.md) to do the work, prints every result and error through [`deps.Deps.Printf`](/docs/References/PublicApi/deps.Deps.md), and returns the process exit code. Every command and flag it understands is listed in [Cli.md](/docs/References/Cli.md).

Putting the whole program behind one library field is what keeps the interface portable: the binary in [cmd/main](/cmd/main/) does nothing but wire an adapter into the library, call this field, and `os.Exit` with what it returns, so the same interface runs anywhere a `deps.Deps` can be filled — a test harness with a captured writer and a fixed argument vector included.

`args` must be the same vector the adapter wired `Deps.VerbLib` over: the parser owns the reading, and `args` is what `Sandboxmain` checks for an empty command line. [`standard.New`](/docs/References/PublicApi/standard.New.md) and `cmd/main` both take it from `os.Args[1:]`, so they agree by construction.

## Parameters

| Name | Type | Description |
| :--- | :--- | :--- |
| `args` | `[]string` | The argument vector, without the program name — the same one the adapter wired the parser over. |

## Returns

| Type | Description |
| :--- | :--- |
| `int` | The process exit code: `api.ExitOk` (0), `api.ExitUsage` (1), or `api.ExitFailure` (2). |

## Examples

```go
package main

import (
	"os"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

func main() {
	// The whole CLI: wire an adapter into the lib, run the interface over
	// the command line, and exit with the code it returns.
	l := agnoslib.New(agnosadapter.New("trackerdata"))
	os.Exit(l.Sandboxmain(os.Args[1:]))
}
```
