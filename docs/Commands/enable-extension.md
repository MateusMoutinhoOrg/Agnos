# `enable-extension`

Turn one generation mechanic on

```bash
agnos enable-extension <Name> [--help] [--path <path>] [--quiet]
```

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the extension to turn on |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--help` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

Extensions · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
