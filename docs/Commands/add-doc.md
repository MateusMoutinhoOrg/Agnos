# `add-doc`

Scaffold a new doc directory under docs/

```bash
agnos add-doc [--theme <theme>...] --description <description> [--path <path>] [--quiet] <name>
```

Creates docs/<name>/ with a doc.md stub and the props.yaml declaring it, then runs build so README.md's index and the parent's Index.md list it. A first-level doc needs at least one --theme of themes.yaml; a nested name (docs/<Parent>/<Name>) creates a sub-doc, which takes no theme. Refuses to overwrite an existing doc.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--theme`, `-t` | string, repeatable |  | a theme id of themes.yaml the doc belongs to (repeatable; first-level docs only) |
| `--description`, `-d` | string, required |  | the one-line summary every index lists the doc with |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the doc directory under docs/, nested with / for a sub-doc (e.g. PublicApi/api.AddDoc) |

```bash
agnos add-doc HandleReports --theme development --description "How a report is written and regenerated"
agnos add-doc PublicApi/api.AddDoc --description "The AddDoc action of the sandbox api"
```

Documentation · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
