# `show-database`

Print one database's whole declaration as a tree

```bash
agnos show-database <Name> [--help] [--path <path>] [--quiet]
```

Reads sandbox/internal/databases/<database>/database.yaml and prints it as a tree: the package it declares and where its keys are written, then every table, the fields it holds, and the signature of every method those fields generate. It is the one command of the database surface that writes nothing and runs no build.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the database (identifier or package name) to print |

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
agnos show-database app-database
```

Database · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
