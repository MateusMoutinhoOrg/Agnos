# CliExamples

Every example of the agnos cli. Each one is a shell session that runs with its own
directory as the working directory and writes only into its own `TestDir`, so it can be read
as documentation and copied line by line. The script types `agnos`, which `exec-test`
resolves to the code in this tree.

`agnos exec-test` runs them all and checks each against the `result.yaml` beside it — the
golden holding the output, the exit code and the sha256 of every `TestDir` file, written by
`exec-test` and never by hand. [Workflow](../Workflow/doc.md) has the commands that add and
remove one; the lib side is [LibExamples](../LibExamples/doc.md).

| Example | Source |
|---|---|
| `add-arg` | [example.sh](../../examples/cli/add-arg/example.sh) |
| `add-cli-example` | [example.sh](../../examples/cli/add-cli-example/example.sh) |
| `add-command` | [example.sh](../../examples/cli/add-command/example.sh) |
| `add-doc` | [example.sh](../../examples/cli/add-doc/example.sh) |
| `add-flag` | [example.sh](../../examples/cli/add-flag/example.sh) |
| `add-lib-example` | [example.sh](../../examples/cli/add-lib-example/example.sh) |
| `build` | [example.sh](../../examples/cli/build/example.sh) |
| `cli-init` | [example.sh](../../examples/cli/cli-init/example.sh) |
| `cli-purge` | [example.sh](../../examples/cli/cli-purge/example.sh) |
| `compile` | [example.sh](../../examples/cli/compile/example.sh) |
| `dep-install` | [example.sh](../../examples/cli/dep-install/example.sh) |
| `dep-list` | [example.sh](../../examples/cli/dep-list/example.sh) |
| `dep-remove` | [example.sh](../../examples/cli/dep-remove/example.sh) |
| `deps-init` | [example.sh](../../examples/cli/deps-init/example.sh) |
| `deps-purge` | [example.sh](../../examples/cli/deps-purge/example.sh) |
| `remove-arg` | [example.sh](../../examples/cli/remove-arg/example.sh) |
| `remove-cli-example` | [example.sh](../../examples/cli/remove-cli-example/example.sh) |
| `remove-command` | [example.sh](../../examples/cli/remove-command/example.sh) |
| `remove-doc` | [example.sh](../../examples/cli/remove-doc/example.sh) |
| `remove-flag` | [example.sh](../../examples/cli/remove-flag/example.sh) |
| `remove-lib-example` | [example.sh](../../examples/cli/remove-lib-example/example.sh) |
| `set-command` | [example.sh](../../examples/cli/set-command/example.sh) |
| `start` | [example.sh](../../examples/cli/start/example.sh) |
| `verify` | [example.sh](../../examples/cli/verify/example.sh) |

