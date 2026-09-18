# `set-table-field`

Rewrite one declared field of a table

```bash
agnos set-table-field --database <database> --table <table> [--parent <parent>] [--rename <rename>] [--type <type>] [--required] [--target <target>] [--clear <clear>...] [--path <path>] [--quiet] <name>
```

Rewrites one declared field in place and runs build. It is add-table-field applied to a declaration that already exists: the keys given are written over the ones there, --clear takes one off, and the result goes through the same constructor — so changing a type never means removing the field and declaring it again.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--database` | string, required |  | the database (identifier or package name) the field is declared on |
| `--table` | string, required |  | the table the field is declared on |
| `--parent` | string |  | the nested database field the field sits inside, instead of the table itself |
| `--rename` | string |  | the name the field is stored under from now on |
| `--type` | string |  | the field type it takes on: key, string, int, float, link or database |
| `--required` | boolean |  | an insert must carry this field |
| `--target` | string |  | the table a link points at (only with --type link) |
| `--clear` | string, repeatable |  | a key to take off again: required or target (repeatable) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the field to edit, by the name it is declared under |

```bash
agnos set-table-field link --database app-database --table url --required
```

Database System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
