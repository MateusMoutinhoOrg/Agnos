# `start`

Initialize a new project in a directory

```bash
agnos start --project-name <project-name> [--force] [--module <module>] [--help] [--path <path>] [--quiet]
```

Scaffolds a new Agnos project in the given directory, creating the required configuration files and folder structure. If no path is provided, the current directory is used.

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--project-name`, `-p` | string, required |  | the name of the project | — |
| `--force`, `-f` | boolean |  | Forces the creation of the project, overwriting existing files | — |
| `--module`, `-m` | string |  | the go module path written into go.mod (required when the target dir has no go.mod yet) | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

```bash
agnos start -p my-project -m github.com/you/my-project
agnos start -p my-project -m github.com/you/my-project --path ./my-project-dir
agnos start -p my-project -m github.com/you/my-project -q
```

Core · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
