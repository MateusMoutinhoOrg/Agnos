# `deps-init`

Initializes the dependency-injection subsystem for the project

```bash
agnos deps-init [--help] [--path <path>] [--quiet]
```

Creates the sandbox/deps and adapters directories and calls build. Run this once before using add-dep.

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
agnos deps-init
agnos deps-init --path ./my-project
```

Deps System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
