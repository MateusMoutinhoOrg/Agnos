# `remove-cli-example`

Delete an example from examples/cli/

```bash
agnos remove-cli-example [--path <path>] [--quiet] <name>
```

Deletes examples/cli/<name>/ whole - the example.sh, the golden result.yaml and any TestDir or AssertDir the last run left behind - and runs build so the example listing of the docs is rewritten without it.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the example directory under examples/cli/ |

```bash
agnos remove-cli-example start
```

Examples · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
