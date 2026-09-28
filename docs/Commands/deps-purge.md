# `deps-purge`

Removes the dependency-injection subsystem from the project

```bash
agnos deps-purge [--help] [--path <path>] [--quiet]
```

Removes the sandbox/deps and adapters directories and calls build.

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
agnos deps-purge
agnos deps-purge --path ./my-project
```

Deps System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
