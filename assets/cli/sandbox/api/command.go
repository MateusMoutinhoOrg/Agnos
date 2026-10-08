package api

import (
	opinatedagnoscli "{{.Module}}/sandbox/deps/OpinatedAgnosCli"
)

// Every type here is the OpinatedAgnosCli contract's own, aliased so a
// handler reads it as api.<Name>: sandbox/deps/OpinatedAgnosCli holds the
// whole doc of each.

// ArgType is what the segments an arg reads have to convert to.
type ArgType = opinatedagnoscli.ArgType

// StringArg takes any segment, bound as a string, or several as a []string.
const StringArg = opinatedagnoscli.StringArg

// IntegerArg takes one segment reading as a whole number, bound as an int.
const IntegerArg = opinatedagnoscli.IntegerArg

// NumberArg takes one segment reading as a number, bound as a float64.
const NumberArg = opinatedagnoscli.NumberArg

// UuidArg takes one segment reading as a canonical uuid, bound as a string.
const UuidArg = opinatedagnoscli.UuidArg

// CommandArg is one entry of `args` in command.yaml.
type CommandArg = opinatedagnoscli.CommandArg

// FlagType is the type a CommandFlag is converted to before it reaches
// Entries.
type FlagType = opinatedagnoscli.FlagType

// StringFlag is bound as a string.
const StringFlag = opinatedagnoscli.StringFlag

// IntegerFlag is bound as an int.
const IntegerFlag = opinatedagnoscli.IntegerFlag

// NumberFlag is bound as a float64.
const NumberFlag = opinatedagnoscli.NumberFlag

// BooleanFlag is bound as a bool, true when one of its keys is present.
const BooleanFlag = opinatedagnoscli.BooleanFlag

// StringArrayFlag is bound as a []string, one element per occurrence.
const StringArrayFlag = opinatedagnoscli.StringArrayFlag

// IntegerArrayFlag is bound as a []int, one element per occurrence.
const IntegerArrayFlag = opinatedagnoscli.IntegerArrayFlag

// CommandFlag is one entry of `flags` in command.yaml.
type CommandFlag = opinatedagnoscli.CommandFlag

// CommandResponse is how a handler answers: what it prints and the exit
// status the process ends with.
type CommandResponse = opinatedagnoscli.CommandResponse

// CommandFailureKind is which Handle* file of sandbox/internal/cli/errors/ a
// CommandFailure is answered by.
type CommandFailureKind = opinatedagnoscli.CommandFailureKind

// HandlerFailure is a failure a handler returned through
// Deps.OpinatedAgnosCli.Fail, or an error it returned without answering.
const HandlerFailure = opinatedagnoscli.HandlerFailure

// NotFoundFailure is a command line no command answered.
const NotFoundFailure = opinatedagnoscli.NotFoundFailure

// BadUsageFailure is a value that would not bind.
const BadUsageFailure = opinatedagnoscli.BadUsageFailure

// UnknownFlagFailure is a token looking like a flag no command consumed.
const UnknownFlagFailure = opinatedagnoscli.UnknownFlagFailure

// UnexpectedArgFailure is any other token no command consumed.
const UnexpectedArgFailure = opinatedagnoscli.UnexpectedArgFailure

// CommandFailure is one way a command line did not get answered; an
// InternalPureHandler refuses a command line by returning one, built by
// Deps.OpinatedAgnosCli.Fail.
type CommandFailure = opinatedagnoscli.CommandFailure

// Command is one command of the project: the whole of what its command.yaml
// declares, plus the handler behind it.
type Command = opinatedagnoscli.Command
