package interview

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
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

// printWelcome opens the session: who is asking, which project the answers
// will be applied to, and what state that project is in — the same reading the
// menus are filtered by, said once in plain words so the first menu is not a
// surprise.
func printWelcome(sandbox *api.Sandbox, io *smartio.SmartIO, path string) {
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
	p("  %sYou do not need to know any command to use this.%s\n", bold+white, reset)
	p("  %sPick what you want from the menu, answer the questions, and I run it.%s\n", dim, reset)
	p("  %sNothing happens before you see the command and say yes.%s\n", dim, reset)
	p("\n")
	p("  %s↑↓%s move   %senter%s choose   %sesc%s go back   %sctrl-c%s quit\n", cyan, reset, cyan, reset, cyan, reset, cyan, reset)
	p("  %sgoing back re-asks the question before this one; nothing is lost by it%s\n", dim, reset)
	p("\n")

	state := readState(sandbox, io)

	p("  %sfolder %s  %s%s%s\n", dim, reset, green, path, reset)
	p("  %sproject%s  %s\n", dim, reset, projectText(sandbox, state))
	p("\n")
}

// projectText is the one line saying what is in the folder: nothing yet, or
// the project's name followed by the layers it has turned on. It is what makes
// the first menu legible — an area is missing because the layer behind it is
// off, and this line is where that is said.
func projectText(sandbox *api.Sandbox, state projectState) string {
	if !state.Started {
		return sandbox.Deps.Std.Sprintf("%snone here yet — creating one is the first step%s", yellow, reset)
	}

	name := state.Name
	if name == "" {
		name = "unnamed"
	}

	text := sandbox.Deps.Std.Sprintf("%s%s%s  ", bold+white, name, reset)
	for _, spec := range utils.ExtensionCatalog() {
		if extensionInit[spec.Name] == "" {
			continue
		}

		mark := sandbox.Deps.Std.Sprintf("%soff%s", red, reset)
		if enabled(state, spec.Name) {
			mark = sandbox.Deps.Std.Sprintf("%son%s", green, reset)
		}
		text += sandbox.Deps.Std.Sprintf("  %s%s%s %s", dim, layerName(sandbox, spec.Name), reset, mark)
	}

	return text
}

// layerName is a mechanic without its sandbox- prefix, which is the generator's
// word for it and not the person's.
func layerName(sandbox *api.Sandbox, extension string) string {
	return sandbox.Deps.Stringsdeps.TrimPrefix(extension, "sandbox-")
}

// printPlan is the confirm screen: the command line the answers add up to,
// spelled exactly as it would be typed. It is what turns an interview into a
// way of learning the cli rather than a way of avoiding it.
//
// Four things are said around that line, and each of them is there because the
// line alone would mislead: that the attempt before this one failed and the
// answers were kept, that a name is written down differently from the way it
// was typed, that questions were left out and why, and that running this takes
// a layer away.
func printPlan(sandbox *api.Sandbox, io *smartio.SmartIO, command api.Command, values map[string][]any, failed int) {
	p := sandbox.Deps.Std.Printf

	p("\n")

	if failed != api.ExitOk {
		p("  %sTHAT DID NOT WORK%s %s— every answer is still here; change the one at fault and run it again%s\n", bold+yellow, reset, dim, reset)
		p("\n")
	}

	if Destructive(sandbox, command, values) {
		p("  %sTHIS TAKES SOMETHING AWAY%s %s— %s%s\n", bold+red, reset, dim, lossText(sandbox, io, command), reset)
		p("\n")
	}

	p("  %sNOTHING HAS RUN YET%s %s— this is the command your answers add up to:%s\n", bold+cyan, reset, dim, reset)
	p("  %s│%s\n", gray, reset)
	p("  %s│%s  %s$%s %s%s%s\n", gray, reset, dim, reset, bold+white, CommandLine(sandbox, command, values), reset)

	for _, note := range NormalizedNotes(sandbox, command, values) {
		p("  %s│%s  %s%s%s\n", gray, reset, yellow, note, reset)
	}

	for _, note := range NotAskedNotes(sandbox, command, values) {
		p("  %s│%s  %s%s%s\n", gray, reset, dim, note, reset)
	}

	p("  %s│%s\n", gray, reset)
	p("  %s│%s  %syou could have typed it yourself — that is all I am doing%s\n", gray, reset, dim, reset)
	p("\n")
}

// lossText names what a destructive command is about to remove. A purge knows
// the units it takes with it, and naming them is the difference between a
// person agreeing to a command line and agreeing to losing four routes they
// wrote.
func lossText(sandbox *api.Sandbox, io *smartio.SmartIO, command api.Command) string {
	unit, names := Losses(sandbox, io, command)
	if len(names) == 0 {
		return "there is no undoing this one"
	}

	return sandbox.Deps.Std.Sprintf("this removes %d %s: %s",
		len(names), plural(unit, len(names)), sandbox.Deps.Stringsdeps.Join(names, ", "))
}

// plural is a unit said of however many there are of it.
func plural(unit string, count int) string {
	if count == 1 {
		return unit
	}
	return unit + "s"
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
