# `remove-page`

Remove an html page from assets/frontend/

```bash
agnos remove-page <Name> [--help] [--path <path>] [--quiet]
```

Deletes assets/frontend/<name>.html. The name is spelled as add-page spells it: the path under assets/frontend/ without .html.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the page to remove |

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
agnos remove-page about
```

Front System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
