# `sandbox/api/generated.command.go`

| Constant | Value | Description |
| --- | --- | --- |
| `ArgString` | `opinionatedagnoscli.ArgString` | ArgString takes any segment, bound as a string, or several as a []string. |
| `ArgInteger` | `opinionatedagnoscli.ArgInteger` | ArgInteger takes one segment reading as a whole number, bound as an int. |
| `ArgNumber` | `opinionatedagnoscli.ArgNumber` | ArgNumber takes one segment reading as a number, bound as a float64. |
| `ArgUuid` | `opinionatedagnoscli.ArgUuid` | ArgUuid takes one segment reading as a canonical uuid, bound as a string. |
| `FlagString` | `opinionatedagnoscli.FlagString` | FlagString is bound as a string. |
| `FlagInteger` | `opinionatedagnoscli.FlagInteger` | FlagInteger is bound as an int. |
| `FlagNumber` | `opinionatedagnoscli.FlagNumber` | FlagNumber is bound as a float64. |
| `FlagBoolean` | `opinionatedagnoscli.FlagBoolean` | FlagBoolean is bound as a bool, true when one of its keys is present. |
| `FlagStringArray` | `opinionatedagnoscli.FlagStringArray` | FlagStringArray is bound as a []string, one element per occurrence. |
| `FlagIntegerArray` | `opinionatedagnoscli.FlagIntegerArray` | FlagIntegerArray is bound as a []int, one element per occurrence. |
| `FailureHandler` | `opinionatedagnoscli.FailureHandler` | FailureHandler is a failure a handler returned through Deps.OpinionatedAgnosCli.Fail, or an error it returned without answering. |
| `FailureNotFound` | `opinionatedagnoscli.FailureNotFound` | FailureNotFound is a command line no command answered. |
| `FailureBadUsage` | `opinionatedagnoscli.FailureBadUsage` | FailureBadUsage is a value that would not bind. |
| `FailureUnknownFlag` | `opinionatedagnoscli.FailureUnknownFlag` | FailureUnknownFlag is a token looking like a flag no command consumed. |
| `FailureUnexpectedArg` | `opinionatedagnoscli.FailureUnexpectedArg` | FailureUnexpectedArg is any other token no command consumed. |

## `ArgType`

ArgType is what the segments an arg reads have to convert to.

`type ArgType = opinionatedagnoscli.ArgType`

## `CommandArg`

CommandArg is one entry of `args` in command.yaml.

`type CommandArg = opinionatedagnoscli.CommandArg`

## `FlagType`

FlagType is the type a CommandFlag is converted to before it reaches Input.

`type FlagType = opinionatedagnoscli.FlagType`

## `CommandFlag`

CommandFlag is one entry of `flags` in command.yaml.

`type CommandFlag = opinionatedagnoscli.CommandFlag`

## `CommandResponse`

CommandResponse is how a handler answers: what it prints and the exit status the process ends with.

`type CommandResponse = opinionatedagnoscli.CommandResponse`

## `CommandFailureKind`

CommandFailureKind is which Handle* file of sandbox/internal/cli/errors/ a CommandFailure is answered by.

`type CommandFailureKind = opinionatedagnoscli.CommandFailureKind`

## `CommandFailure`

CommandFailure is one way a command line did not get answered; a Handle refuses a command line by returning one, built by Deps.OpinionatedAgnosCli.Fail.

`type CommandFailure = opinionatedagnoscli.CommandFailure`

## `Command`

Command is one command of the project: the whole of what its command.yaml declares, plus the handler behind it.

`type Command = opinionatedagnoscli.Command`

[every contract](doc.md)
