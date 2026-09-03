# Public API

## Description
Index of every public-facing entry of the library, grouped by role. Callers hold **structs of function fields** declared in `sandbox/api` and `sandbox/deps`; the implementations that fill those fields live in `sandbox/internal` and are unreachable from outside `sandbox/`. See [StructContracts.md](/docs/References/StructContracts.md). Every symbol is imported from `github.com/MateusMoutinhoOrg/Agnos-Cli/<package path>`.

---

## Entry Points

### [sandbox.New](/docs/References/PublicApi/sandbox.New.md)
Injects a `*deps.Deps` into the sandbox and returns the `*api.Sandbox` carrying every action and the command-line interface.

### [standard.New](/docs/References/PublicApi/standard.New.md)
Builds a `deps.Deps` by running every adapter lib's `Bind` — the default assembly `cmd/main` wires.

---

## Core Interface

### [api.Sandbox](/docs/References/PublicApi/api.Sandbox.md)
What `sandbox.New` hands back: an `Actions` struct and a `Cli` struct, one field per file of `sandbox/api/`.

### [api.Cli](/docs/References/PublicApi/api.Cli.md)
The whole command-line interface as one function, `CliMain(args []string) int`, plus the three exit-code constants.

### [api.Actions](/docs/References/PublicApi/api.Actions.md)
Every operation `agnos` performs, as a function field: `Build`, `Verify`, `Start`, the deps and cli subsystems, and the command editors.

---

## Actions

### [api.Actions.Build](/docs/References/PublicApi/api.Actions.md#build)
Re-renders a project's generated files and hands it to a runtime.

### [api.Actions.Verify](/docs/References/PublicApi/api.Actions.md#verify)
Checks a project against the harness schema without writing.

### [api.Actions.Start](/docs/References/PublicApi/api.Actions.md#start)
Scaffolds a project: configuration, `go.mod`, sandbox skeleton, then `build`.

### [api.Actions.DepsInit / DepsPurge](/docs/References/PublicApi/api.Actions.md#depsinit--depspurge)
Turn the dependency-injection layer on and off.

### [api.Actions.DepInstall / DepRemove / DepList](/docs/References/PublicApi/api.Actions.md#depinstall--depremove--deplist)
Render, delete, and enumerate installable deps.

### [api.Actions.CliInit / CliPurge](/docs/References/PublicApi/api.Actions.md#cliinit--clipurge)
Install and remove the command-line layer.

### [api.Actions.AddCommand / RemoveCommand / SetCommand](/docs/References/PublicApi/api.Actions.md#addcommand--removecommand--setcommand)
Create, delete, and rewrite a command's declaration.

### [api.Actions.AddFlag / RemoveFlag / AddArg / RemoveArg](/docs/References/PublicApi/api.Actions.md#addflag--removeflag--addarg--removearg)
Add and drop a command's fields.

---

## Props

### [api.BuildProps](/docs/References/PublicApi/api.BuildProps.md)
`Path` and `Runtime` (`api.RuntimeGo` / `api.RuntimeNone`) for one build.

### [api.StartProps](/docs/References/PublicApi/api.StartProps.md)
`Path`, `ProjectName`, optional `Module`, `Force` for one scaffold.

### [api.FieldProps](/docs/References/PublicApi/api.FieldProps.md)
One flag or positional argument to add: name, identifiers, type, description, examples, and the raw `Default` / `Min` / `Max` literals.

### [api.CommandProps](/docs/References/PublicApi/api.CommandProps.md)
The command-level keys `SetCommand` may rewrite.

---

## Dependency Contracts

### [deps.Deps](/docs/References/PublicApi/deps.Deps.md)
The dependency contract every adapter fills: one sub-contract struct per directory of `sandbox/deps/`, named mechanically after it.

### [argvdeps.Lib](/docs/References/PublicApi/argvdeps.Lib.md)
`Deps.Argvdeps` — a per-call argument-vector parser, the sandbox's copy of the Verb library.

### [embeddeps.Lib](/docs/References/PublicApi/embeddeps.Lib.md)
`Deps.Embeddeps` — read-only access to the assets compiled into the binary, and template rendering over them.

### [goimportsdeps.Lib](/docs/References/PublicApi/goimportsdeps.Lib.md)
`Deps.Goimportsdeps` — a Go source reader: package clause, imports, top-level declarations.

### [iodeps.Lib](/docs/References/PublicApi/iodeps.Lib.md)
`Deps.Iodeps` — the filesystem: read, write, predicates, directory ops, listings.

### [requestdeps.Lib](/docs/References/PublicApi/requestdeps.Lib.md)
`Deps.Requestdeps` — a per-call HTTP request, with its `Request` and `Response` structs.

### [rundeps.Lib](/docs/References/PublicApi/rundeps.Lib.md)
`Deps.Rundeps` — running one external program to completion.

### [serializables.Lib](/docs/References/PublicApi/serializables.Lib.md)
`Deps.Serializables` — generic JSON/YAML values: create, parse, walk, serialize.

### [std.Lib](/docs/References/PublicApi/std.Lib.md)
`Deps.Std` — the clock and the three output channels `Printf` / `Log` / `Error`, plus `Errorf`.
