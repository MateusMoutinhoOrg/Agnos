# `remove-database`

Delete one database package whole

```bash
agnos remove-database [--path <path>] [--quiet] <name>
```

Deletes sandbox/internal/databases/<name>/ and runs build without compiling: dropping a database may leave hand-written code calling what is gone. It refuses a package carrying a methods_custom.go, which is the project's own and which nothing in the declaration describes.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the database (identifier or package name) to delete |

```bash
agnos remove-database app-database
```

Database System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
