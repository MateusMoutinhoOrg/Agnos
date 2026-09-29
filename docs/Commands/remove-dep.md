# `remove-dep`

Uninstalls one dep from the project

```bash
agnos remove-dep <Dep> [--with-adapters] [--help] [--path <path>] [--quiet]
```

Removes every adapter whose declaration names the dep, its require and its enrollment in every available, then the contract itself, then calls build.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Dep` | string, required |  | the dep to remove from the project |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--with-adapters` | boolean |  | also remove every adapter that fills the dep (without it, a dep with adapters installed is refused) | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos remove-dep embeddeps
agnos remove-dep embeddeps --path ./my-project
```

Deps System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
