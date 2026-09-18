# `add-table-field`

Declare one field on a table of a database

```bash
agnos add-table-field --database <database> --table <table> [--parent <parent>] [--type <type>] [--required] [--target <target>] [--path <path>] [--quiet] <name>
```

Appends one field to a table of the database's specs.yaml and runs build. What the field generates is decided by its type: a key brings a Find, a link a Get, a database the Add/List pair for the collection nested under each record, and every plain field an Update and a place in the table's filtrage.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--database` | string, required |  | the database (identifier or package name) the field is declared on |
| `--table` | string, required |  | the table the field is declared on |
| `--parent` | string |  | the nested database field the new field goes inside, instead of the table itself |
| `--type` | string | `string` | the field type: key, string, int, float, link or database |
| `--required` | boolean |  | an insert must carry this field |
| `--target` | string |  | the table a link points at (only with --type link) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the field name (the key the record stores it under) |

```bash
agnos add-table-field alias --database app-database --table url --type key --required
```

Database System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
