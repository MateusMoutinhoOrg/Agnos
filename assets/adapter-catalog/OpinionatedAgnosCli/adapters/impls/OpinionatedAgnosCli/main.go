package opinionatedagnoscli

import (
	"fmt"
	"strings"

	opinionatedagnoscli "{{.Module}}/sandbox/deps/OpinionatedAgnosCli"
)

// line is one command line being answered: the props it runs with and what
// every command of its chain shares — the CommandProps, the Consumed slice and
// the one response, with the reader of the status it was answered with.
type line struct {
	props    opinionatedagnoscli.MainProps
	args     []string
	consumed []bool
	shared   any
	response *opinionatedagnoscli.CommandResponse
	status   func() (int, bool)
}

// cliMain runs every command of props.Cli.Commands the command line is for,
// in the order the collector put them — lowest `priority` first — and returns
// the exit status the line was answered with. Whether a command is for the
// line is its Matches's to say; binding and running it is its
// Run's.
//
// More than one command may be for one line, which is what a chain is: each one
// runs in turn until one of them **answers** — sets a status, or prints to
// stdout, which answers ExitOk. A middleware that does neither has declined, so
// the next command runs; a strict command that does neither ran silently, and
// answers ExitOk. A handler that returns an error ends the chain too, answered
// or not, through HandleFailure. Every command of one line is handed one
// CommandProps and one Consumed, which is how a middleware hands what it
// learned, and the tokens it read, to the commands after it.
//
// Nothing here prints itself. Every way a line can end without a command
// answering it — nothing matched, a value that will not bind, a token nobody
// read, a panic — is raised through props.Cli.Fail and answered by one of the
// project's own Handle* files.
func cliMain(props opinionatedagnoscli.MainProps) int {
	response, status := tracked(props)
	run := &line{
		props:    props,
		args:     props.Args,
		consumed: make([]bool, len(props.Args)),
		response: response,
		status:   status,
	}
	if props.NewProps != nil {
		run.shared = props.NewProps()
	}

	run.answer()

	code, answered := status()
	if !answered {
		return opinionatedagnoscli.ExitFailure
	}
	return code
}

// bind copies one declaration onto this line: the command line, the shared
// Consumed, CommandProps and response.
func (run *line) bind(declared *opinionatedagnoscli.Command) *opinionatedagnoscli.Command {
	bound := bindCommand(declared)
	bound.Argv = run.args
	bound.Consumed = run.consumed
	bound.Props = run.shared
	bound.Response = run.response
	return bound
}

// answer runs the chain for one command line and, when no command answered
// it, raises the failure that says so.
func (run *line) answer() {
	defer run.recoverCommand()

	for _, declared := range run.props.Cli.Commands {
		bound := run.bind(declared)

		if !run.matches(bound) {
			continue
		}

		err := run.runCommand(bound)

		// An error ends the chain whether or not the handler answered first:
		// a command that printed part of its output and then failed has
		// failed, and the status HandleFailure sets replaces the ExitOk its
		// print implied.
		if err != nil {
			run.raise(bound, opinionatedagnoscli.FailureHandler, opinionatedagnoscli.ExitFailure, "", "", err.Error())
			return
		}
		if _, answered := run.status(); answered {
			return
		}
		// A strict command is the end of the line, not a middleware: one
		// that ran and printed nothing — it wrote a file, removed one — has
		// carried the line out, so its silence is ExitOk, never a decline
		// that would read as "unknown command".
		if bound.Strict {
			run.response.SetStatus(opinionatedagnoscli.ExitOk)
			return
		}
	}

	// A command whose verb the line starts with is the one it was meant
	// for: its args did not fit, which is a usage error naming what was
	// wrong, never "unknown command".
	if near := run.nearCommand(); near != nil {
		bound := run.bind(near)
		run.raise(bound, opinionatedagnoscli.FailureBadUsage, opinionatedagnoscli.ExitUsage, "", usageProblem(bound), "")
		return
	}

	run.failLine(opinionatedagnoscli.FailureNotFound, opinionatedagnoscli.ExitUsage, "")
}

// matches is the bound command's own Matches when it declares one,
// the lib's matcher otherwise.
func (run *line) matches(bound *opinionatedagnoscli.Command) bool {
	if bound.Matches != nil {
		return bound.Matches(bound)
	}
	return matches(run.props.ArgvDeps, bound)
}

// runCommand is the bound command's own Run when it declares
// one, the lib's binder otherwise.
func (run *line) runCommand(bound *opinionatedagnoscli.Command) error {
	if bound.Run != nil {
		return bound.Run(bound)
	}
	return run.handle(bound)
}

// nearCommand is the strict command whose verb — the longest of its
// Identifiers — the command line's segments start with, nil when none does.
func (run *line) nearCommand() *opinionatedagnoscli.Command {
	segments, _ := splitArgv(run.args)
	var near *opinionatedagnoscli.Command
	longest := 0
	for _, declared := range run.props.Cli.Commands {
		if !declared.Strict {
			continue
		}
		for _, identifier := range declared.Identifiers {
			words := strings.Fields(identifier)
			if len(words) == 0 || len(words) > len(segments) || len(words) <= longest {
				continue
			}
			if strings.Join(segments[:len(words)], " ") == identifier {
				near, longest = declared, len(words)
			}
		}
	}
	return near
}

// usageProblem says what keeps a command line from fitting the command it
// names: an arg that will not convert, one missing, or a segment too many.
func usageProblem(bound *opinionatedagnoscli.Command) string {
	segments, _ := splitArgv(bound.Argv)
	usage := fmt.Sprintf(" (usage: %s)", bound.Pattern)

	for _, arg := range bound.Args {
		values, found := argSlice(segments, arg)
		if !found {
			continue
		}
		if _, converts := argValue(arg, values); !converts {
			return fmt.Sprintf("%q is not a valid %s for <%s>%s",
				strings.Join(values, " "), argTypeName(arg.Type), arg.Id, usage)
		}
	}

	reads := 0
	for _, arg := range bound.Args {
		if arg.End < 0 {
			reads = -1
			break
		}
		if arg.End+1 > reads {
			reads = arg.End + 1
		}
	}
	if bound.Segments > 0 {
		reads = bound.Segments
	}

	if reads >= 0 && len(segments) < reads {
		for _, arg := range bound.Args {
			if _, found := argSlice(segments, arg); !found {
				return fmt.Sprintf("missing <%s>%s", arg.Id, usage)
			}
		}
	}
	if reads >= 0 && len(segments) > reads {
		return fmt.Sprintf("unexpected argument %q%s", segments[reads], usage)
	}
	return "the command line does not fit the command" + usage
}

// argTypeName is how an arg type reads in a message.
func argTypeName(kind opinionatedagnoscli.ArgType) string {
	switch kind {
	case opinionatedagnoscli.ArgInteger:
		return "integer"
	case opinionatedagnoscli.ArgNumber:
		return "number"
	case opinionatedagnoscli.ArgUuid:
		return "uuid"
	}
	return "string"
}

// failLine raises a failure that belongs to no command — nothing matched the
// line. The handler still gets a Command carrying the line and the response
// and nothing else, so every Handle* file reads the same whichever failure
// brought it there.
//
// It carries no message on purpose: the wording is the project's, filled in by
// the Handle* file through FailureOf, which is what makes editing that file
// change what the cli says.
func (run *line) failLine(kind opinionatedagnoscli.CommandFailureKind, code int, cause string) {
	command := newCommand(run.props.ArgvDeps)
	command.Argv = run.args
	command.Consumed = run.consumed
	command.Props = run.shared
	command.Response = run.response

	run.raise(command, kind, code, "", "", cause)
}

// recoverCommand turns a panicking handler into one answered command line, so
// the panic is reported once, through HandleFailure, rather than as a stack
// trace. A handler that panicked after printing has still failed: the status
// HandleFailure sets replaces the ExitOk its print implied.
func (run *line) recoverCommand() {
	failure := recover()
	if failure == nil {
		return
	}

	cause := fmt.Sprintf("command panicked: %v", failure)
	run.failLine(opinionatedagnoscli.FailureHandler, opinionatedagnoscli.ExitFailure, cause)
}

// raise records what went wrong on a bound command and hands it to the
// project's own handler for that kind — HandleNotFound, HandleBadUsage,
// HandleUnknownFlag, HandleUnexpectedArg, HandleFailure — through
// props.Cli.Fail.
func (run *line) raise(command *opinionatedagnoscli.Command, kind opinionatedagnoscli.CommandFailureKind, status int, field string, message string, cause string) error {
	return run.raiseFailure(command, &opinionatedagnoscli.CommandFailure{
		Kind:    kind,
		Status:  status,
		Field:   field,
		Message: message,
		Cause:   cause,
	})
}

// raiseFailure raises one failure already built — what a handler returned
// through Fail — on the bound command it was returned from.
func (run *line) raiseFailure(command *opinionatedagnoscli.Command, failure *opinionatedagnoscli.CommandFailure) error {
	command.Failure = failure

	// A cli built without a Fail — a command bound by hand rather than by
	// the generated registry — has no handler to reach, so the failure is
	// still reported rather than swallowed.
	if run.props.Cli.Fail == nil {
		return fmt.Errorf("%s", failure.Message)
	}

	return run.props.Cli.Fail(command)
}

// tracked builds the one response of a command line, so the dispatch can tell
// whether a command answered it. It returns the response to hand to the
// commands and a reader of the status it was answered with, and whether it
// was answered at all.
//
// Answering is what ends a chain, and there are two ways to answer: setting a
// status, or printing to stdout — which answers ExitOk, the status a command
// that printed its result ends with. Error and Log answer nothing, which is how
// a middleware says something and hands the command line on. Every print goes
// through props.StdDeps at the moment it is made, so a middleware silencing
// StdDeps.Logf silences Log here too.
//
// The first status set is the one the line exits with. A print only implies
// ExitOk, so a status set after it — a failure the same command raises once it
// has printed part of its output — still replaces it: a command that printed
// and then failed never exits 0.
func tracked(props opinionatedagnoscli.MainProps) (*opinionatedagnoscli.CommandResponse, func() (int, bool)) {
	status := 0
	answered := false
	explicit := false

	response := &opinionatedagnoscli.CommandResponse{}
	response.SetStatus = func(code int) {
		if explicit {
			return
		}
		status, answered, explicit = code, true, true
	}
	response.Printf = func(format string, a ...any) (int, error) {
		if !answered {
			status, answered = opinionatedagnoscli.ExitOk, true
		}
		return props.StdDeps.Printf(format, a...)
	}
	response.Eprintf = func(format string, a ...any) (int, error) {
		return props.StdDeps.Eprintf(format, a...)
	}
	response.Logf = func(format string, a ...any) (int, error) {
		return props.StdDeps.Logf(format, a...)
	}

	return response, func() (int, bool) {
		return status, answered
	}
}
