# `list-deps`

List the deps the embedded catalog can install

```bash
agnos list-deps [--help] [--path <path>] [--quiet]
```

One row per dep of the embedded catalog: the name, whether the project has the contract installed, the adapters filling it (or the catalog's default-adapter when none is), and what it provides.

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
agnos list-deps
```

Deps · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
