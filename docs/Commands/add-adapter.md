# `add-adapter`

Installs one further adapter for a contract the project already has

```bash
agnos add-adapter <Adapter> [--available <available>] [--help] [--path <path>] [--quiet]
```

Renders assets/adapterlist/<adapter> into the project and writes its declaration to adapters/libs/<adapter>/adapter.yaml. The contract it fills has to be installed already. Installing changes no selection: an available binds one adapter per field, so --available names the one that switches to it.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Adapter` | string, required |  | the adapter to install from assets/adapterlist |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--available` | string |  | the available that should switch to this adapter (installs only when absent) | — |
| `--help` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos add-adapter reflectsort
agnos add-adapter reflectsort --available lambda
```

Deps System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
