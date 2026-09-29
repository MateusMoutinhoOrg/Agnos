# `remove-doc`

Delete a doc directory from docs/

```bash
agnos remove-doc <Name> [--help] [--path <path>] [--quiet]
```

Deletes docs/<name>/ (doc.md, props.yaml, its assets and every sub-doc nested under it) and runs build so the indexes that listed it are rewritten without it. A theme left with no docs simply stops rendering a section in README.md; it is not an error, so themes.yaml can keep it.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the doc directory under docs/, nested with / for a sub-doc (e.g. PublicApi/api.AddDoc) |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos remove-doc HandleReports
agnos remove-doc PublicApi/api.AddDoc --path ./my-project
```

Documentation · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
