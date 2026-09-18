# `database-init`

Add the database layer to the project

```bash
agnos database-init [--path <path>] [--quiet]
```

Installs the store the database layer is built over as a remote dep under sandbox/deps/database, renders sandbox/internal/databaseio and turns the sandbox-database mechanic on. It scaffolds no database of its own: which tables a project wants is a declaration, so 'agnos add-database' is the step that follows.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

```bash
agnos database-init
agnos database-init --path ./my-project
```

Database System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
