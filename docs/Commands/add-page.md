# `add-page`

Scaffold a new html page under assets/frontend/

```bash
agnos add-page [--path <path>] [--quiet] [--title <title>] <name>
```

Writes assets/frontend/<name>.html, plain html the frontend route serves as soon as it exists: /<name> answers it, and index is the page / answers. A name may carry slashes (blog/post). An existing file is refused. A page is a file and nothing else, so one written by hand or by a bundler is as much a page as one this command scaffolded.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |
| `--title` | string |  | the <title> the scaffolded page carries (defaults to the page name) |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the page, as its path under assets/frontend/ without .html (blog/post; index is the one / answers) |

```bash
agnos add-page about --title About
agnos add-page blog/post
```

Front System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
