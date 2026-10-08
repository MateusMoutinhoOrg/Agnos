# `list-routes`

List every route in the order the chain runs them

```bash
agnos list-routes [--help] [--path <path>] [--quiet]
```

Prints one line per declared route — the rung, the methods, the pattern and the name — in the order the dispatch runs them, lowest priority first. Writes nothing.

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
agnos list-routes
```

Server · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
