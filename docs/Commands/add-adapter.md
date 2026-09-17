# `add-adapter`

Installs one further adapter for a contract the project already has

```bash
agnos add-adapter [--available <available>] [--path <path>] [--quiet] <adapter>
```

Renders assets/adapterlist/<adapter> into the project and writes its declaration to adapters/libs/<adapter>/adapter.yaml. The contract it fills has to be installed already. Installing changes no selection: an available binds one adapter per field, so --available names the one that switches to it.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--available` | string |  | the available that should switch to this adapter (installs only when absent) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `adapter` | string, required |  | the adapter to install from assets/adapterlist |

```bash
agnos add-adapter reflectsort
agnos add-adapter reflectsort --available lambda
```

Deps System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
