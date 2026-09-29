# `rebalance-routes`

Lay the chain down again with room between its rungs

```bash
agnos rebalance-routes [--step <step>] [--help] [--path <path>] [--quiet]
```

Gives every route a priority of its own, --step apart, in the order the chain runs them now, so --before and --after have room again. The generated health route keeps its rung.

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--step` | integer | `10` | how many rungs apart two routes land (defaults to 10) | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

Server System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
