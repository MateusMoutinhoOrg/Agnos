# Public API Index
Every public-facing entry of the library, grouped by role, with a page per symbol

| Doc | Description |
| --- | --- |
| [sandbox.New](/docs/PublicApi/sandbox.New/doc.md) | Injects a `*deps.Deps` into the sandbox and returns the `*api.Sandbox` carrying every action and the command-line interface |
| [standard.New](/docs/PublicApi/standard.New/doc.md) | Builds a `deps.Deps` by running every adapter lib's `Bind` — the default assembly `cmd/main` wires |
| [api.Sandbox](/docs/PublicApi/api.Sandbox/doc.md) | What `sandbox.New` hands back: an `Actions` struct and a `Cli` struct, one field per file of `sandbox/api/` |
| [api.Cli](/docs/PublicApi/api.Cli/doc.md) | The whole command-line interface as one function, plus the three exit-code constants |
| [api.Actions](/docs/PublicApi/api.Actions/doc.md) | Every operation `agnos` performs, as a function field on one struct |
| [api.BuildProps](/docs/PublicApi/api.BuildProps/doc.md) | `Path` and `Runtime` for one build |
| [api.CompileProps](/docs/PublicApi/api.CompileProps/doc.md) | `Path` and `Targets` for one cross-compile run |
| [api.StartProps](/docs/PublicApi/api.StartProps/doc.md) | `Path`, `ProjectName`, optional `Module` and `Force` for one scaffold |
| [api.FieldProps](/docs/PublicApi/api.FieldProps/doc.md) | One flag or positional argument to add to a command's declaration |
| [api.CommandProps](/docs/PublicApi/api.CommandProps/doc.md) | The command-level keys `SetCommand` may rewrite |
| [api.DocProps](/docs/PublicApi/api.DocProps/doc.md) | One documentation page to create: its directory, summary and themes |
| [deps.Deps](/docs/PublicApi/deps.Deps/doc.md) | The dependency contract every adapter fills: one sub-contract struct per directory of `sandbox/deps/` |
| [argvdeps.Lib](/docs/PublicApi/argvdeps.Lib/doc.md) | `Deps.Argvdeps` — a per-call argument-vector parser, the sandbox's copy of the Verb library |
| [embeddeps.Lib](/docs/PublicApi/embeddeps.Lib/doc.md) | `Deps.Embeddeps` — read-only access to the assets compiled into the binary, and template rendering over them |
| [goimportsdeps.Lib](/docs/PublicApi/goimportsdeps.Lib/doc.md) | `Deps.Goimportsdeps` — a Go source reader: package clause, imports, top-level declarations |
| [iodeps.Lib](/docs/PublicApi/iodeps.Lib/doc.md) | `Deps.Iodeps` — the filesystem: read, write, predicates, directory ops, listings |
| [requestdeps.Lib](/docs/PublicApi/requestdeps.Lib/doc.md) | `Deps.Requestdeps` — a per-call HTTP request, with its `Request` and `Response` structs |
| [rundeps.Lib](/docs/PublicApi/rundeps.Lib/doc.md) | `Deps.Rundeps` — running one external program to completion |
| [serializables.Lib](/docs/PublicApi/serializables.Lib/doc.md) | `Deps.Serializables` — generic JSON/YAML values: create, parse, walk, serialize |
| [std.Lib](/docs/PublicApi/std.Lib/doc.md) | `Deps.Std` — the clock and the three output channels `Printf` / `Log` / `Error`, plus `Errorf` |
