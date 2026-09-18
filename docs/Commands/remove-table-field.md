# `remove-table-field`

Delete one declared field from a table

```bash
agnos remove-table-field --database <database> --table <table> [--parent <parent>] [--path <path>] [--quiet] <name>
```

Drops one field from a table of the database's specs.yaml and runs build without compiling: the methods that field generated go with it, and hand-written code may still be calling one.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--database` | string, required |  | the database (identifier or package name) the field is declared on |
| `--table` | string, required |  | the table the field is declared on |
| `--parent` | string |  | the nested database field the field sits inside, instead of the table itself |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the field to delete, by the name it is declared under |

```bash
agnos remove-table-field redirects --database app-database --table url
```

Database System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
