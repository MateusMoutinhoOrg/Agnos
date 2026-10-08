# `add-lib-example`

Scaffold a new example under examples/lib/

```bash
agnos add-lib-example <Name> [--help] [--path <path>] [--quiet]
```

Creates examples/lib/<name>/ with an example.go stub (package main) that already runs and exits 0, then runs build so the example listing of the docs names it. The golden result.yaml is written by the first run-examples, never by hand. Refuses an existing name.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the name of the new example (it becomes one directory under examples/lib/) |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

```bash
agnos add-lib-example start
```

Examples · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
