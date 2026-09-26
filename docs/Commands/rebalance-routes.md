# `rebalance-routes`

Lay the chain down again with room between its rungs

```bash
agnos rebalance-routes [--step <step>] [--path <path>] [--quiet]
```

Gives every route a priority of its own, --step apart, in the order the chain runs them now — each phase on its own ladder — so --before and --after have room again. The generated health route keeps its rung.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--step` | int | `10` | how many rungs apart two routes land (defaults to 10) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
