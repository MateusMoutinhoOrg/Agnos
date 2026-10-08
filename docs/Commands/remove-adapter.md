# `remove-adapter`

Uninstall one adapter, leaving the contract it filled

```bash
agnos remove-adapter <Name> [--help] [--path <path>] [--quiet]
```

Removes adapters/impls/<adapter>/ and the require its declaration pins. Refuses an adapter a binding still binds — point that binding at another adapter first — and refuses one the generator wrote as the shim of a remote dep.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the adapter to remove from the project |

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
agnos remove-adapter reflectsort
```

Deps · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
