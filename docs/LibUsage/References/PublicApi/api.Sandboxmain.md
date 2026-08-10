# `api.Lib.Sandboxmain`

**Type:** Field

## Signature

```go
Sandboxmain func(args []string) int
```

## Description

The command-line interface itself, run **inside the sandbox**. It reads the actions and flags of the command line through the injected Verb parser ([`deps.Deps.VerbLib`](/docs/LibUsage/References/PublicApi/verbdeps.Lib.md)), calls the library functions of [`api.Lib`](/docs/LibUsage/References/PublicApi/api.Lib.md), prints every result and error through `deps.Deps.Printf`, and returns the process exit code.

Everything the interface does is a library call: the installed binary in `cmd/main` holds no command, no flag, and no output of its own — it wires an adapter into the library, hands this field the argument vector, and exits with what it returns. Every command, flag, and exit code the interface understands is listed in [Commands.md](/docs/CliUsage/References/Commands.md).

`args` must be the same argument vector the adapter wired `Deps.VerbLib` over: the parser owns the reading, and `args` is only checked here to detect an empty command line, which prints the usage screen. The standard adapter and `cmd/main` both take it from `os.Args[1:]`, so they agree by construction.

## Parameters

| Name | Type | Description |
| :--- | :--- | :--- |
| `args` | `[]string` | The argument vector, without the program name — the same one the adapter wired `Deps.VerbLib` over. |

## Returns

| Type | Description |
| :--- | :--- |
| `int` | The process exit code: `api.ExitOk` (0) when the command ran to completion, `api.ExitUsage` (1) when the command line itself was wrong, `api.ExitFailure` (2) when a well-formed command could not be carried out. |

## Examples

```go
package main

import (
	"os"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

func main() {
	// The adapter wires the Verb parser over the same vector handed below.
	l := agnoslib.New(agnosadapter.New("trackerdata"))

	os.Exit(l.Sandboxmain(os.Args[1:]))
}
```
