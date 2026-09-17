# `set-adapter`

Changes which adapter an available binds for one dep

```bash
agnos set-adapter [--available <available>] [--path <path>] [--quiet] <dep> <adapter>
```

Rewrites one available.yaml so the named adapter is the one bound for that dep, dropping whichever adapter filled the field before. It is the only editor of that choice.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--available` | string |  | the available to change (defaults to standard) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `dep` | string, required |  | the dep whose field is being filled |
| `adapter` | string, required |  | the installed adapter that should fill it |

```bash
agnos set-adapter sortdeps reflectsort
agnos set-adapter sortdeps reflectsort --available lambda
```

Deps System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
