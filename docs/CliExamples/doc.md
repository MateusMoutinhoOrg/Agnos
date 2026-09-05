# CliExamples

Every example of the agnos cli. Each one is a shell session that runs with its own
directory as the working directory and writes only into its own `TestDir`, so it can be read
as documentation and copied line by line. The script types `agnos`, which `exec-test`
resolves to the code in this tree.

`agnos exec-test` runs them all and checks each against the `result.yaml` beside it — the
golden holding the output, the exit code and the sha256 of every `TestDir` file, written by
`exec-test` and never by hand. [Workflow](../Workflow/doc.md) has the commands that add and
remove one; the lib side is [LibExamples](../LibExamples/doc.md).

| Example | Description | Source |
|---|---|---|
| `add-arg` | declare one positional argument on a command | [example.sh](../../examples/cli/add-arg/example.sh) |
| `add-cli-example` | create an example under examples/cli/ | [example.sh](../../examples/cli/add-cli-example/example.sh) |
| `add-command` | declare a new command | [example.sh](../../examples/cli/add-command/example.sh) |
| `add-doc` | create a doc directory under docs/ | [example.sh](../../examples/cli/add-doc/example.sh) |
| `add-flag` | declare one flag on a command | [example.sh](../../examples/cli/add-flag/example.sh) |
| `add-lib-example` | create an example under examples/lib/ | [example.sh](../../examples/cli/add-lib-example/example.sh) |
| `build` | regenerate every generated file of a project | [example.sh](../../examples/cli/build/example.sh) |
| `cli-init` | add the cli layer to a project that has none | [example.sh](../../examples/cli/cli-init/example.sh) |
| `cli-purge` | remove the cli layer and every command in it | [example.sh](../../examples/cli/cli-purge/example.sh) |
| `compile` | cross-compile a project's cmd/main into release/ | [example.sh](../../examples/cli/compile/example.sh) |
| `dep-install` | install one dep of the built-in list | [example.sh](../../examples/cli/dep-install/example.sh) |
| `dep-list` | list the deps available to a project | [example.sh](../../examples/cli/dep-list/example.sh) |
| `dep-remove` | uninstall one installed dep | [example.sh](../../examples/cli/dep-remove/example.sh) |
| `deps-init` | add the dependency layer to a project that has none | [example.sh](../../examples/cli/deps-init/example.sh) |
| `deps-purge` | remove the dependency layer and every installed dep | [example.sh](../../examples/cli/deps-purge/example.sh) |
| `remove-arg` | delete one declared positional argument | [example.sh](../../examples/cli/remove-arg/example.sh) |
| `remove-cli-example` | delete an example of examples/cli/ | [example.sh](../../examples/cli/remove-cli-example/example.sh) |
| `remove-command` | delete a command and unwire its dispatch | [example.sh](../../examples/cli/remove-command/example.sh) |
| `remove-doc` | delete a doc directory | [example.sh](../../examples/cli/remove-doc/example.sh) |
| `remove-flag` | delete one declared flag | [example.sh](../../examples/cli/remove-flag/example.sh) |
| `remove-lib-example` | delete an example of examples/lib/ | [example.sh](../../examples/cli/remove-lib-example/example.sh) |
| `set-command` | rewrite the command-level keys of a command | [example.sh](../../examples/cli/set-command/example.sh) |
| `start` | scaffold a new project | [example.sh](../../examples/cli/start/example.sh) |
| `verify` | check a project against the schema, writing nothing | [example.sh](../../examples/cli/verify/example.sh) |

