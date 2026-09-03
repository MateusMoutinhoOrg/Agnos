# `std.Lib`

**Type:** Struct

## Definition

```go
type Lib struct {
	Now    func() time.Time
	Printf func(format string, a ...any) (n int, err error)
	Log    func(format string, a ...any) (n int, err error)
	Error  func(format string, a ...any) (n int, err error)
	Errorf func(format string, a ...any) error
}
```

## Description

The runtime library injected whole as `Deps.Std`: the clock and the process's output streams, which the sandbox may not reach through `time`, `fmt` or `os`. `deps.Deps` carries no bare function fields, so the loose runtime functions are gathered here and injected as one sub-contract like every other.

Its three writers are the CLI's three **output channels**, one meaning each, so no handler has to choose a stream — see [Commands.md](/docs/References/Commands.md#output-channels). `Printf` carries a command's *result* to standard output and is never silenced. `Log` carries *progress* to standard error and is the only channel `--quiet` turns off: the generated dispatch replaces `deps.Std.Log` with a no-op as soon as it has read a boolean flag named `quiet`. `Error` carries usage errors and failures to standard error. `Errorf` builds the `error` values actions return. Installed by the `std` dep, which `cli-init` pulls in.

## Fields

| Field | Description |
| :--- | :--- |
| `Now()` | The current wall-clock time. |
| `Printf(format, a...)` | stdout. The result. |
| `Log(format, a...)` | stderr. Progress; silenced by `--quiet`. |
| `Error(format, a...)` | stderr. Failures and usage errors. |
| `Errorf(format, a...)` | Formats and returns an `error`; `%w` wraps. |

## Examples

```go
func CommandHandler(deps *deps.Deps, entries *Entries) int {
	deps.Std.Log("listing %s\n", entries.Path)          // progress, off under --quiet
	names, err := list(deps, entries.Path)
	if err != nil {
		deps.Std.Error("%v\n", err)                        // failure, always shown
		return api.ExitFailure
	}
	for _, name := range names {
		deps.Std.Printf("%s\n", name)                      // the result, always shown
	}
	return api.ExitOk
}
```
