package api

import (
	opinionatedagnoscli "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/OpinionatedAgnosCli"
)

// Every type here is the OpinionatedAgnosCli contract's own, aliased so a
// handler reads it as api.<Name>: sandbox/deps/OpinionatedAgnosCli holds the
// whole doc of each.

// ArgType is what the segments an arg reads have to convert to.
type ArgType = opinionatedagnoscli.ArgType

// ArgString takes any segment, bound as a string, or several as a []string.
const ArgString = opinionatedagnoscli.ArgString

// ArgInteger takes one segment reading as a whole number, bound as an int.
const ArgInteger = opinionatedagnoscli.ArgInteger

// ArgNumber takes one segment reading as a number, bound as a float64.
const ArgNumber = opinionatedagnoscli.ArgNumber

// ArgUuid takes one segment reading as a canonical uuid, bound as a string.
const ArgUuid = opinionatedagnoscli.ArgUuid

// CommandArg is one entry of `args` in command.yaml.
type CommandArg = opinionatedagnoscli.CommandArg

// FlagType is the type a CommandFlag is converted to before it reaches
// Input.
type FlagType = opinionatedagnoscli.FlagType

// FlagString is bound as a string.
const FlagString = opinionatedagnoscli.FlagString

// FlagInteger is bound as an int.
const FlagInteger = opinionatedagnoscli.FlagInteger

// FlagNumber is bound as a float64.
const FlagNumber = opinionatedagnoscli.FlagNumber

// FlagBoolean is bound as a bool, true when one of its keys is present.
const FlagBoolean = opinionatedagnoscli.FlagBoolean

// FlagStringArray is bound as a []string, one element per occurrence.
const FlagStringArray = opinionatedagnoscli.FlagStringArray

// FlagIntegerArray is bound as a []int, one element per occurrence.
const FlagIntegerArray = opinionatedagnoscli.FlagIntegerArray

// CommandFlag is one entry of `flags` in command.yaml.
type CommandFlag = opinionatedagnoscli.CommandFlag

// CommandResponse is how a handler answers: what it prints and the exit
// status the process ends with.
type CommandResponse = opinionatedagnoscli.CommandResponse

// CommandFailureKind is which Handle* file of sandbox/internal/cli/errors/ a
// CommandFailure is answered by.
type CommandFailureKind = opinionatedagnoscli.CommandFailureKind

// FailureHandler is a failure a handler returned through
// Deps.OpinionatedAgnosCli.Fail, or an error it returned without answering.
const FailureHandler = opinionatedagnoscli.FailureHandler

// FailureNotFound is a command line no command answered.
const FailureNotFound = opinionatedagnoscli.FailureNotFound

// FailureBadUsage is a value that would not bind.
const FailureBadUsage = opinionatedagnoscli.FailureBadUsage

// FailureUnknownFlag is a token looking like a flag no command consumed.
const FailureUnknownFlag = opinionatedagnoscli.FailureUnknownFlag

// FailureUnexpectedArg is any other token no command consumed.
const FailureUnexpectedArg = opinionatedagnoscli.FailureUnexpectedArg

// CommandFailure is one way a command line did not get answered; a
// Handle refuses a command line by returning one, built by
// Deps.OpinionatedAgnosCli.Fail.
type CommandFailure = opinionatedagnoscli.CommandFailure

// Command is one command of the project: the whole of what its command.yaml
// declares, plus the handler behind it.
type Command = opinionatedagnoscli.Command
