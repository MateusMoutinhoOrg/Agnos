# `remove-binding`

Delete one binding

```bash
agnos remove-binding <Name> [--help] [--path <path>] [--quiet]
```

Removes adapters/bindings/<name>/ whole. The standard binding is refused: cmd/main/main.go imports it.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the binding to remove |

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
agnos remove-binding lambda
```

Deps · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
