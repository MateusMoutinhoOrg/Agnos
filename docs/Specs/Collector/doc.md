# Collector Specification

## Description
Defines the required shape of a **collector** — a function in `sandbox/internal/actions/build/collect_*.go` that turns one directory listing of the target project into the data a `{{range}}` in an asset template iterates over. Collectors are why a generated file's contents are predictable from a listing; see [BuildPipeline](/docs/BuildPipeline/doc.md#collectors). Adding one is [HandleAssetGroups](/docs/HandleAssetGroups/doc.md#add-a-variable-through-a-collector).

### Rules
- One file per collector, named `collect_<plural>.go`, holding one exported function `Collect<Plural>(io *smartio.SmartIO) []T` (`(deps, io)` when it must parse content).
- The one-line shape: list **one** directory with `io.ListDirs` or `io.ListFiles`; derive each name from the last path segment (dropping `.go` for files); title-case it; append; return the slice in listing order. Nothing is read, nothing is filtered beyond that.
- When a template needs more than a name, the element is a small struct or `map[string]any` with **precomputed** fields (`Title`, `Name`, `GoField`, `RangeCheck`) — the template does formatting, never logic.
- A collector that must parse content (`CollectCommands`, which reads each `entries.yaml` through `parsables/commandconf`) is the documented exception: it still lists one directory and still returns one element per entry, and it returns an `error` for an unparsable file naming that file.
- The result is added to the `vars` map in `build_internal.go` under the key the templates range over, and that key is documented in the variables table of [BuildPipeline](/docs/BuildPipeline/doc.md#template-variables).

## Structure
1. **Package clause**: `package build`.
2. **Imports**: `strings`, `sandbox/internal/smartio`, and — for the parsing exception — `deps` and the parsable.
3. **`Collect<Plural>`**: list, derive, title-case, return.

> **Note**: For a concrete example, refer to [sample.go](/docs/Specs/Collector/sample.go).
