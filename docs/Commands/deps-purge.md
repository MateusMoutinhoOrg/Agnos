# `deps-purge`

Removes the dependency-injection subsystem from the project

```bash
agnos deps-purge [--path <path>] [--quiet]
```

Removes the sandbox/deps and adapters directories and calls build.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

```bash
agnos deps-purge
agnos deps-purge --path ./my-project
```

Deps System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
