package opinatedagnoscli

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	opinatedagnoscli "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/OpinatedAgnosCli"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/argvdeps"
)

// Bind fills deps.Deps.OpinatedAgnosCli with the agnos cli mechanic: the
// dispatch chain of climain.go, the binder of command_handler.go, the matcher
// of is_actionable.go and trigger.go. Nothing here holds a dep: what a run
// reaches the outside world through comes in its MainProps, read as it runs.
func Bind(deps *deps.Deps) {
	deps.OpinatedAgnosCli = opinatedagnoscli.Sandbox{
		CliMain:       cliMain,
		NewCommand:    newCommand,
		BindCommand:   bindCommand,
		Fail:          fail,
		FailWithCause: failWithCause,
		FailureOf:     failureOf,
		MatchTrigger:  matchTrigger,
	}
}

// newCommand returns an empty Command with every slice open and IsActionable
// closed over the parser it reads flags with, so a caller holding a command —
// the help screen — matches it the way the dispatch does.
func newCommand(parsers argvdeps.Sandbox) *opinatedagnoscli.Command {
	command := &opinatedagnoscli.Command{
		Identifiers: []string{},
		Strict:      true,
		Examples:    []string{},
		Args:        []opinatedagnoscli.CommandArg{},
		Flags:       []opinatedagnoscli.CommandFlag{},
	}
	command.IsActionable = func(bound *opinatedagnoscli.Command) bool {
		return isActionable(parsers, bound)
	}
	return command
}

// bindCommand copies one declaration into the command a single command line
// runs on, with no command line, response or failure yet. The dispatch calls
// it once per command of the chain, so what Cli.Commands holds is never
// written to.
func bindCommand(command *opinatedagnoscli.Command) *opinatedagnoscli.Command {
	bound := *command
	bound.Argv = nil
	bound.Consumed = nil
	bound.Props = nil
	bound.Response = nil
	bound.Failure = nil
	return &bound
}

// fail builds the failure an InternalPureHandler refuses a command line with.
func fail(status int, field string, message string) error {
	return failWithCause(status, field, message, "")
}

// failWithCause is fail carrying what went wrong underneath.
func failWithCause(status int, field string, message string, cause string) error {
	return &opinatedagnoscli.CommandFailure{
		Kind:    opinatedagnoscli.HandlerFailure,
		Status:  status,
		Field:   field,
		Message: message,
		Cause:   cause,
	}
}

// failureOf merges the failure a Handle* file is answering with that file's
// own status and wording: what the failure carries wins, what it leaves empty
// the file fills.
func failureOf(command *opinatedagnoscli.Command, status int, message string) opinatedagnoscli.CommandFailure {
	if command.Failure == nil {
		return opinatedagnoscli.CommandFailure{Status: status, Message: message}
	}

	failure := *command.Failure
	if failure.Status == 0 {
		failure.Status = status
	}
	if failure.Message == "" {
		failure.Message = message
	}
	return failure
}
