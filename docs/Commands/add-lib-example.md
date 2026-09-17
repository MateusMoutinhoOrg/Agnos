# `add-lib-example`

Scaffold a new example under examples/lib/

```bash
agnos add-lib-example [--path <path>] [--quiet] <name>
```

Creates examples/lib/<name>/ with an example.go stub (package main) that already runs and exits 0, then runs build so the example listing of the docs names it. The golden result.yaml is written by the first exec-test, never by hand. Refuses an existing name.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the name of the new example (it becomes one directory under examples/lib/) |

```bash
agnos add-lib-example start
```

Examples · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
