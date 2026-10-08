# `add-binding`

Declare one further binding

```bash
agnos add-binding <Name> [--help] [--path <path>] [--quiet]
```

Creates adapters/bindings/<name>/binding.yaml as a copy of the standard selection, so it starts filling every field, and build generates its new.go. Point it at another adapter with set-adapter --binding.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the name of the new binding (it becomes one directory under adapters/bindings/) |

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
agnos add-binding lambda
```

Deps · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
