# Parsable Specification

## Description
Defines the required shape of a **parsable** — a package under `sandbox/internal/parsables/<name>conf/` that parses one configuration file into a struct, exposes operations on it as function fields, and renders it back. Adding one is [HandleParsables.md](/docs/Tutorials/HandleParsables.md).

### Rules
- `package <name>conf`, five files with fixed names: `api.go`, `new.go`, `new_empty.go`, `bind_methods.go`, `render.go`. Further operations go in `<topic>_methods.go` files (`moduleconf/require_methods.go`); a second constructor in `new_<source>.go` (`moduleconf/new_from_path.go`).
- `api.go` declares one exported struct `<Name>Conf`: the data fields first, then the operations as **function fields**, `Render func() string` last. Supporting structs (`Field`, `Theme`, `PathReplacerEntry`) are declared above it.
- `new.go` holds `New(deps *deps.Deps, content string) (*<Name>Conf, error)`: parse through `deps.Serializables.ParseYaml` (or by hand, for a non-YAML format), fill the fields, call `bindMethods`, return. A malformed file is an error naming what was wrong.
- `new_empty.go` holds `NewEmpty(deps *deps.Deps) *<Name>Conf` returning the defaults with methods bound.
- `bind_methods.go` holds the unexported `bindMethods(deps, self)` assigning **every** function field as a closure over `self`.
- `render.go` holds the unexported `render(deps, self) string` building a value tree and calling `deps.Serializables.SerializeToYaml` (or writing text by hand for a non-YAML format). `New(Render())` must yield an equal struct.
- The package imports `sandbox/deps` and the standard library; never a YAML or JSON module.

## Structure
1. **`api.go`**: supporting structs, then `type <Name>Conf struct` — data, operations, `Render`.
2. **`new.go`**: `New`.
3. **`new_empty.go`**: `NewEmpty`.
4. **`bind_methods.go`**: `bindMethods`.
5. **`render.go`**: `render`.

> **Note**: For a concrete example, refer to [sample.go](/docs/References/Specs/Parsable/sample.go), which shows the five files in one listing.
