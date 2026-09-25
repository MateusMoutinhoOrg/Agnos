# `show-route`

Print one route's whole declaration as a tree

```bash
agnos show-route [--path <path>] [--quiet] <route>
```

Reads sandbox/internal/routeslist/<route>/route.yaml and prints it as a tree: the request line the route answers, then every place the declaration holds something — its paths, its parameters and the json-schema of its body, property by property with the keywords declared on each. It is the one command of the route surface that writes nothing and runs no build.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `route` | string, required |  | the route (identifier or package name) to print |

```bash
agnos show-route create-user
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
