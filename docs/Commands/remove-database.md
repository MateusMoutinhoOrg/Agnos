# `remove-database`

Delete one database package whole

```bash
agnos remove-database <Name> [--help] [--path <path>] [--quiet]
```

Deletes sandbox/internal/databases/<name>/ and runs build without compiling: dropping a database may leave hand-written code calling what is gone. It refuses a package carrying a methods_custom.go, which is the project's own and which nothing in the declaration describes.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the database (identifier or package name) to delete |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--help` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos remove-database app-database
```

Database System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
