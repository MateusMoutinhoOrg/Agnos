# `add-dep`

Installs one dep of the embedded catalog into the project

```bash
agnos add-dep [--adapter <adapter>] [--path <path>] [--quiet] [--as <as>] [--remote-available <remote-available>] <dep>
```

Renders the contract of assets/deplist/<dep> and the adapter that fills it, then calls build. The adapter is the dep's default-adapter unless --adapter names another; it is enrolled in every available.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--adapter` | string |  | the adapter to fill the dep's contract with (defaults to the dep's default-adapter) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |
| `--as` | string |  | the name the copied contract takes under sandbox/deps/ (remote deps only; defaults to the last segment of the module path) |
| `--remote-available` | string |  | the available of the remote repo the generated shim builds its sandbox from |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `dep` | string, required |  | the dep to install from assets/deplist |

```bash
agnos add-dep embeddeps
agnos add-dep embeddeps --path ./my-project
```

Deps System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
