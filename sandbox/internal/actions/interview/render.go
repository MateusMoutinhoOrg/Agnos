package interview

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// ─── ANSI escape sequences ──────────────────────────────────────────────────
//
// The palette sandbox/internal/commands/help/handler.go prints its screens
// with. It is restated rather than shared: help/handler.go is rendered into
// every project from assets/sandbox-cli/, and the interview is agnos's own.

const (
	bold   = "\033[1m"
	dim    = "\033[2m"
	reset  = "\033[0m"
	cyan   = "\033[36m"
	green  = "\033[32m"
	yellow = "\033[33m"
	white  = "\033[97m"
	gray   = "\033[90m"
	red    = "\033[31m"
)

// binaryName is the executable's name as a user types it, the same spelling
// the help screens use.
func binaryName(sandbox *api.Sandbox) string {
	return sandbox.Deps.Stringsdeps.ToLower(sandbox.Config.ProjectName)
}

// notice reports one unusable answer, between the question and the next
// attempt at it.
func notice(sandbox *api.Sandbox, format string, a ...any) {
	sandbox.Deps.Std.Printf("  %s%s%s\n", yellow, sandbox.Deps.Std.Sprintf(format, a...), reset)
}

// printWelcome opens the session: who is asking, and which project the answers
// will be applied to.
func printWelcome(sandbox *api.Sandbox, path string) {
	p := sandbox.Deps.Std.Printf

	title := sandbox.Deps.Std.Sprintf("%s  %s", sandbox.Config.ProjectName, sandbox.Config.Version)
	width := len(title) + 4
	if width < 46 {
		width = 46
	}

	p("\n")
	p("  %s╭%s╮%s\n", cyan, sandbox.Deps.Stringsdeps.Repeat("─", width), reset)
	p("  %s│%s  %s%s%s%s%s│%s\n",
		cyan, reset, bold+white, title, reset,
		sandbox.Deps.Stringsdeps.Repeat(" ", width-2-len(title)), cyan, reset,
	)
	p("  %s╰%s╯%s\n", cyan, sandbox.Deps.Stringsdeps.Repeat("─", width), reset)
	p("\n")
	p("  %sWelcome to %s. I am your development assistant.%s\n", bold+white, sandbox.Config.ProjectName, reset)
	p("  %sAnswer the questions and I will run the command for you —%s\n", dim, reset)
	p("  %severy answer is a flag or an argument you could have typed.%s\n", dim, reset)
	p("\n")
	p("  %sproject%s  %s%s%s\n", dim, reset, green, path, reset)
	p("\n")
}

// printPlan is the confirm screen: the command line the answers add up to,
// spelled exactly as it would be typed. It is what turns an interview into a
// way of learning the cli rather than a way of avoiding it.
func printPlan(sandbox *api.Sandbox, command api.Command, values map[string][]any) {
	p := sandbox.Deps.Std.Printf

	p("\n")
	p("  %sABOUT TO RUN%s\n", bold+cyan, reset)
	p("  %s│%s\n", gray, reset)
	p("  %s│%s  %s$%s %s%s%s\n", gray, reset, dim, reset, bold+white, CommandLine(sandbox, command, values), reset)
	p("  %s│%s\n", gray, reset)
	p("\n")
}

// printOutcome reports what the command answered with, in the words the exit
// code carries.
func printOutcome(sandbox *api.Sandbox, command api.Command, exit int) {
	p := sandbox.Deps.Std.Printf

	if exit == api.ExitOk {
		p("\n  %s✔%s %s%s%s finished\n\n", green+bold, reset, bold+white, verbOf(command), reset)
		return
	}

	p("\n  %s✘%s %s%s%s exited %s%d%s\n\n", red+bold, reset, bold+white, verbOf(command), reset, yellow, exit, reset)
}

// CommandLine spells the bound values back as the command line that produces
// them. A flag holding exactly its declared default is left off: it changes
// nothing and only makes the line harder to read. A boolean is the flag's
// presence, so a false one is left off too.
func CommandLine(sandbox *api.Sandbox, command api.Command, values map[string][]any) string {
	line := sandbox.Deps.Std.Sprintf("%s %s", binaryName(sandbox), verbOf(command))

	for _, flag := range command.Flags {
		for _, value := range values[flag.Id] {
			if flag.Type == typeBoolean {
				if truth, ok := value.(bool); ok && truth {
					line += " " + flagName(flag)
				}
				continue
			}

			text := valueText(sandbox, value)
			if flag.HasDefault && text == flag.Default {
				continue
			}
			line += sandbox.Deps.Std.Sprintf(" %s %s", flagName(flag), quoted(sandbox, text))
		}
	}

	return line + argsText(sandbox, command, values)
}

// argsText spells the positional half of the line. Positionals bind by order,
// so every arg up to the last one answered is printed — an unanswered one as
// an empty string, which is what keeps the ones after it in their places.
func argsText(sandbox *api.Sandbox, command api.Command, values map[string][]any) string {
	last := -1
	for index, arg := range command.Args {
		if len(values[arg.Id]) > 0 {
			last = index
		}
	}

	text := ""
	for index := 0; index <= last; index++ {
		arg := command.Args[index]
		bound := values[arg.Id]

		if len(bound) == 0 {
			text += ` ""`
			continue
		}
		for _, value := range bound {
			text += " " + quoted(sandbox, valueText(sandbox, value))
		}
	}

	return text
}

// flagName is the spelling a flag is written with, the first of its
// identifiers.
func flagName(flag api.CommandFlag) string {
	if len(flag.Identifiers) == 0 {
		return "--" + flag.Id
	}
	return flag.Identifiers[0]
}

// valueText writes one bound value back as the text a command line carries it
// as.
func valueText(sandbox *api.Sandbox, value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case bool:
		if typed {
			return "true"
		}
		return "false"
	case int:
		return sandbox.Deps.Stringsdeps.FormatInt(int64(typed), 10)
	case float64:
		return sandbox.Deps.Stringsdeps.FormatFloat(typed, 'g', -1, 64)
	}
	return sandbox.Deps.Std.Sprintf("%v", value)
}

// quoted wraps a value the shell would otherwise split, so the printed line is
// one that can be pasted.
func quoted(sandbox *api.Sandbox, text string) string {
	if text == "" || sandbox.Deps.Stringsdeps.Contains(text, " ") {
		return sandbox.Deps.Std.Sprintf("%q", text)
	}
	return text
}
