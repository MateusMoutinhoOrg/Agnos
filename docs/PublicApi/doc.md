# Public API

## Description
Index of every public-facing entry of the library, grouped by role. Callers hold **structs of function fields** declared in `sandbox/api` and `sandbox/deps`; the implementations that fill those fields live in `sandbox/internal` and are unreachable from outside `sandbox/`. See [StructContracts](/docs/StructContracts/doc.md). Every symbol is imported from `github.com/MateusMoutinhoOrg/Agnos/<package path>`. Each symbol below is a sub-doc of this page; the generated [Index.md](/docs/PublicApi/Index.md) lists them flat.

---

## Entry Points

### [sandbox.New](/docs/PublicApi/sandbox.New/doc.md)
Injects a `*deps.Deps` into the sandbox and returns the `*api.Sandbox` carrying every action and the command-line interface.

### [standard.New](/docs/PublicApi/standard.New/doc.md)
Builds a `deps.Deps` by running every adapter lib's `Bind` — the default assembly `cmd/main` wires.

---

## Core Interface

### [api.Sandbox](/docs/PublicApi/api.Sandbox/doc.md)
What `sandbox.New` hands back: an `Actions` struct and a `Cli` struct, one field per file of `sandbox/api/`.

### [api.Cli](/docs/PublicApi/api.Cli/doc.md)
The whole command-line interface as one function, `CliMain(args []string) int`, plus the three exit-code constants.

### [api.Actions](/docs/PublicApi/api.Actions/doc.md)
Every operation `agnos` performs, as a function field: `Build`, `Compile`, `Verify`, `Start`, the deps and cli subsystems, and the command editors.

---

## Actions

### [api.Actions.Build](/docs/PublicApi/api.Actions/doc.md#build)
Re-renders a project's generated files and hands it to a runtime.

### [api.Actions.Compile](/docs/PublicApi/api.Actions/doc.md#compile)
Runs `Build`, then cross-compiles `./cmd/main` into `release/` once per target.

### [api.Actions.Verify](/docs/PublicApi/api.Actions/doc.md#verify)
Checks a project against the harness schema without writing.

### [api.Actions.Start](/docs/PublicApi/api.Actions/doc.md#start)
Scaffolds a project: configuration, `go.mod`, sandbox skeleton, then `build`.

### [api.Actions.DepsInit / DepsPurge](/docs/PublicApi/api.Actions/doc.md#depsinit--depspurge)
Turn the dependency-injection layer on and off.

### [api.Actions.DepInstall / DepRemove / DepList](/docs/PublicApi/api.Actions/doc.md#depinstall--depremove--deplist)
Render, delete, and enumerate installable deps.

### [api.Actions.CliInit / CliPurge](/docs/PublicApi/api.Actions/doc.md#cliinit--clipurge)
Install and remove the command-line layer.

### [api.Actions.AddCommand / RemoveCommand / SetCommand](/docs/PublicApi/api.Actions/doc.md#addcommand--removecommand--setcommand)
Create, delete, and rewrite a command's declaration.

### [api.Actions.AddFlag / RemoveFlag / AddArg / RemoveArg](/docs/PublicApi/api.Actions/doc.md#addflag--removeflag--addarg--removearg)
Add and drop a command's fields.

### [api.Actions.AddDoc / RemoveDoc](/docs/PublicApi/api.Actions/doc.md#adddoc--removedoc)
Create and delete a doc directory of `docs/`, leaving every index to the build.

---

## Props

### [api.BuildProps](/docs/PublicApi/api.BuildProps/doc.md)
`Path` and `Runtime` (`api.RuntimeGo` / `api.RuntimeNone`) for one build.

### [api.CompileProps](/docs/PublicApi/api.CompileProps/doc.md)
`Path` and `Targets` (target names, or `"all"`) for one cross-compile run.

### [api.StartProps](/docs/PublicApi/api.StartProps/doc.md)
`Path`, `ProjectName`, optional `Module`, `Force` for one scaffold.

### [api.FieldProps](/docs/PublicApi/api.FieldProps/doc.md)
One flag or positional argument to add: name, identifiers, type, description, examples, and the raw `Default` / `Min` / `Max` literals.

### [api.CommandProps](/docs/PublicApi/api.CommandProps/doc.md)
The command-level keys `SetCommand` may rewrite.

### [api.DocProps](/docs/PublicApi/api.DocProps/doc.md)
One doc to create: its directory under `docs/`, its one-line description, and its themes.

---

## Dependency Contracts

### [deps.Deps](/docs/PublicApi/deps.Deps/doc.md)
The dependency contract every adapter fills: one sub-contract struct per directory of `sandbox/deps/`, named mechanically after it.

### [argvdeps.Lib](/docs/PublicApi/argvdeps.Lib/doc.md)
`Deps.Argvdeps` — a per-call argument-vector parser, the sandbox's copy of the Verb library.

### [embeddeps.Lib](/docs/PublicApi/embeddeps.Lib/doc.md)
`Deps.Embeddeps` — read-only access to the assets compiled into the binary, and template rendering over them.

### [goimportsdeps.Lib](/docs/PublicApi/goimportsdeps.Lib/doc.md)
`Deps.Goimportsdeps` — a Go source reader: package clause, imports, top-level declarations.

### [iodeps.Lib](/docs/PublicApi/iodeps.Lib/doc.md)
`Deps.Iodeps` — the filesystem: read, write, predicates, directory ops, listings.

### [requestdeps.Lib](/docs/PublicApi/requestdeps.Lib/doc.md)
`Deps.Requestdeps` — a per-call HTTP request, with its `Request` and `Response` structs.

### [rundeps.Lib](/docs/PublicApi/rundeps.Lib/doc.md)
`Deps.Rundeps` — running one external program to completion.

### [serializables.Lib](/docs/PublicApi/serializables.Lib/doc.md)
`Deps.Serializables` — generic JSON/YAML values: create, parse, walk, serialize.

### [std.Lib](/docs/PublicApi/std.Lib/doc.md)
`Deps.Std` — the clock and the three output channels `Printf` / `Log` / `Error`, plus `Errorf`.
