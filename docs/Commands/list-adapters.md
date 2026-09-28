# `list-adapters`

Lists the adapters of the catalog and of the project

```bash
agnos list-adapters [--help] [--path <path>] [--quiet]
```

One row per adapter: the name, the dep it fills, whether it is installed, the availables binding it, and what backs it.

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
agnos list-adapters
```

Deps System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
