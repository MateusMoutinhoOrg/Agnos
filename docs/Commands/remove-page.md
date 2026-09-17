# `remove-page`

Remove an html page

```bash
agnos remove-page [--path <path>] [--quiet] <name>
```

Deletes the page's route package and its html template both. A route with no html beside it is not a page: remove-route is the editor for those.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the page to remove |

```bash
agnos remove-page about
```

Front System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
