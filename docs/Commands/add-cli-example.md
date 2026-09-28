# `add-cli-example`

Scaffold a new example under examples/cli/

```bash
agnos add-cli-example <Name> [--help] [--path <path>] [--quiet]
```

Creates examples/cli/<name>/ with an example.sh stub that already runs and exits 0, then runs build so the example listing of the docs names it. The golden result.yaml is written by the first exec-test, never by hand. Refuses an existing name, and refuses outright in a project with no cli.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the name of the new example (it becomes one directory under examples/cli/) |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos add-cli-example start
```

Examples · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
