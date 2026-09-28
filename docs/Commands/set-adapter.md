# `set-adapter`

Changes which adapter an available binds for one dep

```bash
agnos set-adapter <Dep> <Adapter> [--available <available>] [--help] [--path <path>] [--quiet]
```

Rewrites one available.yaml so the named adapter is the one bound for that dep, dropping whichever adapter filled the field before. It is the only editor of that choice.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Dep` | string, required |  | the dep whose field is being filled |
| `Adapter` | string, required |  | the installed adapter that should fill it |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--available` | string |  | the available to change (defaults to standard) | — |
| `--help` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos set-adapter sortdeps reflectsort
agnos set-adapter sortdeps reflectsort --available lambda
```

Deps System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
