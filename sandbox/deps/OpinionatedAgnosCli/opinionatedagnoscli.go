package opinionatedagnoscli

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/argvdeps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/stddeps"
)

// This package is the contract of an *opinionated* lib: unlike every other dep,
// which restates a library's raw capability and nothing more, it carries the
// agnos cli mechanic itself — the declaration a command.yaml becomes, the
// chain that runs a command line against those declarations, how a value is
// bound and how a failure is raised. The project's sandbox/api aliases its
// types, so a handler still reads them as api.Command, api.CommandResponse.
//
// What stays in the project is everything it declares or edits: each
// command's command.yaml and handler.go, the generated
// declarations built from them, and the Handle* files of
// sandbox/internal/cli/errors/ that answer a failure.

// TriggerType is how a Trigger compares the text it is handed.
type TriggerType int

const (
	// TriggerEqual matches a text that is exactly the trigger's Value.
	TriggerEqual TriggerType = iota
	// TriggerPrefix matches a text that is the Value or continues it with a
	// new segment: "/admin" matches "/admin" and "/admin/users", never
	// "/administrator". A Value of "/" matches every path.
	TriggerPrefix
	// TriggerTextPrefix matches a text that begins with the Value, whatever
	// follows it: "/admin" matches "/administrator" too.
	TriggerTextPrefix
	// TriggerSuffix matches a text that ends with the Value.
	TriggerSuffix
	// TriggerRegex matches a text the Value, a regular expression, matches.
	TriggerRegex
	// TriggerOneOf matches a text that is exactly one of the trigger's
	// Values — how a command answers to more than one name.
	TriggerOneOf
)

// Trigger is the condition one declared slice or value has to meet for its
// unit to run at all — the parsed form of one `trigger:` of a route.yaml or a
// command.yaml. Failing it is a non-match, never a usage error or a 400: the
// input is for some other unit.
type Trigger struct {
	// Set tells a declared trigger from none at all; an entry with none
	// matches whatever the request brought.
	Set bool
	// Type is how Value is compared.
	Type TriggerType
	// Value is what the text is compared against.
	Value string
	// Values are the texts a TriggerOneOf accepts; nil on every other type.
	Values []string
	// Negate inverts the comparison: the trigger holds when the text does
	// not match.
	Negate bool
	// IgnoreCase compares without regard to case.
	IgnoreCase bool
}

// ArgType is what the segments an arg reads have to convert to. A segment
// that will not is a non-match: the command line is for some other command.
type ArgType int

const (
	// ArgString takes any segment, bound as a string — or any slice of them,
	// bound as a []string, when the arg reads more than one.
	ArgString ArgType = iota
	// ArgInteger takes one segment reading as a whole number, bound as an
	// int.
	ArgInteger
	// ArgNumber takes one segment reading as a number, bound as a float64.
	ArgNumber
	// ArgUuid takes one segment reading as a canonical uuid, bound as a
	// string.
	ArgUuid
)

// CommandArg is one entry of `args` in command.yaml: the segments of the
// command line from Start to End, both inclusive, End -1 standing for the last
// one. A command line's segments are its leading words — every token before the
// first one starting with "-" — plus every token after a bare "--". The slice
// is bound to the Input field tagged with its Id: one segment in its Type,
// several as a []string.
type CommandArg struct {
	// Id is the Input field the slice is bound to.
	Id string
	// Start is the index of the first segment of the slice.
	Start int
	// End is the index of the last segment of the slice, -1 for the last
	// segment of the command line.
	End int
	// Type is what a one-segment slice converts to.
	Type ArgType
	// Required reports that a command line matching the command without
	// this slice is a usage error.
	Required bool
	// Default is the value bound when the slice is absent, spelled as
	// command.yaml writes it; HasDefault tells an empty default from none.
	Default    string
	HasDefault bool
	// Description is the one-line help text.
	Description string
	// Trigger is what the slice, its segments joined by a space, has to
	// match for the command to run.
	Trigger Trigger
}

// FlagType is the type a CommandFlag is converted to before it reaches
// Input.
type FlagType int

const (
	// FlagString is bound as a string.
	FlagString FlagType = iota
	// FlagInteger is bound as an int.
	FlagInteger
	// FlagNumber is bound as a float64.
	FlagNumber
	// FlagBoolean is bound as a bool: true when one of its keys is on the
	// command line. It takes no value.
	FlagBoolean
	// FlagStringArray is bound as a []string, one element per occurrence.
	FlagStringArray
	// FlagIntegerArray is bound as a []int, one element per occurrence.
	FlagIntegerArray
)

// CommandFlag is one entry of `flags` in command.yaml: the value that follows
// one of Keys on the command line, converted to Type and bound to the Input
// field tagged with its Id.
type CommandFlag struct {
	// Id is the Input field the value is bound to.
	Id string
	// Keys are the spellings a user types it under ("--command", "-c").
	Keys []string
	// Type is what the value converts to.
	Type FlagType
	// Required reports that a command line without it is a usage error.
	Required bool
	// Default is the value bound when the command line brings none,
	// spelled as command.yaml writes it; HasDefault tells an empty default
	// from none.
	Default    string
	HasDefault bool
	// Min and Max bound a numeric value; HasMin and HasMax tell a bound of
	// zero from no bound.
	Min    float64
	HasMin bool
	Max    float64
	HasMax bool
	// Enum is every value accepted, nil for any.
	Enum []string
	// Pattern is a regular expression every value has to match, "" for any.
	Pattern string
	// Trigger is what the value has to match for the command to run.
	Trigger Trigger
	// Description is the one-line help text.
	Description string
}

// CommandResponse is how a handler answers: what it prints and the exit
// status the process ends with. The dispatch hands one per command line to
// every command of the chain.
//
// Answering is what ends a chain: SetStatus answers, and so does Printf, with
// ExitOk when no status was set before it. Error and Log print without
// answering, which is how a middleware says something and still hands the
// command line on.
type CommandResponse struct {
	// SetStatus fixes the exit status; the first status set is the one the
	// process ends with.
	SetStatus func(code int)
	// Printf writes to stdout, and answers ExitOk unless a status was set.
	Printf func(format string, a ...any) (int, error)
	// Eprintf writes to stderr.
	Eprintf func(format string, a ...any) (int, error)
	// Logf writes a progress notice to stderr, silenced by a quiet run.
	Logf func(format string, a ...any) (int, error)
}

// CommandFailureKind is which of the project's own Handle* files of
// sandbox/internal/cli/errors/ a CommandFailure is answered by.
type CommandFailureKind int

const (
	// FailureHandler is a failure a handler returned through Fail, or an
	// error it returned without answering.
	FailureHandler CommandFailureKind = iota
	// FailureNotFound is a command line no command answered.
	FailureNotFound
	// FailureBadUsage is a value that would not bind: a required arg or flag
	// missing, a value that will not convert or is out of its bounds.
	FailureBadUsage
	// FailureUnknownFlag is a token looking like a flag that no command of
	// the chain consumed.
	FailureUnknownFlag
	// FailureUnexpectedArg is any other token no command of the chain
	// consumed.
	FailureUnexpectedArg
)

// CommandFailure is one way a command line did not get answered: which Handle*
// file answers it, the exit status, and why. It is an error too: a
// Handle refuses a command line by returning one, built by Fail.
type CommandFailure struct {
	// Kind is the Handle* file that answers it.
	Kind CommandFailureKind
	// Status is the exit status it carries: ExitFailure or ExitUsage.
	Status int
	// Field is the arg or flag that failed, "" when the failure names none.
	Field string
	// Message is the one-line reason, written for whoever typed the command.
	Message string
	// Cause is what went wrong underneath — an error's text, or the value a
	// handler panicked with — "" when there is nothing below the message.
	Cause string
}

// Error is the failure's Message, which is what makes a CommandFailure an
// error a handler can return.
func (failure *CommandFailure) Error() string {
	return failure.Message
}

// Command is one command of the project, as the sandbox offers it: the whole
// of what its command.yaml declares, plus the handler behind it. Cli.Commands
// holds one per directory under sandbox/internal/commands holding a
// command.yaml, at any depth, each built by that package's generated
// NewCommand, in run order — lowest Priority first.
//
// What Cli.Commands holds is the declaration alone: nothing is ever bound onto
// it. The dispatch copies it with BindCommand, puts the command line and the
// response on the copy, and hands the copy to Matches and Run.
type Command struct {
	// Name is the package directory of the command, snake_case.
	Name string
	// Identifiers are the literal words its first arg answers to ("add-flag"),
	// read off that arg's trigger by the build; the first is the one the help
	// screen prints. A command whose first arg is no literal has none.
	Identifiers []string
	// Priority is the rung this command runs on when several match one
	// command line: the dispatch runs them from the lowest upwards and stops
	// at the first one that answers.
	Priority int
	// Segments is how many segments the command line has to have for the
	// command to run, 0 for any count.
	Segments int
	// Strict reports that every token of the command line has to be consumed
	// by the command or a command run before it: a middleware is not strict.
	Strict bool
	// Pattern is its command line as it reads in help and docs.
	Pattern string
	// Category groups it on the general help screen.
	Category string
	// Summary is the one-line description.
	Summary string
	// Description is the paragraph the per-command help screen prints.
	Description string
	// Examples are whole command lines the help screen prints.
	Examples []string
	// Hidden keeps it off the general help screen without disabling it.
	Hidden bool
	// Args are the slices of the segments it reads, in declaration order.
	Args []CommandArg
	// Flags are the flags it reads, in declaration order.
	Flags []CommandFlag

	// Handle is the command package's own Handle,
	// closed over the sandbox: a func(props *commandprops.CommandProps,
	// input *Input, response *CommandResponse) error whose Input is that
	// package's generated struct. It is held as any because every command's
	// Input is a type of its own; the dispatch builds and fills one by
	// reflection and calls it.
	Handle any

	// Argv is the command line this copy was bound from.
	Argv []string
	// Consumed tracks, index for index against Argv, the tokens a command of
	// this chain read. It is one slice shared by every copy of one command
	// line, so what a middleware consumed, the strict command after it does
	// not have to.
	Consumed []bool
	// Props is one command line's *commandprops.CommandProps, shared by every
	// command of the chain: what a middleware sets on it, the commands after
	// it read. It is held as any because the project types it under
	// sandbox/internal, which no contract may name; a Handle* file reads it
	// back with command.Props.(*commandprops.CommandProps).
	Props any
	// Response is the one response of the command line.
	Response *CommandResponse

	// Failure is why this command is being handed to one of the project's
	// Handle* files, nil on a normal run. The dispatch sets it as it raises.
	Failure *CommandFailure

	// Matches reports whether one bound copy — its Argv set — is for
	// this command: the segment count, and every arg and every flag declaring
	// a trigger, match it. NewCommand sets it to the lib's own matcher; set it
	// to decide by hand. The dispatch reads a nil one with the lib's matcher
	// too.
	Matches func(bound *Command) bool
	// Run binds one bound copy's command line onto a fresh Input
	// and runs Handle with it. It returns the failure the handler
	// did not answer itself, nil otherwise; a handler answering nothing hands
	// the command line to the next command of the chain. Nil binds with the
	// lib's own binder; set it to bind by hand.
	Run func(bound *Command) error
}

// Cli is the CLI surface of the sandbox: every command the project declares,
// and the dispatch that reads a command line against them. It is built by
// sandbox/internal/generated/cli/new.go, generated by the build.
type Cli struct {
	// Main is the dispatch-and-parse entry point: it hands the command
	// line to the lib's Main against these Commands, and returns the exit
	// status the command line was answered with.
	Main func(args []string) int
	// Commands is every command the project declares, in run order —
	// lowest `priority` first, then by name — each built by the generated
	// NewCommand of its own package. The dispatch reads the command line
	// against these declarations; a caller holding the sandbox reads the
	// same surface without one.
	Commands []*Command
	// Fail answers one failure with the project's own handler for it: it
	// reads command.Failure and calls the matching Handle* of
	// sandbox/internal/cli/errors/. The dispatch raises through it; a handler
	// returns a failure built by Fail rather than calling it directly.
	Fail func(command *Command) error
}

const (
	// ExitOk reports that the command did what it was asked to do.
	ExitOk = 0
	// ExitFailure reports that a well-formed command could not be carried
	// out.
	ExitFailure = 1
	// ExitUsage reports that the command line itself was wrong — an unknown
	// command or flag, a missing operand, an unparsable amount. Every such
	// error exits with this one code, whichever command produced it.
	ExitUsage = 2
)

// MainProps is everything one run of Main reads: the surface, the command
// line and the two deps the dispatch reaches the outside world through. The
// generated registry builds it per call, so nothing the lib holds goes stale.
type MainProps struct {
	// Cli is the surface the command line is read against. It is a pointer,
	// read as the line runs, so a caller that replaced a field of
	// sandbox.Cli is followed.
	Cli *Cli
	// Args is the command line, without the program name.
	Args []string
	// NewProps builds the one *commandprops.CommandProps every command of
	// the chain is handed — a type of the project, which is why the lib is
	// handed a constructor rather than naming it.
	NewProps func() any
	// StdDeps is the channel every print of the response goes through, read at
	// the moment of each print: a middleware that silences StdDeps.Logf silences
	// the response's Log too.
	StdDeps *stddeps.Contract
	// ArgvDeps is the parser the flags are read with.
	ArgvDeps argvdeps.Contract
}

// Contract is the cli lib injected whole as the Deps.OpinionatedAgnosCli field.
type Contract struct {
	// Main is the whole dispatch layer: it runs every command of
	// props.Cli.Commands the command line is for, in order, and returns the
	// exit status the line was answered with. Each command runs in turn until
	// one of them answers — sets a status, or prints to stdout, which answers
	// ExitOk; a middleware that does neither has declined, and a strict
	// command that does neither ran silently and answers ExitOk. Every way a
	// line can end without an answer — nothing matched, a value that will not
	// bind, a token nobody read, a panic — is raised through props.Cli.Fail
	// and answered by one of the project's own Handle* files.
	Main func(props MainProps) int

	// NewCommand returns an empty Command with every slice open and
	// Matches set to the lib's own matcher, reading flags through
	// argvdeps — the base a command's generated NewCommand fills its
	// declaration on.
	NewCommand func(argvdeps argvdeps.Contract) *Command

	// BindCommand copies one declaration into the command a single command
	// line runs on: the same declared fields — the slices are read-only and
	// shared — with no command line, response or failure yet.
	BindCommand func(command *Command) *Command

	// Fail is how a Handle refuses a command line: it builds
	// the failure and returns it as an error, and the handler returns it in
	// turn, which the dispatch raises on the command:
	//
	//	return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", err.Error())
	//
	// A message left empty is filled by handle_failure.go's own wording.
	Fail func(status int, field string, message string) error

	// FailWithCause is Fail carrying what went wrong underneath — an error's
	// text, meant for the log, never for whoever typed the command.
	FailWithCause func(status int, field string, message string, cause string) error

	// FailureOf is the failure one of the project's Handle* files is
	// answering. A file stands for one kind and one wording, and passes its
	// status and wording as the fallback: a failure that carries its own
	// answers with that, and one that carries none — "no command" is raised
	// with no message at all — answers with the file's.
	FailureOf func(command *Command, status int, message string) CommandFailure

	// MatchTrigger reports whether one text meets a trigger: the comparison
	// its Type names, without regard to case when it declares IgnoreCase,
	// inverted when it declares Negate. Segmented is true for a text a prefix
	// reads segment by segment — a route's path slice, a command's segments.
	// Every layer matching on a declaration shares it, so a trigger reads the
	// same in a route.yaml and in a command.yaml.
	MatchTrigger func(trigger Trigger, text string, segmented bool) bool
}
