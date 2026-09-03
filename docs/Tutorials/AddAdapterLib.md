# Add a Capability to a Project

## Description
Covers giving a scaffolded project an effect none of the installable deps provides — a clock you control, a database driver, a message queue — by writing the pair every dep is made of: a contract under `sandbox/deps/<x>/` and an adapter lib under `adapters/libs/<x>/`. Installing a shipped capability instead is [ManageDeps.md](/docs/Tutorials/ManageDeps.md); why the pair is split across the sandbox wall is [SandboxIsolation.md](/docs/References/SandboxIsolation.md); why contracts are structs and not interfaces is [StructContracts.md](/docs/References/StructContracts.md).

### Rules
- The contract package declares one `Lib` struct of function fields and imports only the standard library and other `sandbox/deps` packages. `verify` enforces both.
- The adapter lib is a package exporting exactly one binder, `Bind(deps *deps.Deps)`, that fills its one `Deps` field. `agnos build` generates a call to it under that name, so the name is not a choice.
- The `Deps` field name is mechanical: the title-cased contract directory. `sandbox/deps/clockdeps/` is reached as `deps.Clockdeps`, and sandbox code must use that spelling, because the struct is regenerated from the directory listing on every build.
- The contract and the adapter directories usually share a name; they need not (`argvdeps` is implemented by `verb`). The collectors iterate the two trees independently.
- Both files must follow their specifications — DepsContract and AdapterLib, located in [Specs.md](/docs/References/Specs.md).

---

## Workflow
1. Declare the contract. One struct, function fields, doc comments saying what each field promises; nothing that names an implementation:
   ```go
   // sandbox/deps/clockdeps/clockdeps.go
   package clockdeps

   import "time"

   // Lib is the clock injected whole as the Deps.Clockdeps field.
   type Lib struct {
       // Now returns the current wall-clock time.
       Now func() time.Time
       // Sleep pauses the calling goroutine for the given duration.
       Sleep func(d time.Duration)
   }
   ```
2. Implement it outside the wall. The adapter lib imports `sandbox/deps` for the struct to fill and its own contract for the type, and may import anything else — this is the only place `time.Now` or a driver is allowed:
   ```go
   // adapters/libs/clockdeps/Clock.go
   package clockdeps

   import (
       "time"

       "github.com/you/my-tool/sandbox/deps"
       clockdeps "github.com/you/my-tool/sandbox/deps/clockdeps"
   )

   // Bind fills deps.Deps.Clockdeps with the real wall clock.
   func Bind(deps *deps.Deps) {
       deps.Clockdeps = clockdeps.Lib{
           Now:   time.Now,
           Sleep: time.Sleep,
       }
   }
   ```
3. Rebuild. The collectors list `sandbox/deps/*/` and `adapters/libs/*/`, so `Deps` gains a `Clockdeps clockdeps.Lib` field and `standard.New` a `clockdeps.Bind(&deps)` call without you editing either:
   ```bash
   agnos build
   grep Clockdeps sandbox/deps/deps.go
   grep clockdeps adapters/availables/standard/new.go
   ```
4. Use it from a handler through the mechanical name:
   ```go
   started := deps.Clockdeps.Now()
   ```
5. Give an implementation that is created per call a constructor field instead of a set of operations, the way `requestdeps.Lib.NewRequest` and `argvdeps.Lib.New` do: the contract holds one `New func(...) Thing` and a `Thing` struct of function fields the adapter fills per call.
6. Write a second implementation of the same contract as another adapter lib when a test or a different deployment needs it, and compose it into a `Deps` of your own following [ComposeDeps.md](/docs/Tutorials/ComposeDeps.md). `standard.New` always binds **every** lib under `adapters/libs/`, so two libs filling the same field would have the last one in listing order win — keep the alternative outside `adapters/libs/` or in its own `availables/<name>/`.
