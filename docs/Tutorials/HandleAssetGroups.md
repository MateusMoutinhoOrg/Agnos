# Handle Asset Groups and Collectors

## Description
Covers changing what `agnos build` renders: adding a template to an asset group under `assets/<group>/**`, and feeding it a new variable through a collector under `sandbox/internal/actions/build/collect_*.go`. Which group renders when, and in what order, is [BuildPipeline.md](/docs/References/BuildPipeline.md); the shape a template and a collector must have is the AssetTemplate and Collector specifications in [Specs.md](/docs/References/Specs.md).

### Rules
- A template lives at the path it will be written to inside the target project: `assets/all/sandbox/new.go` renders to `sandbox/new.go`. The group is the first path segment and nothing else marks it.
- Every file in a group is rendered with the **same** `vars` map, built once by `BuildInternal`. A template may use any key of it: `Module`, `Name`, `Description`, `Version`, `ProjectName`, `HasDeps`, `HasCli`, `Binds`, `Constructors`, `DepsLibs`, `AdapterLibs`, `Commands`.
- Templates are rendered through `io.WriteFileOverwrite`: every file of a group is regenerated on every build, so nothing hand-written may live at a group's path.
- A collector lists one directory with `io.List*`, derives names from the last path segment, title-cases them, and returns the slice. `CollectCommands` is the one exception, because it parses each `entries.yaml`.
- After changing a template or a collector, regenerate this tree with the bootstrap binary — see [BootstrapAgnos.md](/docs/Tutorials/BootstrapAgnos.md). `assets/` is never compiled, so a broken template only fails there.

---

## Workflow

### Add a template to a group

1. Decide the group by when the file must exist: `all` for every project, `deps` for projects with `sandbox/deps/`, `cli` for projects with `sandbox/internal/cli/`, `start` for the configuration `agnos start` writes once.
2. Create the file at its target path inside the group, as a Go `text/template`, and take the module path from `{{.Module}}`:
   ```go
   // assets/all/sandbox/internal/config/build_info.go
   package config

   // Rendered by `agnos build` — do not edit by hand.
   const Description = "{{.Description}}"
   ```
3. Bootstrap and check the file appears in this tree and in a scratch project:
   ```bash
   go build -o release/bootstrap.bin ./cmd/main && ./release/bootstrap.bin build
   cat sandbox/internal/config/build_info.go
   ```
4. Register the file in [GeneratedFiles.md](/docs/References/GeneratedFiles.md) and in the group's row of [Structure.md](/docs/References/Structure.md#assets).

### Add a variable through a collector

5. Write the collector beside the others, following the one-line shape:
   ```go
   // sandbox/internal/actions/build/collect_parsables.go
   package build

   import (
       "strings"

       "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
   )

   // CollectParsables lists sandbox/internal/parsables/<x>/ and returns the
   // title-cased package names, one per dir, for {{range .Parsables}}.
   func CollectParsables(io *smartio.SmartIO) []string {
       var names []string
       for _, dir := range io.ListDirs("sandbox/internal/parsables") {
           parts := strings.Split(dir, "/")
           name := parts[len(parts)-1]
           names = append(names, strings.ToUpper(name[:1])+name[1:])
       }
       return names
   }
   ```
6. Add its result to the `vars` map in `build_internal.go`, under the key the templates will range over:
   ```go
   "Parsables": CollectParsables(io),
   ```
7. Use it from a template with `{{range .Parsables}}…{{end}}`, bootstrap, and prove idempotence with a second build that changes nothing.
8. Document the new key in the variables table of [BuildPipeline.md](/docs/References/BuildPipeline.md#template-variables).

### Add a single-file scaffold

9. Put a template that is rendered to **one** destination outside any group, under `assets/templates/`, and render it with `utils.RenderTemplateToDest(deps, io, "templates/<file>", vars, dest)` from the action that owns it — `command_entries.yaml` and `command_handler.go` (rendered by `add-command`) and `entries.go` (rendered per command by `build`) are the existing ones.
