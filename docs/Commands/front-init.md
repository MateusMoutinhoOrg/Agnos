# `front-init`

Add the html front layer to the project

```bash
agnos front-init [--help] [--path <path>] [--quiet]
```

Installs the deps the front layer needs — the OpinatedAgnosFront lib, its file layer, among them — and writes, once, the frontend route serving every file of assets/frontend and that tree's index.html and 404.html. A project with no server layer is given one first: the front is answered over http.

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
agnos front-init
agnos front-init --path ./my-project
```

Front System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
