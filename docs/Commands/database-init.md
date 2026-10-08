# `database-init`

Add the database layer to the project

```bash
agnos database-init [--help] [--path <path>] [--quiet]
```

Installs the store the database layer is built over as a remote dep under sandbox/deps/databasedeps, the OpinionatedAgnosDatabase lib every generated methods.go reads it through, and turns the database mechanic on. It scaffolds no database of its own: which tables a project wants is a declaration, so 'agnos add-database' is the step that follows.

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
agnos database-init
agnos database-init --path ./my-project
```

Database · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
