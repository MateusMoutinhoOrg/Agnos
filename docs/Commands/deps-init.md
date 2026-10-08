# `deps-init`

Initialize the dependency-injection subsystem for the project

```bash
agnos deps-init [--help] [--path <path>] [--quiet]
```

Creates the sandbox/deps and adapters directories and calls build. Run this once before using add-dep.

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
agnos deps-init
agnos deps-init --path ./my-project
```

Deps · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
