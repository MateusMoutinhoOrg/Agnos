# `show-database`

Print one database's whole declaration as a tree

```bash
agnos show-database <Database> [--help] [--path <path>] [--quiet]
```

Reads sandbox/internal/databases/<database>/specs.yaml and prints it as a tree: the package it declares and where its keys are written, then every table, the fields it holds, and the signature of every method those fields generate. It is the one command of the database surface that writes nothing and runs no build.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Database` | string, required |  | the database (identifier or package name) to print |

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
agnos show-database app-database
```

Database System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
