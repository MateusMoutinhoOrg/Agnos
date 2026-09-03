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
}

type Result struct {
	Output   string
	ExitCode int
}
```

## Description

The process runner injected whole as `Deps.Rundeps` — the sandbox's copy of what `os/exec` provides. `Run` executes one program to completion and returns what it wrote, standard output and standard error merged in the order they were written. A non-zero exit status is reported in `Result.ExitCode`, **not** as the error: the error is reserved for a program that could not be started at all (missing binary, unreadable directory).

It is what the build pipeline reaches for when it hands a rendered project to a real toolchain — `go mod tidy`, `go build ./cmd/... ./sandbox/... ./adapters/...` — and what `start` asks `go env GOVERSION`; see [BuildPipeline.md](/docs/References/BuildPipeline.md#persist-then-runtime). Installed by the `rundeps` dep.

## Fields

| Field | Description |
| :--- | :--- |
| `Run(props)` | Runs `props.Program` with `props.Args` in `props.Dir` (`""` is the current directory). |
| `RunProps.Program` | The executable, looked up on `PATH`. |
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
