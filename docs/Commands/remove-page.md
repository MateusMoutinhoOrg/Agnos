# `remove-page`

Remove an html page from assets/frontend/

```bash
agnos remove-page [--path <path>] [--quiet] <name>
```

Deletes assets/frontend/<name>.html. The name is spelled as add-page spells it: the path under assets/frontend/ without .html.

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
