# Handle CLI Commands

## Description
Covers adding a command or a flag to the command-line interface — the dispatch behind `api.Lib.Sandboxmain`, which lives in [sandbox/internal/cli/](/sandbox/internal/cli/). Add the *library* function the command calls first, following [HandleLibElements.md](/docs/Development/Tutorials/HandleLibElements.md); a command that needs a new OS-bound effect needs [HandleDependencies.md](/docs/Development/Tutorials/HandleDependencies.md) before either.

### Rules
- The interface is inside the closed sandbox: it may not import `adapters/`, `cmd/`, a third-party module, or an OS-bound standard-library package. It reads the command line through `l.Deps.VerbLib` and prints through `l.Deps.Printf`, and nothing else. See [SandboxIsolation.md](/docs/Development/References/SandboxIsolation.md).
- A command does no work of its own: it parses its operands, calls the library functions on `api.Lib`, and reports. Behavior worth having belongs on `api.Lib`, where a Go caller can reach it too.
- Every command returns one of `api.ExitOk`, `api.ExitUsage`, or `api.ExitFailure` — a wrong command line is `ExitUsage`, a well-formed command that could not be carried out is `ExitFailure`.
- Adding a command requires updating the `Usage` screen in `sandbox/internal/cli/cli.go` and [Commands.md](/docs/CliUsage/References/Commands.md) in the same commit.

---

## Add a CLI Command

### Workflow
1. Add the command to the `Usage` constant in `sandbox/internal/cli/cli.go`, in the same column layout as the commands already there:
   ```go
   const Usage = `agnos — a financial tracker on the command line
   …
     largest [category]                            print the largest transaction
   `
   ```
2. Write the handler beside the other command helpers, draining its operands from the injected parser and printing through the injected writer:
   ```go
   // largest runs the `largest` command, printing the biggest single
   // transaction of the budget or of one category.
   func largest(l *api.Lib) int {
       verb := l.Deps.VerbLib

       listed := l.ListTransactions()
       if name, err := verb.GetNextStringArg(); err == nil {
           stored, found := l.GetCategory(name)
           if !found {
               return failure(l, "no category named "+quote(name))
           }
           listed = stored.ListTransactions()
       }

       if len(listed) == 0 {
           l.Deps.Printf("no transactions yet\n")
           return api.ExitOk
       }
       biggest := listed[0]
       for _, written := range listed {
           if written.Amount > biggest.Amount {
               biggest = written
           }
       }
       l.Deps.Printf("%s\n", biggest.String())
       return api.ExitOk
   }
   ```
3. Dispatch to it from `Run`, in the `switch` over the command word:
   ```go
   case "largest":
       return largest(l)
   ```
4. Read any flag the command adds **before** the positional arguments are drained, in `Run` — Verb marks a matched flag used, so reading flags first is what leaves only the command words behind:
   ```go
   quiet := verb.IsPresent(QuietFlags)
   ```
5. Build and try it:
   ```bash
   go build ./... && go run ./cmd/main largest groceries
   ```
6. Add the command to the Commands table of [Commands.md](/docs/CliUsage/References/Commands.md), and any flag to its Flags table.
7. Demonstrate it in a script under [examples/cliExamples/](/examples/cliExamples/) when it is worth showing, following [HandleCliExamples.md](/docs/Development/Tutorials/HandleCliExamples.md).
