# Extensions

An extension is one **generation mechanic**: a set of files `agnos` writes and
keeps up to date for you. `AgnosConfig/extensions.yaml` declares which ones are on, and it
is the whole of what `build` reads to decide what to render — nothing is inferred from the
directories the project happens to carry.

```yaml
backoffice: false
cli: true
database: false
deps: true
doc: true
example: true
front: false
readme: true
sandbox: true
server: false
```

| Key | What agnos generates and looks after |
|---|---|
| `sandbox` | the core: `sandbox/new.go`, `sandbox/api/sandbox.go`, `sandbox/constructors/<x>/constructor.go`, `sandbox/api/config.go`, `sandbox/internal/generated/config/new.go` |
| `deps` | `sandbox/deps/deps.go`, each binding's `new.go`; `add-dep` and the rest of the dependency commands |
| `cli` | `cmd/main`, `help`, `version`, `sandbox/api/{cli,command,trigger}.go` and the `OpinionatedAgnosCli` lib — the dispatch; `add-command` and the rest |
| `server` | `sandbox/internal/{server,routes}`, `sandbox/api/{server,route}.go` and the `OpinionatedAgnosServer` lib; `add-route` and the rest |
| `front` | the `OpinionatedAgnosFront` lib and the `front` route serving `assets/front/`; `add-page` and `remove-page` |
| `database` | `sandbox/internal/databases` and the `OpinionatedAgnosDatabase` lib; `add-database`, `add-table`, `add-table-field` and the rest; `sandbox/api/databaseconfig.go` and, with a cli, the `--database` middleware |
| `example` | the `examples/` suite; `add-cli-example`, `add-lib-example`, `run-examples`, `update-example` |
| `backoffice` | nothing on `build`: `backoffice-init` writes the admin backoffice once (routes, `backoffice-db`, the `backup` database and its snapshots, `add-backoffice-user`), and the key turns on its doc. Needs `server`, `front` and `database` |
| `doc` | the `docs/` tree and every `Index.md`; `add-doc` and `remove-doc` |
| `readme` | `README.md`, built from `AgnosConfig/docs/ReadmeHeader.md` and the doc index |

Everything that renders into the sandbox needs `sandbox` on. `doc` and `readme` stand on
their own.

## Turning one on and off

```bash
agnos list-extensions                  # every mechanic and whether it is on
agnos enable-extension readme          # start generating README.md again
agnos disable-extension doc            # stop generating docs/
```

`false` means **stop generating**, never **delete**. What the mechanic already wrote stays
exactly where it is and becomes yours, to keep or to edit by hand; the next `build` does not
touch it. Deleting those files is what an `<x>-purge` is for — and that is the same command
that writes the `false`:

```bash
agnos cli-init                         # cli: true, plus the deps the layer needs
agnos cli-purge                        # removes the files and writes cli: false
```

A mechanic that needs another one gets it: `server-init` runs `cli-init` first when the
project has no cli, `front-init` runs `server-init`, and `backoffice-init` runs whichever of
`server-init`, `front-init` and `database-init` the project is missing.

Never edit `extensions.yaml` by hand — the commands above re-render it with the keys in
alphabetical order. A key the catalog gained since this project was scaffolded is filled in
with its default by the next `build`, so a new mechanic never breaks an older tree.
