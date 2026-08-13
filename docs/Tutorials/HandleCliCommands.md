# Handle CLI Commands

## Description
Covers adding a command or a flag to the command-line interface — the dispatch triggered by `api.Lib.Sandboxmain`, which lives in [sandbox/cli/](/sandbox/cli/). Add the *library* function the command calls first, following [HandleLibElements.md](/docs/Tutorials/HandleLibElements.md); a command that needs a new OS-bound effect needs [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md) before either.

### Rules
- The interface is inside the closed sandbox: it may not import `adapters/`, `cmd/`, a third-party module, or an OS-bound standard-library package. It reads the command line through `l.Deps.VerbLib`, takes the messages it prints from [sandbox/config/](/sandbox/config/), and prints through `l.Deps.Printf`, and nothing else. See [SandboxIsolation.md](/docs/References/SandboxIsolation.md).
- No display text is hardcoded in the command files. Every line a command prints is a constant defined in [sandbox/config/cli.go](/sandbox/config/cli.go).
- A command does no work of its own: it parses its operands, calls the library functions on `api.Lib`, and reports. Behavior worth having belongs on `api.Lib`, where a Go caller can reach it too.
- Every command returns one of `api.ExitOk`, `api.ExitUsage`, or `api.ExitFailure` — a wrong command line is `ExitUsage`, a well-formed command that could not be carried out is `ExitFailure`.
- Adding a command requires updating the usage screen in `config.Usages` at [sandbox/config/cli.go](/sandbox/config/cli.go) and [Commands.md](/docs/References/Commands.md) in the same commit.

---

## Add CLI Command

### Workflow
1. Add the command to the `Usages` constant in [sandbox/config/cli.go](/sandbox/config/cli.go), in the same column layout as the commands already there:
   ```text
   agnos-cli — a financial tracker on the command line
   …
     largest [category]                            print the largest transaction
   ```
2. Create the command file in [sandbox/cli/commands/](/sandbox/cli/commands/), draining its operands from the injected parser and printing the configured messages through the injected writer:
   ```go
   package commands
   
   import (
       "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
       "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/config"
   )
   
   // Largest runs the `largest` command, printing the biggest single
   // transaction of the budget or of one category.
   func Largest(l *api.Lib) int {
       verb := l.Deps.VerbLib

       listed := l.ListTransactions()
       if name, err := verb.GetNextStringArg(); err == nil {
           stored, found := l.GetCategory(name)
           if !found {
               return Failure(l, config.CategoryNotFound, name)
           }
           listed = stored.ListTransactions()
       }

       if len(listed) == 0 {
           l.Deps.Printf("%s\n", config.NoTransactions)
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
3. Dispatch to it from `Run` in [sandbox/cli/run.go](/sandbox/cli/run.go), in the `switch` over the command word:
   ```go
   case "largest":
       return commands.Largest(l)
   ```
4. Read any flag the command adds **before** the positional arguments are drained, in `Run` — Verb marks a matched flag used, so reading flags first is what leaves only the command words behind:
   ```go
   quiet := verb.IsPresent(config.QuietFlags)
   ```
   Any line the command prints that is not already a message constant needs one, added to [sandbox/config/cli.go](/sandbox/config/cli.go).
5. Build and try it:
   ```bash
   go build ./... && go run ./cmd/main largest groceries
   ```
6. Add the command to the Commands table of [Commands.md](/docs/References/Commands.md), and any flag to its Flags table.
7. Demonstrate it in a script under [examples/cliExamples/](/examples/cliExamples/) when it is worth showing, following [HandleCliExamples.md](/docs/Tutorials/HandleCliExamples.md).

---

## Remove CLI Command

### Workflow
1. Remove the command file from [sandbox/cli/commands/](/sandbox/cli/commands/).
2. Remove the dispatch case for the command from `Run` in [sandbox/cli/run.go](/sandbox/cli/run.go).
3. Remove the command from the usage screen in [sandbox/config/cli.go](/sandbox/config/cli.go).
4. Remove any message constants exclusively used by this command from [sandbox/config/cli.go](/sandbox/config/cli.go).
5. Remove the command from the Commands table of [Commands.md](/docs/References/Commands.md).
6. Update any CLI examples in [examples/cliExamples/](/examples/cliExamples/) that were demonstrating the command.
