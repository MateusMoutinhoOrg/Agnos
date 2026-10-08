# `remove-page`

Remove an html page from assets/front/

```bash
agnos remove-page <Name> [--help] [--path <path>] [--quiet]
```

Deletes assets/front/<name>.html. The name is spelled as add-page spells it: the path under assets/front/ without .html.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the page to remove |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

```bash
agnos remove-page about
```

Front · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
