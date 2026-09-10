# Adapters

Three units, and the relation between them is declared, never inferred.

| Unit | Is | Lives in | How many |
|---|---|---|---|
| **dep** | the contract: one field of `deps.Deps` | `sandbox/deps/<dep>/` (closed) | one per field |
| **adapter** | one implementation of it, exporting `Bind` | `adapters/libs/<adapter>/` | any number per dep |
| **available** | a selection: exactly one adapter per dep | `adapters/availables/<name>/` | any number per project |

`adapters/libs/` is what the project **has**. `adapters/availables/<name>/available.yaml` is
which of them **wins** for each field. `cmd/main/main.go` imports one available — `standard` —
and that import is the whole of how a program picks its implementations.

## The invariant

**Every available fills every field of `Deps` exactly once.** Zero leaves a nil func that
panics on first use; two is a silent overwrite in which whichever bound last wins. `verify`
reports both, and names them differently. What each adapter fills is read from its own
`adapters/libs/<adapter>/adapter.yaml`:

```yaml
dep: sortdeps
help: Insertion sort over reflect, no stdlib sort
module: ""
name: reflectsort
origin: catalog
```

`module` is the versioned module that adapter imports — `""` for one that needs nothing beyond
the stdlib — and it is what `add-dep` and `add-adapter` put in `go.mod`, filed under the
adapter that actually imports it. `origin` is `catalog` for one the catalogue installs.

`adapters/availables/<name>/available.yaml` lists the winners, and `new.go` beside it is
generated from that list:

```yaml
adapters:
    - reflectsort
    - std
```

An available directory with no `available.yaml` is a hand-written mix, and no build touches it.

## Two implementations of one contract

```bash
{{.GeneratorName}} add-dep sortdeps                          # contract + its default-adapter, bound everywhere
{{.GeneratorName}} add-adapter reflectsort                   # a second implementation, bound nowhere yet
{{.GeneratorName}} add-available lambda                      # a second selection, a copy of standard's
{{.GeneratorName}} set-adapter sortdeps reflectsort --available lambda
{{.GeneratorName}} list-adapters                             # who is installed, and who binds whom
```

`standard` still binds `sortdeps`; `lambda` binds `reflectsort`. Nothing in `sandbox/` can tell
the two apart — it calls `deps.Sortdeps` either way — so the choice lives entirely in the
entry point that picks an available.

## Removing

| Command | Refuses when |
|---|---|
| `{{.GeneratorName}} remove-adapter <adapter>` | an available still binds it (point that available elsewhere first), or the generator wrote it as the shim of a remote dep |
| `{{.GeneratorName}} remove-dep <dep>` | an adapter still fills it — `--with-adapters` takes them all |
| `{{.GeneratorName}} remove-available <name>` | it is `standard`: `cmd/main/main.go` imports it |

Both refusals name what is holding the unit, so the answer is in the message.
