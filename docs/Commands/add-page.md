# `add-page`

Declare a new html page

```bash
agnos add-page [--path <path>] [--quiet] [--trigger <trigger>] [--title <title>] [--help <help>] <name>
```

Declares the route that answers the page and writes the html template it renders under assets/frontend/pages/. The trigger defaults to /<name>; --trigger / declares the home page. An html file already there is kept, which is the way back from a front-purge.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |
| `--trigger` | string |  | the literal segment the page answers on, / included (defaults to /<name>) |
| `--title` | string |  | the <title> the scaffolded page carries (defaults to the page name) |
| `--help` | string |  | one-line description of the page, for docs/Routes (defaults to one derived from the name) |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the page name (becomes the route, its Go package and the html file) |

```bash
agnos add-page home --trigger / --title Home
agnos add-page about --title About
```

Front System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
