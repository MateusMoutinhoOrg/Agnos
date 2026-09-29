# `remove-adapter`

Uninstalls one adapter, leaving the contract it filled

```bash
agnos remove-adapter <Adapter> [--help] [--path <path>] [--quiet]
```

Removes adapters/libs/<adapter>/ and the require its declaration pins. Refuses an adapter an available still binds — point that available at another adapter first — and refuses one the generator wrote as the shim of a remote dep.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Adapter` | string, required |  | the adapter to remove from the project |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos remove-adapter reflectsort
```

Deps System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
