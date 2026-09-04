# `api.Cli`

**Type:** Struct

## Definition

```go
type Cli struct {
	CliMain func(args []string) int
}

const (
	ExitOk      = 0
	ExitFailure = 1
	ExitUsage   = 2
)
```

## Description

The command-line interface of the sandbox as a single function. `CliMain` takes an argument vector — what `cmd/main` passes as `os.Args[1:]` — reads the verb, parses the remaining arguments into the matching command's typed `Entries`, rejects anything unknown or invalid, calls the command's `CommandHandler`, and returns the process exit code. There is no command registry to consult: the function is **generated** by `agnos build` from every `sandbox/internal/commands/<name>/entries.yaml` — see [CommandDispatch](/docs/CommandDispatch/doc.md). `sandbox/api/cli.go` itself belongs to the `cli` asset group and is regenerated too.

The exit-code constants are the same three everywhere and are mirrored inside the generated `cli` package, so that layer depends on no contract package.

## Fields

| Field | Description |
| :--- | :--- |
| `CliMain func(args []string) int` | Runs the whole interface over `args` and returns `ExitOk` (`0`), `ExitFailure` (`1`, a well-formed command that failed) or `ExitUsage` (`2`, the command line itself was wrong). An empty `args` prints the general help and returns `ExitUsage`. |

## Examples

```go
lib := agnoslib.New(&deps)

code := lib.Cli.CliMain([]string{"add-flag", "name", "--command", "greet", "--required", "--path", "./my-tool"})
switch code {
case api.ExitOk:
	// the flag was added and the project rebuilt
case api.ExitUsage:
	// e.g. a misspelled flag — nothing was written
case api.ExitFailure:
	// e.g. the command does not exist, or the rebuild did not compile
}
```
