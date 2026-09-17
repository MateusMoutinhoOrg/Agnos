# `set-dep`

Moves one remote dep to another version of its module

```bash
agnos set-dep --version <version> [--remote-available <remote-available>] [--path <path>] [--quiet] <dep>
```

Re-copies the remote repo's sandbox/api into sandbox/deps/<dep>/ at the given version and regenerates the shim that converts it, then calls build. The module comes from the shim's own adapter.yaml.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--version` | string, required |  | the module version to copy the contract from |
| `--remote-available` | string |  | the available of the remote repo the regenerated shim builds its sandbox from |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `dep` | string, required |  | the remote dep to move |

```bash
agnos set-dep mathlib --version v1.3.0
```

Deps System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
