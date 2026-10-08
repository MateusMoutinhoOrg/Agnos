# `remove-lib-example`

Delete an example from examples/lib/

```bash
agnos remove-lib-example <Name> [--help] [--path <path>] [--quiet]
```

Deletes examples/lib/<name>/ whole - the example.go, the golden result.yaml and any test-dir or assert-dir the last run left behind - and runs build so the example listing of the docs is rewritten without it.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the example directory under examples/lib/ |

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
agnos remove-lib-example start
```

Examples · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
