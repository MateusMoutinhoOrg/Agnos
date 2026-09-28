# `set-table-field`

Rewrite one declared field of a table

```bash
agnos set-table-field <Name> --database <database> --table <table> [--parent <parent>] [--rename <rename>] [--type <type>] [--required] [--target <target>] [--clear <clear>...] [--help] [--path <path>] [--quiet]
```

Rewrites one declared field in place and runs build. It is add-table-field applied to a declaration that already exists: the keys given are written over the ones there, --clear takes one off, and the result goes through the same constructor — so changing a type never means removing the field and declaring it again.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the field to edit, by the name it is declared under |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--database` | string, required |  | the database (identifier or package name) the field is declared on | — |
| `--table` | string, required |  | the table the field is declared on | — |
| `--parent` | string |  | the nested database field the field sits inside, instead of the table itself | — |
| `--rename` | string |  | the name the field is stored under from now on | — |
| `--type` | string |  | the field type it takes on: key, string, int, float, link or database | — |
| `--required` | boolean |  | an insert must carry this field | — |
| `--target` | string |  | the table a link points at (only with --type link) | — |
| `--clear` | string-array |  | a key to take off again: required or target (repeatable) | — |
| `--help` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos set-table-field link --database app-database --table url --required
```

Database System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
