# `show-database`

Print one database's whole declaration as a tree

```bash
agnos show-database [--path <path>] [--quiet] <database>
```

Reads sandbox/internal/databases/<database>/specs.yaml and prints it as a tree: the package it declares and where its keys are written, then every table, the fields it holds, and the signature of every method those fields generate. It is the one command of the database surface that writes nothing and runs no build.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `database` | string, required |  | the database (identifier or package name) to print |

```bash
agnos show-database app-database
```

Database System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
