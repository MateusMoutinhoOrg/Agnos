package api

// ArgType is what the segments an arg reads have to convert to. A segment
// that will not is a non-match: the command line is for some other command.
type ArgType int

const (
	// StringArg takes any segment, bound as a string — or any slice of them,
	// bound as a []string, when the arg reads more than one.
	StringArg ArgType = iota
	// IntegerArg takes one segment reading as a whole number, bound as an
	// int.
	IntegerArg
	// NumberArg takes one segment reading as a number, bound as a float64.
	NumberArg
	// UuidArg takes one segment reading as a canonical uuid, bound as a
	// string.
	UuidArg
)

// CommandArg is one entry of `args` in command.yaml: the segments of the
// command line from Start to End, both inclusive, End -1 standing for the last
// one. A command line's segments are its leading words — every token before the
// first one starting with "-" — plus every token after a bare "--". The slice
// is bound to the Entries field tagged with its Id: one segment in its Type,
// several as a []string.
type CommandArg struct {
	// Id is the Entries field the slice is bound to.
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
// Entries.
type FlagType int

const (
	// StringFlag is bound as a string.
	StringFlag FlagType = iota
	// IntegerFlag is bound as an int.
	IntegerFlag
	// NumberFlag is bound as a float64.
	NumberFlag
	// BooleanFlag is bound as a bool: true when one of its keys is on the
	// command line. It takes no value.
	BooleanFlag
	// StringArrayFlag is bound as a []string, one element per occurrence.
	StringArrayFlag
	// IntegerArrayFlag is bound as a []int, one element per occurrence.
	IntegerArrayFlag
)

// CommandFlag is one entry of `flags` in command.yaml: the value that follows
// one of Keys on the command line, converted to Type and bound to the Entries
// field tagged with its Id.
type CommandFlag struct {
	// Id is the Entries field the value is bound to.
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
	// Error writes to stderr.
	Error func(format string, a ...any) (int, error)
	// Log writes a progress notice to stderr, silenced by a quiet run.
	Log func(format string, a ...any) (int, error)
}

// CommandFailureKind is which of the project's own Handle* files of
// sandbox/internal/cli/errors/ a CommandFailure is answered by.
type CommandFailureKind int

const (
	// HandlerFailure is a failure a handler returned through cliio.Fail, or
	// an error it returned without answering.
	HandlerFailure CommandFailureKind = iota
	// NotFoundFailure is a command line no command answered.
	NotFoundFailure
	// BadUsageFailure is a value that would not bind: a required arg or flag
	// missing, a value that will not convert or is out of its bounds.
	BadUsageFailure
	// UnknownFlagFailure is a token looking like a flag that no command of
	// the chain consumed.
	UnknownFlagFailure
	// UnexpectedArgFailure is any other token no command of the chain
	// consumed.
	UnexpectedArgFailure
)

// CommandFailure is one way a command line did not get answered: which Handle*
// file answers it, the exit status, and why. It is an error too: an
// InternalPureHandler refuses a command line by returning one, built by
// cliio.Fail.
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
// response on the copy, and hands the copy to IsActionable and CommandHandler.
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
	// Help is the one-line description.
	Help string
	// LongDescription is the paragraph the per-command help screen prints.
	LongDescription string
	// Examples are whole command lines the help screen prints.
	Examples []string
	// Hidden keeps it off the general help screen without disabling it.
	Hidden bool
	// Args are the slices of the segments it reads, in declaration order.
	Args []CommandArg
	// Flags are the flags it reads, in declaration order.
	Flags []CommandFlag

	// InternalPureHandler is the command package's own InternalPureHandler,
	// closed over the sandbox: a func(props *commandprops.CommandProps,
	// entries *Entries, response *CommandResponse) error whose Entries is that package's
	// generated struct. It is held as any because every command's Entries is
	// a type of its own; CommandHandler builds and fills one through
	// Deps.Reflectdeps and calls it.
	InternalPureHandler any

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
	// sandbox/internal, which sandbox/api may not name; a Handle* file reads
	// it back with command.Props.(*commandprops.CommandProps).
	Props any
	// Response is the one response of the command line.
	Response *CommandResponse

	// Failure is why this command is being handed to one of the project's
	// Handle* files, nil on a normal run. It is set by cliio.Raise.
	Failure *CommandFailure

	// IsActionable reports whether one bound copy — its Argv set — is for
	// this command: the segment count, and every arg and every flag declaring
	// a trigger, match it.
	IsActionable func(bound *Command) bool
	// CommandHandler binds one bound copy's command line onto a fresh Entries
	// and runs InternalPureHandler with it. It returns the failure the handler
	// did not answer itself, nil otherwise; a handler answering nothing hands
	// the command line to the next command of the chain.
	CommandHandler func(bound *Command) error
}

// NewCommand returns an empty Command with every slice open. The generic
// sandbox/internal/generated/cli/command.NewCommand fills the matcher and the
// handler on top of it, and a command's generated NewCommand its declaration.
func NewCommand() *Command {
	return &Command{
		Identifiers: []string{},
		Strict:      true,
		Examples:    []string{},
		Args:        []CommandArg{},
		Flags:       []CommandFlag{},
	}
}

// BindCommand copies one declaration into the command a single command line
// runs on: the same declared fields — the slices are read-only and shared —
// with no command line, response or failure yet. The dispatch calls it once
// per command of the chain, so what Cli.Commands holds is never written to.
func BindCommand(command *Command) *Command {
	bound := *command
	bound.Argv = nil
	bound.Consumed = nil
	bound.Props = nil
	bound.Response = nil
	bound.Failure = nil
	return &bound
}
