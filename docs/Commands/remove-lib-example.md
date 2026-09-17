# `remove-lib-example`

Delete an example from examples/lib/

```bash
agnos remove-lib-example [--path <path>] [--quiet] <name>
```

Deletes examples/lib/<name>/ whole - the example.go, the golden result.yaml and any TestDir or AssertDir the last run left behind - and runs build so the example listing of the docs is rewritten without it.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the example directory under examples/lib/ |

```bash
agnos remove-lib-example start
```

Examples · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
