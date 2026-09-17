# `deps-init`

Initializes the dependency-injection subsystem for the project

```bash
agnos deps-init [--path <path>] [--quiet]
```

Creates the sandbox/deps and adapters directories and calls build. Run this once before using add-dep.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

```bash
agnos deps-init
agnos deps-init --path ./my-project
```

Deps System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
