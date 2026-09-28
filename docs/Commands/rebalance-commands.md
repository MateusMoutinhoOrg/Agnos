# `rebalance-commands`

Lay the cli chain down again with room between its rungs

```bash
agnos rebalance-commands [--step <step>] [--help] [--path <path>] [--quiet]
```

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--step` | integer, >= 1 | `10` | how many rungs apart two commands land (defaults to 10) | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

Cli System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
