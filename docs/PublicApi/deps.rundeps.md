# `deps.Rundeps`

`sandbox/deps/rundeps`

## `Sandbox`

Sandbox is the process runner injected whole as the Deps.Rundeps field. It is what the build action reaches for when it has to hand the rendered project to a real toolchain (`go mod tidy`, `go build ./...`) and report whether that toolchain accepted it.

| Field | Type | Description |
| --- | --- | --- |
| `Run` | `func(props RunProps) (Result, error)` | Run executes one program to completion and returns what it wrote. A non-zero exit status is reported in Result.ExitCode, not as the error: the error is reserved for a program that could not be started at all (missing binary, unreadable directory). |

## `RunProps`

RunProps describes one program invocation.

| Field | Type | Description |
| --- | --- | --- |
| `Dir` | `string` | Dir is the working directory the program runs in. "" means the current directory. |
| `Program` | `string` | Program is the executable to run, looked up on PATH. |
| `Args` | `[]string` | Args are the arguments handed to the program, excluding its own name. |
| `Env` | `[]string` | Env is a list of "KEY=VALUE" entries added on top of the current process environment for this one invocation (later entries win). Empty means "inherit the environment unchanged" — the common case. `agnos compile` uses it to set GOOS/GOARCH/CGO_ENABLED per cross-compile. |
| `PathPrefix` | `[]string` | PathPrefix are directories prepended to the PATH the program sees, ahead of the inherited one, and searched first when Program itself is looked up. A PATH entry cannot be expressed through Env: the adapter is what reads the current PATH and joins it, because the sandbox cannot. `agnos exec-test` uses it to put the project's own cli alias in front of the PATH an example runs with. |

## `Result`

Result is what one finished invocation produced.

| Field | Type | Description |
| --- | --- | --- |
| `Output` | `string` | Output is the program's standard output and standard error, merged in the order they were written. |
| `ExitCode` | `int` | ExitCode is the program's exit status; 0 means success. |

[every contract](doc.md)
