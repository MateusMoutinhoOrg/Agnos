# `backoffice-init`

Add the admin backoffice to the project

```bash
agnos backoffice-init [--help] [--path <path>] [--quiet]
```

Writes the admin backoffice into the project: the /admin pages (login, home, backoffice users, API tokens), the /api/admin JSON api, the backoffice-client-ip, backoffice-security-headers and backoffice-same-origin middlewares, the backoffice-db database, the add-backoffice-user command and the backoffice-start-server middleware in front of start-server. The server, front and database layers are installed first when any is missing, and so are the catalog deps the backoffice calls into (envdeps, jwtdeps, passworddeps, randdeps, ratelimitdeps, timedeps). Every file it writes is the project's from then on: a second backoffice-init keeps each one already there. Nothing the project wrote is edited: the backoffice's part of RouteProps and of api.Config are files of their own. start-server then reads the session secret from <NAME>_BACKOFFICE_SECRET, the project's name upper-cased (MEUSITE_BACKOFFICE_SECRET for meusite), and generates one for the run when it is unset.

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
agnos backoffice-init
agnos backoffice-init --path ./my-project
```

Backoffice · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
