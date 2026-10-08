# `add-page`

Scaffold a new html page under assets/front/

```bash
agnos add-page <Name> [--title <title>] [--help] [--path <path>] [--quiet]
```

Writes assets/front/<name>.html, plain html the front route serves as soon as it exists: /<name> answers it, and index is the page / answers. A name may carry slashes (blog/post). An existing file is refused. A page is a file and nothing else, so one written by hand or by a bundler is as much a page as one this command scaffolded.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the page, as its path under assets/front/ without .html (blog/post; index is the one / answers) |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--title` | string |  | the <title> the scaffolded page carries (defaults to the page name) | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

```bash
agnos add-page about --title About
agnos add-page blog/post
```

Front · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
