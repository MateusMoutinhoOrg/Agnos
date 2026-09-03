# `rundeps.Lib`

**Type:** Struct

## Definition

```go
type Lib struct {
	Run func(props RunProps) (Result, error)
}

type RunProps struct {
	Dir     string
	Program string
	Args    []string
	Env     []string
}

type Result struct {
	Output   string
	ExitCode int
}
```

## Description

The process runner injected whole as `Deps.Rundeps` — the sandbox's copy of what `os/exec` provides. `Run` executes one program to completion and returns what it wrote, standard output and standard error merged in the order they were written. A non-zero exit status is reported in `Result.ExitCode`, **not** as the error: the error is reserved for a program that could not be started at all (missing binary, unreadable directory).

It is what the build pipeline reaches for when it hands a rendered project to a real toolchain — `go mod tidy`, `go build ./cmd/... ./sandbox/... ./adapters/...` — what `start` asks `go env GOVERSION`, and what the `compile` action runs `go build` through once per cross-compile target; see [BuildPipeline.md](/docs/References/BuildPipeline.md#persist-then-runtime). Installed by the `rundeps` dep.

`RunProps.Env` is a list of `KEY=VALUE` entries layered on top of the current process environment for that one invocation (later entries win); an empty slice inherits the environment unchanged. `compile` uses it to set `CGO_ENABLED=0`, `GOOS` and `GOARCH`.

## Fields

| Field | Description |
| :--- | :--- |
| `Run(props)` | Runs `props.Program` with `props.Args` in `props.Dir` (`""` is the current directory). |
| `RunProps.Program` | The executable, looked up on `PATH`. |
| `RunProps.Env` | `KEY=VALUE` entries added on top of the current environment; empty inherits it unchanged. |
| `Result.Output` | Merged stdout and stderr. |
| `Result.ExitCode` | The exit status; `0` is success. |

## Examples

```go
result, err := deps.Rundeps.Run(rundeps.RunProps{
	Dir:     path,
	Program: "go",
	Args:    []string{"build", "./cmd/..."},
})
if err != nil {
	return deps.Std.Errorf("could not run go: %w", err) // no toolchain
}
if result.ExitCode != 0 {
	return deps.Std.Errorf("go build failed:\n%s", result.Output) // a compile error
}
```
