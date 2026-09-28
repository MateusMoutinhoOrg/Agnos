# `publish`

Builds, compiles and publishes a release via gh

```bash
agnos publish [--release-name <release-name>] [--draft] [--target <target>] [--publisher <publisher>] [--help] [--path <path>] [--quiet]
```

Runs build, then compile (every target by default), and publishes every file of release/ as a gh release named --release-name, defaulting to the version in AgnosConfig/project.yaml.

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--release-name`, `-rn` | string |  | The name of the release | — |
| `--draft` | boolean |  | Create a draft release | — |
| `--target`, `-t` | string | `all` | The target to compile for (defaults to all) | — |
| `--publisher`, `-pub` | string | `gh` | The publisher to use (defaults to gh) | — |
| `--help` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

Core Commands · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
