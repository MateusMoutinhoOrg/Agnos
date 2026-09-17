# `exec-test`

Run the project's examples and check them against their goldens

```bash
agnos exec-test [--only <only>] [--update] [--path <path>] [--quiet]
```

Runs every example of examples/cli/ and examples/lib/ in alphabetical order, cli side first, each with its own directory as the working directory and the project's own cli in front of the PATH. Every run starts from a removed TestDir and a removed AssertDir, and what it produced - the merged output, the exit status and the sha256 of every file the example copied out of TestDir into AssertDir - is compared against the example's result.yaml, or written there when that golden does not exist yet. An example that copied nothing out fails: it asserted nothing. An example declared on both sides must leave the same tree and exit the same way: the cli is only a wrapper over the lib.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--only` | string |  | run a single example by name, both sides (defaults to every example) |
| `--update` | boolean |  | rewrite every golden result.yaml with what this run produced instead of comparing |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

```bash
agnos exec-test
agnos exec-test --only start
agnos exec-test --update
```

Examples · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
