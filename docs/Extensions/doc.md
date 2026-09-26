# Extensions

An extension is one **generation mechanic**: a set of files `agnos` writes and
keeps up to date for you. `AgnosConfig/extensions.yaml` declares which ones are on, and it
is the whole of what `build` reads to decide what to render — nothing is inferred from the
directories the project happens to carry.

```yaml
doc: true
readme: true
sandbox: true
sandbox-cli: true
sandbox-database: false
sandbox-deps: true
sandbox-example: true
sandbox-front: false
sandbox-server: false
```

| Key | What agnos generates and looks after |
|---|---|
| `sandbox` | the core: `sandbox/new.go`, `sandbox/api/sandbox.go`, `sandbox/constructors/<x>/constructor.go`, `sandbox/api/config.go`, `sandbox/internal/generated/config/new.go` |
| `sandbox-deps` | `sandbox/deps/deps.go`, each available's `new.go`; `add-dep` and the rest of the dependency commands |
| `sandbox-cli` | `cmd/main`, the dispatch, `help`, `version`, `sandbox/api/{cli,command}.go`; `add-command` and the rest |
| `sandbox-server` | `sandbox/internal/{server,routeslist,routeio}`, `sandbox/api/{server,route}.go`; `add-route` and the rest |
| `sandbox-front` | `sandbox/internal/generated/frontio` and the `frontend` route serving `assets/frontend/`; `add-page` and `remove-page` |
| `sandbox-database` | `sandbox/internal/{databaseio,databases}`; `add-database`, `add-table`, `add-table-field` and the rest |
| `sandbox-example` | the `examples/` suite; `add-cli-example`, `add-lib-example`, `exec-test`, `update-test` |
| `doc` | the `docs/` tree and every `Index.md`; `add-doc` and `remove-doc` |
| `readme` | `README.md`, built from `AgnosConfig/docs/ReadmeHeader.md` and the doc index |

Everything that renders into the sandbox is spelled `sandbox-<mechanic>` and needs `sandbox`
on. `doc` and `readme` stand on their own.

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
agnos cli-init                         # sandbox-cli: true, plus the deps the layer needs
agnos cli-purge                        # removes the files and writes sandbox-cli: false
```

A mechanic that needs another one gets it: `server-init` runs `cli-init` first when the
project has no cli, and `front-init` runs `server-init`.

Never edit `extensions.yaml` by hand — the commands above re-render it with the keys in
alphabetical order. A key the catalog gained since this project was scaffolded is filled in
with its default by the next `build`, so a new mechanic never breaks an older tree.
