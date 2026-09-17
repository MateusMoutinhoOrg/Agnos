# `start`

Initialize a new project in a directory

```bash
agnos start [--path <path>] --project-name <project-name> [--quiet] [--force] [--module <module>]
```

Scaffolds a new Agnos project in the given directory, creating the required configuration files and folder structure. If no path is provided, the current directory is used.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--project-name`, `-p` | string, required |  | the name of the project |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |
| `--force`, `-f` | boolean |  | Forces the creation of the project, overwriting existing files |
| `--module`, `-m` | string |  | the go module path written into go.mod (required when the target dir has no go.mod yet) |

```bash
agnos start -p my-project
agnos start -p my-project --path ./my-project-dir
agnos start -p my-project -q
```

Core Commands · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
