# `list-routes`

List every route in the order the chain runs them

```bash
agnos list-routes [--path <path>] [--quiet]
```

Prints one line per declared route — the rung, the methods, the pattern and the name — in the order the dispatch runs them: the chain first, lowest priority first, then the routes of the after phase. Writes nothing.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

```bash
agnos list-routes
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
