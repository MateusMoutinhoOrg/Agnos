# Add a Parsable Config

## Description
Covers adding a package under `sandbox/internal/parsables/<name>conf/` — a small parser for one YAML file (or `go.mod`) an action reads and writes back — following the fixed five-file shape every existing parsable has. The files `agnos start` writes and `agnos build` reads back are listed in [Structure.md](/docs/References/Structure.md#agnosconfig); the shape is the Parsable specification in [Specs.md](/docs/References/Specs.md).

### Rules
- Five files, always the same names: `api.go` (the struct, with its method fields), `new.go` (`New(deps, content) (*T, error)` parsing a string), `new_empty.go` (`NewEmpty(deps) *T` with defaults), `bind_methods.go` (assigns every method field), `render.go` (the `Render` implementation). A parsable with more methods adds `<topic>_methods.go` files, as `moduleconf/require_methods.go` does.
- The struct's operations are **function fields** assigned by `bind_methods.go`, not Go methods — the same struct-contract shape as everything else; see [StructContracts.md](/docs/References/StructContracts.md).
- Parsing goes through `deps.Serializables` (`ParseYaml`, `SerializeToYaml`) so the sandbox imports no YAML library. `moduleconf` is the one parsable that parses by hand, because `go.mod` is not YAML.
- `Render` must reproduce a file that `New` parses back into an equal struct: an action loads, mutates and renders, and the round trip is what keeps `entries.yaml` editable from the command line.

---

## Workflow
1. Create the package and declare the struct in `api.go`. Data fields first, then the operations as function fields, `Render` last:
   ```go
   // sandbox/internal/parsables/hooksconf/api.go
   package hooksconf

   type HooksConf struct {
       Before []string
       After  []string

       AddBefore func(command string)
       AddAfter  func(command string)
       Render    func() string
   }
   ```
2. Write `new_empty.go` — the defaults, with the methods bound:
   ```go
   package hooksconf

   import "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"

   func NewEmpty(deps *deps.Deps) *HooksConf {
       self := &HooksConf{}
       bindMethods(deps, self)
       return self
   }
   ```
3. Write `new.go` — parse a string through the injected serializer, fill the struct, bind the methods:
   ```go
   package hooksconf

   import "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"

   func New(deps *deps.Deps, content string) (*HooksConf, error) {
       root, err := deps.Serializables.ParseYaml(content)
       if err != nil {
           return nil, err
       }
       self := NewEmpty(deps)
       self.Before, err = stringList(root, "before")
       if err != nil {
           return nil, err
       }
       self.After, err = stringList(root, "after")
       if err != nil {
           return nil, err
       }
       return self, nil
   }
   ```
4. Write `bind_methods.go` — one assignment per function field, each a closure over `self`:
   ```go
   package hooksconf

   import "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"

   func bindMethods(deps *deps.Deps, self *HooksConf) {
       self.AddBefore = func(command string) { self.Before = append(self.Before, command) }
       self.AddAfter = func(command string) { self.After = append(self.After, command) }
       self.Render = func() string { return render(deps, self) }
   }
   ```
5. Write `render.go` — build the serializable tree and hand it to `SerializeToYaml`:
   ```go
   package hooksconf

   import "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"

   func render(deps *deps.Deps, self *HooksConf) string {
       root := deps.Serializables.CreateObject()
       root.AddItemToObject("before", self.Before)
       root.AddItemToObject("after", self.After)
       return deps.Serializables.SerializeToYaml(root)
   }
   ```
6. Add a loader/saver pair to `sandbox/internal/utils/` when more than one action reads the file, the way `LoadCommandConf` / `SaveCommandConf` and `LoadModuleConf` do, so every action resolves the same project-relative path.
7. If `agnos start` should write the file, add its template to the `start` asset group following [HandleAssetGroups.md](/docs/Tutorials/HandleAssetGroups.md), and register the file in [Structure.md](/docs/References/Structure.md#agnosconfig).
8. Compile and verify:
   ```bash
   go build ./cmd/... ./sandbox/... ./adapters/... && go run ./cmd/main verify
   ```
