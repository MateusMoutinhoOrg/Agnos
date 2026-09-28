# `list-deps`

Lists the deps the embedded catalog can install

```bash
agnos list-deps [--help] [--path <path>] [--quiet]
```

One row per dep of the embedded catalog: the name, whether the project has the contract installed, the adapters filling it (or the catalog's default-adapter when none is), and what it provides.

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--help` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos list-deps
```

Deps System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
