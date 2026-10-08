# `add-adapter`

Install one further adapter for a contract the project already has

```bash
agnos add-adapter <Name> [--binding <binding>] [--help] [--path <path>] [--quiet]
```

Renders assets/adapter-catalog/<adapter> into the project and writes its declaration to adapters/impls/<adapter>/adapter.yaml. The contract it fills has to be installed already. Installing changes no selection: a binding binds one adapter per field, so --binding names the one that switches to it.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the adapter to install from assets/adapter-catalog |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--binding` | string |  | the binding that should switch to this adapter (installs only when absent) | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

```bash
agnos add-adapter reflectsort
agnos add-adapter reflectsort --binding lambda
```

Deps · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
