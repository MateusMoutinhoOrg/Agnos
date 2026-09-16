# Adapters

Three units, and the relation between them is declared, never inferred.

| Unit | Is | Lives in | How many |
|---|---|---|---|
| **dep** | the contract: one field of `deps.Deps`, reached as `sandbox.Deps.<Dep>` | `sandbox/deps/<dep>/` (closed) | one per field |
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
the two apart — it calls `sandbox.Deps.Sortdeps` either way — so the choice lives entirely in the
entry point that picks an available.

## From another agnos repo

Every agnos repo is installable as a dep by construction: nothing is declared
and nothing is turned on. The repo publishes two halves, and only one of them
is copied.

| Half | Is | Becomes, in the consumer |
|---|---|---|
| contract | `sandbox/api/`, which imports nothing but `sandbox/deps`, by the rule every agnos repo lives under | `sandbox/deps/<name>/`, the same files with the package clause rewritten and `Sandbox.Deps` dropped |
| wiring | `Sandbox.Deps`, how the repo reaches the outside world | nothing — a consumer installs an api, never the wiring behind it |
| adapter | `sandbox.New` over one of its own availables | nothing — it runs compiled, out of the remote module |

```bash
{{.GeneratorName}} add-dep github.com/user/MathLib@v1.2.0 --as mathlib
{{.GeneratorName}} set-dep mathlib --version v1.3.0     # re-copy and regenerate
{{.GeneratorName}} remove-dep mathlib --with-adapters
```

An argument holding a `/` is a module path, one without is a name of the
catalogue — the same disambiguation `go get` makes.

The adapter that fills the copied contract is **generated**, and it has to be:
`mathlib.Sandbox` lives at an import path inside the consumer's module, one the
remote repo does not know when it publishes, so any Go that names that type has
to compile in the consumer. A top-level cast will not stand in for it either —
Go's type identity does not reach through named types, so two copies of a struct
whose fields are named types are never the same type, however identical the
source. The generated shim builds the remote sandbox out of the remote repo's
own adapters and converts the result field by field, with the parameters of a
`func` field converted in the opposite direction.

**The convertibility rule** is one criterion — the underlying type has to be
identical in both copies — and every case follows from it:

| Form | Treatment |
|---|---|
| identical underlying: a builtin, a slice/map/pointer of builtins, `any`, `error`, a `func` of builtins, a named interface of builtin-only methods | assignment or a plain conversion |
| a struct with a field of a named type of the same package | a generated converter, field by field — the only case that generates code |
| a `func` field whose parameters or results fall in the case above | a closure, parameters in the reverse direction |
| a slice, map or pointer of a convertible named type | a generated loop |
| a type of another package (`time.Time`, `io.Reader`) | impossible: the api imports nothing but `sandbox/deps`, and `Sandbox.Deps` — the one field naming it — never crosses |
| generics, `chan`, an anonymous struct or interface, an embedded field | rejected — outside the generator, not outside Go |

This is not a new rule: it is the discipline `sandbox/deps/` contracts already
state — only builtin types cross the boundary — extended to `sandbox/api/` and
made checkable. `{{.GeneratorName}} verify` applies it to every repo, always and
with no flag: what varies between repos is the shape of an api, and a shape is
verified rather than declared. An opt-in defaulting to off would only move the
violation to the consumer's install, which is what the check exists to prevent.

The copy is checked against the module cache on every `verify`, byte for byte,
whenever that module is already on the machine — the same rule
`assets/deplist/` lives under, with the cache in place of the assets. No digest
is stored: the cache is immutable per version and `go.sum` already signs it.

## Removing

| Command | Refuses when |
|---|---|
| `{{.GeneratorName}} remove-adapter <adapter>` | an available still binds it (point that available elsewhere first), or the generator wrote it as the shim of a remote dep |
| `{{.GeneratorName}} remove-dep <dep>` | an adapter still fills it — `--with-adapters` takes them all |
| `{{.GeneratorName}} remove-available <name>` | it is `standard`: `cmd/main/main.go` imports it |

Both refusals name what is holding the unit, so the answer is in the message.
