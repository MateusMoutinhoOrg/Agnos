# `set-dep`

Move one remote dep to another version of its module

```bash
agnos set-dep <Name> --version <version> [--remote-binding <remote-binding>] [--help] [--path <path>] [--quiet]
```

Re-copies the remote repo's sandbox/api into sandbox/deps/<dep>/ at the given version and regenerates the shim that converts it, then calls build. The module comes from the shim's own adapter.yaml.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the remote dep to move |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--version` | string, required |  | the module version to copy the contract from | — |
| `--remote-binding` | string |  | the binding of the remote repo the regenerated shim builds its sandbox from | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

```bash
agnos set-dep mathlib --version v1.3.0
```

Deps · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
