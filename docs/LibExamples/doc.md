# LibExamples

Every example of agnos used as a Go module. Each one is a `package main` program that
runs with its own directory as the working directory and writes only into its own `TestDir`,
so it can be read as documentation and copied as a starting point.

`agnos exec-test` runs them all and checks each against the `result.yaml` beside it — the
golden holding the output, the exit code and the sha256 of every `TestDir` file, written by
`exec-test` and never by hand. [Workflow](../Workflow/doc.md) has the commands that add and
remove one; the cli side is [CliExamples](../CliExamples/doc.md).

| Example | Description | Source |
|---|---|---|
| `add-arg` | declare one positional argument on a command | [example.go](../../examples/lib/add-arg/example.go) |
| `add-cli-example` | create an example under examples/cli/ | [example.go](../../examples/lib/add-cli-example/example.go) |
| `add-command` | declare a new command | [example.go](../../examples/lib/add-command/example.go) |
| `add-doc` | create a doc directory under docs/ | [example.go](../../examples/lib/add-doc/example.go) |
| `add-flag` | declare one flag on a command | [example.go](../../examples/lib/add-flag/example.go) |
| `add-lib-example` | create an example under examples/lib/ | [example.go](../../examples/lib/add-lib-example/example.go) |
| `build` | regenerate every generated file of a project | [example.go](../../examples/lib/build/example.go) |
| `cli-init` | add the cli layer to a project that has none | [example.go](../../examples/lib/cli-init/example.go) |
| `cli-purge` | remove the cli layer and every command in it | [example.go](../../examples/lib/cli-purge/example.go) |
| `compile` | cross-compile a project's cmd/main into release/ | [example.go](../../examples/lib/compile/example.go) |
| `dep-install` | install one dep of the built-in list | [example.go](../../examples/lib/dep-install/example.go) |
| `dep-list` | list the deps available to a project | [example.go](../../examples/lib/dep-list/example.go) |
| `dep-remove` | uninstall one installed dep | [example.go](../../examples/lib/dep-remove/example.go) |
| `deps-init` | add the dependency layer to a project that has none | [example.go](../../examples/lib/deps-init/example.go) |
| `deps-purge` | remove the dependency layer and every installed dep | [example.go](../../examples/lib/deps-purge/example.go) |
| `remove-arg` | delete one declared positional argument | [example.go](../../examples/lib/remove-arg/example.go) |
| `remove-cli-example` | delete an example of examples/cli/ | [example.go](../../examples/lib/remove-cli-example/example.go) |
| `remove-command` | delete a command and unwire its dispatch | [example.go](../../examples/lib/remove-command/example.go) |
| `remove-doc` | delete a doc directory | [example.go](../../examples/lib/remove-doc/example.go) |
| `remove-flag` | delete one declared flag | [example.go](../../examples/lib/remove-flag/example.go) |
| `remove-lib-example` | delete an example of examples/lib/ | [example.go](../../examples/lib/remove-lib-example/example.go) |
| `set-command` | rewrite the command-level keys of a command | [example.go](../../examples/lib/set-command/example.go) |
| `start` | scaffold a new project | [example.go](../../examples/lib/start/example.go) |
| `verify` | check a project against the schema, writing nothing | [example.go](../../examples/lib/verify/example.go) |

