# LibExamples

Every example of agnos used as a Go module. Each one is a `package main` program that
runs with its own directory as the working directory and writes only into its own `TestDir`,
so it can be read as documentation and copied as a starting point.

`agnos exec-test` runs them all and checks each against the `result.yaml` beside it — the
golden holding the output, the exit code and the sha256 of every `TestDir` file, written by
`exec-test` and never by hand. [Workflow](../Workflow/doc.md) has the commands that add and
remove one; the cli side is [CliExamples](../CliExamples/doc.md).

| Example | Source |
|---|---|
| `add-arg` | [example.go](../../examples/lib/add-arg/example.go) |
| `add-cli-example` | [example.go](../../examples/lib/add-cli-example/example.go) |
| `add-command` | [example.go](../../examples/lib/add-command/example.go) |
| `add-doc` | [example.go](../../examples/lib/add-doc/example.go) |
| `add-flag` | [example.go](../../examples/lib/add-flag/example.go) |
| `add-lib-example` | [example.go](../../examples/lib/add-lib-example/example.go) |
| `build` | [example.go](../../examples/lib/build/example.go) |
| `cli-init` | [example.go](../../examples/lib/cli-init/example.go) |
| `cli-purge` | [example.go](../../examples/lib/cli-purge/example.go) |
| `compile` | [example.go](../../examples/lib/compile/example.go) |
| `dep-install` | [example.go](../../examples/lib/dep-install/example.go) |
| `dep-list` | [example.go](../../examples/lib/dep-list/example.go) |
| `dep-remove` | [example.go](../../examples/lib/dep-remove/example.go) |
| `deps-init` | [example.go](../../examples/lib/deps-init/example.go) |
| `deps-purge` | [example.go](../../examples/lib/deps-purge/example.go) |
| `remove-arg` | [example.go](../../examples/lib/remove-arg/example.go) |
| `remove-cli-example` | [example.go](../../examples/lib/remove-cli-example/example.go) |
| `remove-command` | [example.go](../../examples/lib/remove-command/example.go) |
| `remove-doc` | [example.go](../../examples/lib/remove-doc/example.go) |
| `remove-flag` | [example.go](../../examples/lib/remove-flag/example.go) |
| `remove-lib-example` | [example.go](../../examples/lib/remove-lib-example/example.go) |
| `set-command` | [example.go](../../examples/lib/set-command/example.go) |
| `start` | [example.go](../../examples/lib/start/example.go) |
| `verify` | [example.go](../../examples/lib/verify/example.go) |

