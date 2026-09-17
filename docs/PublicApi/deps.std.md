# `deps.Std`

`sandbox/deps/std`

## `Sandbox`

Sandbox is the runtime library injected whole as the Deps.Std field.

| Field | Type | Description |
| --- | --- | --- |
| `Now` | `func() int64` | Now returns the current wall-clock time as nanoseconds since the Unix epoch, UTC. The sandbox may not name a `time.Time`, so an instant crosses this boundary as a plain integer. |
| `Printf` | `func(format string, a ...any) (n int, err error)` | Printf writes one formatted message to standard output. It carries the command's result — the data a script would read — so it is never silenced. |
| `Log` | `func(format string, a ...any) (n int, err error)` | Log writes one formatted progress message to standard error. It is the channel every "… started with path …" notice goes through, so a caller can keep stdout free of log noise, and it is what --quiet turns off. |
| `Error` | `func(format string, a ...any) (n int, err error)` | Error writes one formatted message to standard error. |
| `Errorf` | `func(format string, a ...any) error` | Errorf formats an error message and returns it as an error. |
| `Sprintf` | `func(format string, a ...any) string` | Sprintf formats a message and returns it as a string. It is the one formatting entry point the sandbox has: every string it builds out of values rather than out of concatenation goes through here. |
| `Goos` | `func() string` | Goos is the name of the operating system the process runs on, in the spelling the Go toolchain uses ("darwin", "linux", "windows", …). |

[every contract](doc.md)
