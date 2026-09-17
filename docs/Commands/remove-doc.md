# `remove-doc`

Delete a doc directory from docs/

```bash
agnos remove-doc [--path <path>] [--quiet] <name>
```

Deletes docs/<name>/ (doc.md, props.yaml, its assets and every sub-doc nested under it) and runs build so the indexes that listed it are rewritten without it. A theme left with no docs simply stops rendering a section in README.md; it is not an error, so themes.yaml can keep it.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the doc directory under docs/, nested with / for a sub-doc (e.g. PublicApi/api.AddDoc) |

```bash
agnos remove-doc HandleReports
agnos remove-doc PublicApi/api.AddDoc --path ./my-project
```

Documentation · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
