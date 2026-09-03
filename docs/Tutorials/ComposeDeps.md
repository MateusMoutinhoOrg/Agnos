# Compose a Custom Deps

## Description
Covers building the `Deps` a scaffolded project's sandbox runs on from parts other than `standard.New()` — replacing one field for a test, or assembling a different mix of adapter libs for another deployment. Writing a new capability is [AddAdapterLib.md](/docs/Tutorials/AddAdapterLib.md); the trade-offs of struct contracts, which make this a plain assignment, are [StructContracts.md](/docs/References/StructContracts.md).

### Rules
- Patch the `deps.Deps` value **before** calling `sandbox.New(&deps)`. The binders capture the pointer they were handed, so a later reassignment is seen; a struct copy made after `New` is not what the sandbox reads.
- An unfilled field compiles and panics on first call. Start from `standard.New()` and patch, rather than from an empty struct, unless you are sure which fields the code you run reaches.
- `adapters/availables/standard/new.go` is regenerated on every build — a hand-written assembly goes in its own `adapters/availables/<name>/new.go`, which `verify` allows and `build` leaves alone.

---

## Workflow
1. Start from the standard assembly and replace the one behavior you want to control. A fixed clock and a captured output are the everyday testing pair:
   ```go
   package main

   import (
       "bytes"
       "fmt"
       "time"

       "github.com/you/my-tool/adapters/availables/standard"
       "github.com/you/my-tool/sandbox"
   )

   func main() {
       deps := standard.New()

       var out bytes.Buffer
       deps.Std.Printf = func(format string, a ...any) (int, error) {
           return fmt.Fprintf(&out, format, a...)
       }
       deps.Clockdeps.Now = func() time.Time { return time.Unix(0, 0) }

       lib := sandbox.New(&deps)
       lib.Cli.CliMain([]string{"greet", "-n", "bob"})
       fmt.Print(out.String()) // the command's output, captured
   }
   ```
2. Assemble a different mix as an adapter of your own when the replacement is permanent. Bind the libs you want and nothing else; the package is yours, so the name of its constructor is too:
   ```go
   // adapters/availables/testing/new.go
   package testing

   import (
       "github.com/you/my-tool/adapters/libs/iodeps"
       "github.com/you/my-tool/adapters/libs/std"
       "github.com/you/my-tool/adapters/libs/verb"
       "github.com/you/my-tool/sandbox/deps"
   )

   // New assembles a Deps for tests: the real libs the sandbox cannot do
   // without, and every other field left for the test to fill.
   func New() deps.Deps {
       deps := deps.Deps{}
       iodeps.Bind(&deps)
       std.Bind(&deps)
       verb.Bind(&deps)
       return deps
   }
   ```
3. Wire it where an adapter and the sandbox meet — a `cmd/<name>/main.go` of its own or a test file. `cmd/main/main.go` belongs to the `cli` asset group and is regenerated, so leave it on `standard`:
   ```go
   deps := testing.New()
   lib := sandbox.New(&deps)
   ```
4. Compare your assembly against `sandbox/deps/deps.go` field by field before relying on it: the compiler does not report a field you left nil.
