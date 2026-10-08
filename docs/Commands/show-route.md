# `show-route`

Print one route's whole declaration as a tree

```bash
agnos show-route <Name> [--help] [--path <path>] [--quiet]
```

Reads the route's route.yaml, in whatever folder of sandbox/internal/routes it sits, and prints it as a tree: the request line the route answers, then every place the declaration holds something — its paths, its parameters and the json-schema of its body, property by property with the keywords declared on each. It is the one command of the route surface that writes nothing and runs no build.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the route (identifier or package name) to print |

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
agnos show-route create-user
```

Server · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
