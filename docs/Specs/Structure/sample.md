# Project Structure

This document maps the project **schema** — the kinds of files a project is built from — not every concrete file. A slot with a **Spec** name is governed by a specification; resolve the name through [Specs](/docs/Specs/doc.md). A slot marked *generated* is written by `agnos build`.

```
adapters/  ──▶  sandbox/  ◀──  cmd/
(reaches the OS)  (closed)     (wires the two together)
```

## Root

| File | Description | Spec |
|------|-------------|------|
| `README.md` | Project overview and the Doc Index | Readme |
| `LICENSE` | License terms for the project | |
| `go.mod` | Go module definition | |

---

## `/sandbox/`
The closed sandbox. It imports nothing outside itself.

| File | Description | Spec |
|------|-------------|------|
| `new.go` | *generated* — `New(deps *deps.Deps) *api.Sandbox` running every binder | |

### `/sandbox/api/`

| File | Description | Spec |
|------|-------------|------|
| `sandbox.go` | *generated* — one field per other file of this directory | |
| `<contract>.go` | A struct of function fields the sandbox hands back | Contract |

### `/sandbox/binds/`

| File | Description | Spec |
|------|-------------|------|
| `<contract>.go` | `<Contract>Bind`, assigning the implementation onto the mirrored `api` struct | Binder |

### `/sandbox/deps/`

| File | Description | Spec |
|------|-------------|------|
| `deps.go` | *generated* — one `<Title> <dir>.Lib` field per sub-directory | |
| `<x>/<x>.go` | One `Lib` struct of function fields per sub-contract | DepsContract |

### `/sandbox/internal/commands/`

| File | Description | Spec |
|------|-------------|------|
| `<name>/entries.yaml` | The command's declaration | CommandEntries |
| `<name>/entries.go` | *generated* — the typed `Entries` struct | |
| `<name>/handler.go` | `CommandHandler(deps, entries) int` | CommandHandler |

---

## `/adapters/`
Outside the sandbox. The only place OS-bound and third-party code is allowed.

| File | Description | Spec |
|------|-------------|------|
| `libs/<lib>/<File>.go` | One package per sub-contract exporting `Bind(deps *deps.Deps)` | AdapterLib |
| `availables/standard/new.go` | *generated* — `New() deps.Deps` running every lib's `Bind` | |

---

## `/cmd/`

| File | Description | Spec |
|------|-------------|------|
| `main/main.go` | *generated* — wires `standard.New()` into `sandbox.New` and exits with `CliMain` | CliMain |
