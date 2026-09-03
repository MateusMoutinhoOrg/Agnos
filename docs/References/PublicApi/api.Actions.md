# `api.Actions`

**Type:** Struct

## Definition

```go
type Actions struct {
	Build         func(props BuildProps) error
	Compile       func(props CompileProps) error
	Verify        func(path string) error
	Start         func(props StartProps) error
	DepsInit      func(path string) error
	DepsPurge     func(path string) error
	DepInstall    func(path string, dep string) error
	DepRemove     func(path string, dep string) error
	DepList       func(path string) ([]string, error)
	CliInit       func(path string) error
	CliPurge      func(path string) error
	AddCommand    func(path string, name string, help string, category string) error
	RemoveCommand func(path string, name string) error
	SetCommand    func(props CommandProps) error
	AddFlag       func(props FieldProps) error
	RemoveFlag    func(path string, command string, name string) error
	AddArg        func(props FieldProps) error
	RemoveArg     func(path string, command string, name string) error
}
```

## Description

Every operation `agnos` performs, as function fields filled by `sandbox/binds/actions.go` from the packages under `sandbox/internal/actions/`. Each command's handler calls the same function its field points at, so the library and the CLI cannot drift. Every action takes the project directory (`path`, or `Path` on a props struct) and scopes every read and write to it — the `--path` flag by another name — reports progress through `deps.Std.Log`, and returns an `error` where the command would exit `1`. The actions that add something end by running `Build` with the `go` runtime; the ones that remove something run it with `none`, because hand-written code may still refer to what is gone. The full behavior of each is in [Commands.md](/docs/References/Commands.md).

## Fields

### Build
`Build(props BuildProps) error` — re-renders every generated file of the project and hands it to `props.Runtime`. Unlike the `build` command it does **not** run `Verify` first, so a mid-refactor tree can still regenerate. See [BuildPipeline.md](/docs/References/BuildPipeline.md).

### Compile
`Compile(props CompileProps) error` — runs `Build` with the `go` runtime, creates `release/`, then cross-compiles `./cmd/main` once per name in `props.Targets` into `release/<name>` with `CGO_ENABLED=0` and the target's `GOOS`/`GOARCH`. `"all"` expands to every target; an unknown or empty target is an error. See [Commands.md](/docs/References/Commands.md) for the name table.

### Verify
`Verify(path string) error` — checks the tree against the harness schema and writes nothing. The error lists every violation at once. Does not run a runtime; the `verify` command does.

### Start
`Start(props StartProps) error` — renders the `start` group into `AgnosConfig/`, writes `go.mod` when `props.Module` is given (required when the directory has none), then renders and compiles like `Build` with the `go` runtime.

### DepsInit / DepsPurge
`DepsInit(path) error` creates `sandbox/deps/` and `adapters/`, then builds; `DepsPurge(path) error` removes both whole, then builds with `none`.

### DepInstall / DepRemove / DepList
`DepInstall(path, dep) error` renders `assets/deplist/<dep>/**`, syncs `go.mod`, persists, builds; an unknown dep is an error. `DepRemove(path, dep) error` deletes those files and emptied directories, strips the `require`, builds with `none`. `DepList(path) ([]string, error)` returns one name per dep — the one action with a result, which the command prints.

### CliInit / CliPurge
`CliInit(path) error` installs the `std` and `argvdeps` deps, renders the `cli` group, builds. `CliPurge(path) error` removes the group's files plus `sandbox/internal/cli/` and `sandbox/internal/commands/` whole, builds with `none`.

### AddCommand / RemoveCommand / SetCommand
`AddCommand(path, name, help, category) error` normalizes `name`, refuses `help` and an existing command, writes `entries.yaml` and a stub `handler.go`, builds. `RemoveCommand(path, name) error` deletes the package, builds with `none`. `SetCommand(props CommandProps) error` rewrites the command-level keys, builds.

### AddFlag / RemoveFlag / AddArg / RemoveArg
`AddFlag(props FieldProps) error` and `AddArg(props FieldProps) error` validate the field (type, literals, `Required` against booleans and defaults, an array arg staying last), insert it at `props.Position`, build. `RemoveFlag(path, command, name) error` matches by name or identifier; `RemoveArg(path, command, name) error` by name; both build with `none`.

## Examples

```go
lib := agnoslib.New(&deps)
path := "./my-tool"

deps, err := lib.Actions.DepList(path)
if err != nil {
	log.Fatal(err)
}
for _, dep := range deps {
	fmt.Println(dep) // argvdeps, dbdeps, embeddeps, …
}

if err := lib.Actions.DepInstall(path, "iodeps"); err != nil {
	log.Fatal(err) // unknown dep, or the rebuild did not compile
}
```
