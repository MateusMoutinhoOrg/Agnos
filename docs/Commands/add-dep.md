# `add-dep`

Installs one dep of the embedded catalog into the project

```bash
agnos add-dep <Dep> [--adapter <adapter>] [--as <as>] [--remote-available <remote-available>] [--help] [--path <path>] [--quiet]
```

Renders the contract of assets/deplist/<dep> and the adapter that fills it, then calls build. The adapter is the dep's default-adapter unless --adapter names another; it is enrolled in every available.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Dep` | string, required |  | the dep to install from assets/deplist |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--adapter` | string |  | the adapter to fill the dep's contract with (defaults to the dep's default-adapter) | — |
| `--as` | string |  | the name the copied contract takes under sandbox/deps/ (remote deps only; defaults to the last segment of the module path) | — |
| `--remote-available` | string |  | the available of the remote repo the generated shim builds its sandbox from | — |
| `--help` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos add-dep embeddeps
agnos add-dep embeddeps --path ./my-project
```

Deps System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
