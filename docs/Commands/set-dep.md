# `set-dep`

Moves one remote dep to another version of its module

```bash
agnos set-dep <Dep> --version <version> [--remote-available <remote-available>] [--help] [--path <path>] [--quiet]
```

Re-copies the remote repo's sandbox/api into sandbox/deps/<dep>/ at the given version and regenerates the shim that converts it, then calls build. The module comes from the shim's own adapter.yaml.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Dep` | string, required |  | the remote dep to move |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--version` | string, required |  | the module version to copy the contract from | — |
| `--remote-available` | string |  | the available of the remote repo the regenerated shim builds its sandbox from | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos set-dep mathlib --version v1.3.0
```

Deps System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
