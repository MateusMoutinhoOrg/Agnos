# `remove-adapter`

Uninstalls one adapter, leaving the contract it filled

```bash
agnos remove-adapter [--path <path>] [--quiet] <adapter>
```

Removes adapters/libs/<adapter>/ and the require its declaration pins. Refuses an adapter an available still binds — point that available at another adapter first — and refuses one the generator wrote as the shim of a remote dep.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `adapter` | string, required |  | the adapter to remove from the project |

```bash
agnos remove-adapter reflectsort
```

Deps System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
