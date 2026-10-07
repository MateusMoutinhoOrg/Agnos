# `sandbox/api/command.go`

| Constant | Value | Description |
| --- | --- | --- |
| `StringArg` | `opinatedagnoscli.StringArg` | StringArg takes any segment, bound as a string, or several as a []string. |
| `IntegerArg` | `opinatedagnoscli.IntegerArg` | IntegerArg takes one segment reading as a whole number, bound as an int. |
| `NumberArg` | `opinatedagnoscli.NumberArg` | NumberArg takes one segment reading as a number, bound as a float64. |
| `UuidArg` | `opinatedagnoscli.UuidArg` | UuidArg takes one segment reading as a canonical uuid, bound as a string. |
| `StringFlag` | `opinatedagnoscli.StringFlag` | StringFlag is bound as a string. |
| `IntegerFlag` | `opinatedagnoscli.IntegerFlag` | IntegerFlag is bound as an int. |
| `NumberFlag` | `opinatedagnoscli.NumberFlag` | NumberFlag is bound as a float64. |
| `BooleanFlag` | `opinatedagnoscli.BooleanFlag` | BooleanFlag is bound as a bool, true when one of its keys is present. |
| `StringArrayFlag` | `opinatedagnoscli.StringArrayFlag` | StringArrayFlag is bound as a []string, one element per occurrence. |
| `IntegerArrayFlag` | `opinatedagnoscli.IntegerArrayFlag` | IntegerArrayFlag is bound as a []int, one element per occurrence. |
| `HandlerFailure` | `opinatedagnoscli.HandlerFailure` | HandlerFailure is a failure a handler returned through Deps.OpinatedAgnosCli.Fail, or an error it returned without answering. |
| `NotFoundFailure` | `opinatedagnoscli.NotFoundFailure` | NotFoundFailure is a command line no command answered. |
| `BadUsageFailure` | `opinatedagnoscli.BadUsageFailure` | BadUsageFailure is a value that would not bind. |
| `UnknownFlagFailure` | `opinatedagnoscli.UnknownFlagFailure` | UnknownFlagFailure is a token looking like a flag no command consumed. |
| `UnexpectedArgFailure` | `opinatedagnoscli.UnexpectedArgFailure` | UnexpectedArgFailure is any other token no command consumed. |

## `ArgType`

ArgType is what the segments an arg reads have to convert to.

`type ArgType = opinatedagnoscli.ArgType`

## `CommandArg`

CommandArg is one entry of `args` in command.yaml.

`type CommandArg = opinatedagnoscli.CommandArg`

## `FlagType`

FlagType is the type a CommandFlag is converted to before it reaches Entries.

`type FlagType = opinatedagnoscli.FlagType`

## `CommandFlag`

CommandFlag is one entry of `flags` in command.yaml.

`type CommandFlag = opinatedagnoscli.CommandFlag`

## `CommandResponse`

CommandResponse is how a handler answers: what it prints and the exit status the process ends with.

`type CommandResponse = opinatedagnoscli.CommandResponse`

## `CommandFailureKind`

CommandFailureKind is which Handle* file of sandbox/internal/cli/errors/ a CommandFailure is answered by.

`type CommandFailureKind = opinatedagnoscli.CommandFailureKind`

## `CommandFailure`

CommandFailure is one way a command line did not get answered; an InternalPureHandler refuses a command line by returning one, built by Deps.OpinatedAgnosCli.Fail.

`type CommandFailure = opinatedagnoscli.CommandFailure`

## `Command`

Command is one command of the project: the whole of what its command.yaml declares, plus the handler behind it.

`type Command = opinatedagnoscli.Command`

[every contract](doc.md)
