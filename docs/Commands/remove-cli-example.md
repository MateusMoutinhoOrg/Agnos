# `remove-cli-example`

Delete an example from examples/cli/

```bash
agnos remove-cli-example <Name> [--help] [--path <path>] [--quiet]
```

Deletes examples/cli/<name>/ whole - the example.sh, the golden result.yaml and any TestDir or AssertDir the last run left behind - and runs build so the example listing of the docs is rewritten without it.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the example directory under examples/cli/ |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--help` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos remove-cli-example start
```

Examples · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
