# `remove-table-field`

Delete one declared field from a table

```bash
agnos remove-table-field <Name> --database <database> --table <table> [--parent <parent>] [--help] [--path <path>] [--quiet]
```

Drops one field from a table of the database's database.yaml and runs build without compiling: the methods that field generated go with it, and hand-written code may still be calling one.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the field to delete, by the name it is declared under |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--database` | string, required |  | the database (identifier or package name) the field is declared on | — |
| `--table` | string, required |  | the table the field is declared on | — |
| `--parent` | string |  | the nested object field the field sits inside, instead of the table itself | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

```bash
agnos remove-table-field redirects --database app-database --table url
```

Database · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
