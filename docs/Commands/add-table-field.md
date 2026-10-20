# `add-table-field`

Declare one field on a table of a database

```bash
agnos add-table-field <Name> --database <database> --table <table> [--parent <parent>] [--type <type>] [--required] [--target <target>] [--default <default>] [--fill-by-hand] [--help] [--path <path>] [--quiet]
```

Appends one field to a table of the database's database.yaml and runs build. What the field generates is decided by its type: a key brings a Find, a link a Get, an object the Add/List pair for the collection nested under each record, and every plain field a Set and a place in the table's filter. An optional field reads as empty on every record stored before it; a --required one opens a version of the data, whose migration fills the records already stored from --default or from the hand-written step --fill-by-hand writes, and refuses to finish while one holds no value.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the field name (the key the record stores it under) |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--database` | string, required |  | the database (identifier or package name) the field is declared on | — |
| `--table` | string, required |  | the table the field is declared on | — |
| `--parent` | string |  | the nested object field the new field goes inside, instead of the table itself | — |
| `--type` | string | `string` | the field type: key, string, integer, number, bytes, link or object | — |
| `--required` | boolean |  | an insert must carry this field | — |
| `--target` | string |  | the table a link points at (only with --type link) | — |
| `--default` | string |  | the value the migration gives every record already stored that holds none, for a --required field (never a key's) | — |
| `--fill-by-hand` | boolean |  | write the version's hand-written step, migrations/v<N>.go, to fill a --required field instead of a --default | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

```bash
agnos add-table-field alias --database app-database --table url --type key --required
```

Database · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
