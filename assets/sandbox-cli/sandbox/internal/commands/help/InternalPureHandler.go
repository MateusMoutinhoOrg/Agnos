package help

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/internal/commandprops"
)

// help is a command like any other — command.yaml, generated new.go and
// entries.go, and this InternalPureHandler.go — except that
// `{{.GeneratorName}} build` writes them instead of the user writing two of
// them. Nothing about the command set is baked in here: the screens below are
// printed from Cli.Commands, the same declarations the dispatch binds a command
// line against.

// identifiedBy reports whether name is one of the identifiers a command
// answers to, its aliases included.
func identifiedBy(identifiers []string, name string) bool {
	for _, identifier := range identifiers {
		if identifier == name {
			return true
		}
	}
	return false
}

// binaryName is the executable's name as a user types it: the configured
// project name, lowercased. Usage lines show what to type, not the display
// name of the project.
func binaryName(sandbox *api.Sandbox) string {
	return sandbox.Deps.Stringsdeps.ToLower(sandbox.Config.ProjectName)
}

// ─── ANSI escape sequences ──────────────────────────────────────────────────

const (
	bold    = "\033[1m"
	dim     = "\033[2m"
	italic  = "\033[3m"
	reset   = "\033[0m"
	cyan    = "\033[36m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	magenta = "\033[35m"
	white   = "\033[97m"
	gray    = "\033[90m"
	red     = "\033[31m"
)

// ─── Entry points ───────────────────────────────────────────────────────────

// InternalPureHandler backs `help`: with no argument it prints the general help
// screen, with a command name it prints that command's detailed help.
func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	name := sandbox.Deps.Stringsdeps.Join(entries.Name, " ")
	if name == "" {
		PrintGeneralHelp(sandbox, response)
		return nil
	}

	if declared := FindCommand(sandbox, name); declared != nil {
		PrintCommandHelp(sandbox, response, declared)
		return nil
	}

	response.SetStatus(api.ExitUsage)
	e := response.Error
	e("\n")
	e("  %s%s✘%s Unknown command: %s%s%s\n", bold, red, reset, bold+white, name, reset)
	e("  %sRun '%s help' to see available commands.%s\n", dim, binaryName(sandbox), reset)
	e("\n")
	return nil
}

// FindCommand is the command one of whose Identifiers is name — or, failing
// that, whose package is (`echo_all`, or `echo-all` for it) — nil when none is.
func FindCommand(sandbox *api.Sandbox, name string) *api.Command {
	for _, declared := range sandbox.Cli.Commands {
		if identifiedBy(declared.Identifiers, name) {
			return declared
		}
	}
	pkg := sandbox.Deps.Stringsdeps.ReplaceAll(name, "-", "_")
	for _, declared := range sandbox.Cli.Commands {
		if declared.Name == pkg {
			return declared
		}
	}
	return nil
}

// ─── General help ──────────────────────────────────────────────────────────

// PrintGeneralHelp lists every command grouped by category — the middlewares
// left out, since nobody types one. It is also the usage screen shown when the
// binary is run with no arguments.
func PrintGeneralHelp(sandbox *api.Sandbox, response *api.CommandResponse) {
	p := response.Printf

	printBanner(sandbox, response)

	p("  %s%sUSAGE%s\n", bold, cyan, reset)
	p("  %s│%s\n", gray, reset)
	p("  %s│%s  %s$%s %s %s<command>%s %s[args]%s %s[flags]%s\n",
		gray, reset, dim, reset, binaryName(sandbox),
		green, reset, dim, reset, yellow, reset,
	)
	p("  %s│%s\n", gray, reset)
	p("\n")

	categoryOrder := []string{}
	categorized := map[string][]*api.Command{}
	for _, cmd := range sandbox.Cli.Commands {
		if !listed(cmd) {
			continue
		}
		cat := cmd.Category
		if cat == "" {
			cat = "Other"
		}
		if _, exists := categorized[cat]; !exists {
			categoryOrder = append(categoryOrder, cat)
		}
		categorized[cat] = append(categorized[cat], cmd)
	}

	maxNameLen := 0
	for _, cmd := range sandbox.Cli.Commands {
		if !listed(cmd) {
			continue
		}
		if n := len(cmd.Identifiers[0]); n > maxNameLen {
			maxNameLen = n
		}
	}

	for _, cat := range categoryOrder {
		p("  %s%s%s%s\n", bold, cyan, sandbox.Deps.Stringsdeps.ToUpper(cat), reset)
		p("  %s│%s\n", gray, reset)
		for _, cmd := range categorized[cat] {
			name := cmd.Identifiers[0]

			aliasTag := ""
			if len(cmd.Identifiers) > 1 {
				aliasTag = sandbox.Deps.Std.Sprintf("  %s[%s]%s", dim, sandbox.Deps.Stringsdeps.Join(cmd.Identifiers[1:], ", "), reset)
			}

			dotsNeeded := (maxNameLen + 20) - len(name)
			if dotsNeeded < 4 {
				dotsNeeded = 4
			}
			dots := " " + sandbox.Deps.Stringsdeps.Repeat("·", dotsNeeded-2) + " "

			p("  %s│%s  %s%s%s%s%s%s%s%s\n",
				gray, reset, green+bold, name, reset, gray, dots, reset, cmd.Help, aliasTag,
			)
		}
		p("  %s│%s\n", gray, reset)
		p("\n")
	}

	p("  %s%s─── %sTip%s%s ──────────────────────────────%s\n",
		dim, gray, italic, reset+dim+gray, gray, reset,
	)
	p("  %sRun %s%s help <command>%s%s for detailed info on any command.%s\n",
		dim, reset+cyan, binaryName(sandbox), reset, dim, reset,
	)
	p("\n")
}

// listed reports a command the general help lists: a visible, strict one
// with a verb to type.
func listed(cmd *api.Command) bool {
	return !cmd.Hidden && cmd.Strict && len(cmd.Identifiers) > 0
}

// ─── Per-command help ──────────────────────────────────────────────────────

// PrintCommandHelp prints one command's detailed help: its description, the
// line to type, and every arg and flag it reads — the flags of the middlewares
// that run in front of it included, since the user types those too.
func PrintCommandHelp(sandbox *api.Sandbox, response *api.CommandResponse, cmd *api.Command) {
	p := response.Printf

	name := cmd.Pattern
	if len(cmd.Identifiers) > 0 {
		name = cmd.Identifiers[0]
	}

	titleLine := sandbox.Deps.Std.Sprintf("%s %s", binaryName(sandbox), name)
	innerW := len(titleLine) + 4
	if w := len(cmd.Help) + 4; w > innerW {
		innerW = w
	}
	if innerW < 42 {
		innerW = 42
	}

	p("\n")
	p("  %s╭%s╮%s\n", cyan, sandbox.Deps.Stringsdeps.Repeat("─", innerW), reset)
	p("  %s│%s  %s%s%s%s%s│%s\n",
		cyan, reset, bold+white, titleLine, reset,
		sandbox.Deps.Stringsdeps.Repeat(" ", innerW-2-len(titleLine)), cyan, reset,
	)
	p("  %s│%s  %s%s%s%s%s│%s\n",
		cyan, reset, dim, cmd.Help, reset,
		sandbox.Deps.Stringsdeps.Repeat(" ", innerW-2-len(cmd.Help)), cyan, reset,
	)
	p("  %s╰%s╯%s\n", cyan, sandbox.Deps.Stringsdeps.Repeat("─", innerW), reset)
	p("\n")

	if cmd.LongDescription != "" {
		for _, line := range sandbox.Deps.Stringsdeps.Split(cmd.LongDescription, "\n") {
			p("  %s%s%s\n", dim, line, reset)
		}
		p("\n")
	}

	flags := InheritedFlags(sandbox, cmd)

	printSection(p, "USAGE")
	flagPart := ""
	if len(cmd.Flags)+len(flags) > 0 {
		flagPart = sandbox.Deps.Std.Sprintf(" %s[flags]%s", yellow, reset)
	}
	p("  %s│%s  %s$%s %s %s%s\n", gray, reset, dim, reset, binaryName(sandbox), cmd.Pattern, flagPart)
	p("  %s│%s\n", gray, reset)
	p("\n")

	if len(cmd.Identifiers) > 1 {
		printSection(p, "ALIASES")
		for _, alias := range cmd.Identifiers {
			bullet := gray + "◦" + reset
			if alias == name {
				bullet = green + "●" + reset
			}
			p("  %s│%s  %s %s%s%s\n", gray, reset, bullet, cyan, alias, reset)
		}
		p("  %s│%s\n", gray, reset)
		p("\n")
	}

	args := []api.CommandArg{}
	for _, arg := range cmd.Args {
		if !arg.Trigger.Exist {
			args = append(args, arg)
		}
	}
	if len(args) > 0 {
		printSection(p, "ARGUMENTS")
		for i, arg := range args {
			printField(p, arg.Id, arg.Description, argTypeLabel(arg), arg.Default, arg.Required, "")
			if i < len(args)-1 {
				p("  %s│%s\n", gray, reset)
			}
		}
		p("  %s│%s\n", gray, reset)
		p("\n")
	}

	all := []inheritedFlag{}
	for _, flag := range cmd.Flags {
		all = append(all, inheritedFlag{Flag: flag})
	}
	all = append(all, flags...)
	if len(all) > 0 {
		printSection(p, "FLAGS")
		for i, flag := range all {
			label := sandbox.Deps.Stringsdeps.Join(flag.Flag.Keys, gray+", "+reset+yellow+bold)
			from := ""
			if flag.From != "" {
				from = "from " + flag.From
			}
			printField(p, label, flag.Flag.Description, flagTypeLabel(flag.Flag.Type), flag.Flag.Default, flag.Flag.Required, from)
			if i < len(all)-1 {
				p("  %s│%s\n", gray, reset)
			}
		}
		p("  %s│%s\n", gray, reset)
		p("\n")
	}

	if len(cmd.Examples) > 0 {
		printSection(p, "EXAMPLES")
		for _, ex := range cmd.Examples {
			p("  %s│%s  %s$%s %s %s\n", gray, reset, dim, reset, binaryName(sandbox), ex)
		}
		p("  %s│%s\n", gray, reset)
		p("\n")
	}
}

// inheritedFlag is one flag a command's help lists, and the middleware that
// declares it — "" for one of the command's own.
type inheritedFlag struct {
	Flag api.CommandFlag
	From string
}

// InheritedFlags is every flag a middleware in front of the command declares:
// a non-strict command on a lower rung whose args match the command's own
// literal verb. What the user types after the command is not known here, so a
// middleware whose trigger reads a flag value is listed as well.
func InheritedFlags(sandbox *api.Sandbox, cmd *api.Command) []inheritedFlag {
	inherited := []inheritedFlag{}
	argv := []string{}
	if len(cmd.Identifiers) > 0 {
		argv = sandbox.Deps.Stringsdeps.Fields(cmd.Identifiers[0])
	}

	for _, declared := range sandbox.Cli.Commands {
		if declared.Strict || declared.Priority >= cmd.Priority || declared.Name == cmd.Name {
			continue
		}
		bound := api.BindCommand(declared)
		bound.Argv = argv
		bound.Flags = []api.CommandFlag{}
		if !declared.IsActionable(bound) {
			continue
		}
		for _, flag := range declared.Flags {
			if declaresKey(cmd, flag.Keys) {
				continue
			}
			inherited = append(inherited, inheritedFlag{Flag: flag, From: declared.Name})
		}
	}
	return inherited
}

// declaresKey reports whether a command declares a flag under one of keys.
func declaresKey(cmd *api.Command, keys []string) bool {
	for _, flag := range cmd.Flags {
		for _, own := range flag.Keys {
			for _, key := range keys {
				if own == key {
					return true
				}
			}
		}
	}
	return false
}

// ─── Helpers ───────────────────────────────────────────────────────────────

func printField(p func(string, ...any) (int, error), label, description, kind, def string, required bool, from string) {
	reqLabel := dim + "optional" + reset
	if required {
		reqLabel = yellow + bold + "required" + reset
	}

	p("  %s│%s  %s%s%s\n", gray, reset, green+bold, label, reset)
	if description != "" {
		p("  %s│%s    %s\n", gray, reset, description)
	}
	p("  %s│%s    %s%s%s %s│%s %s\n",
		gray, reset, magenta, kind, reset, gray, reset, reqLabel,
	)
	if def != "" {
		p("  %s│%s    %sdefault:%s %s%s%s\n", gray, reset, dim, reset, white+bold, def, reset)
	}
	if from != "" {
		p("  %s│%s    %s%s%s\n", gray, reset, dim, from, reset)
	}
}

func printBanner(sandbox *api.Sandbox, response *api.CommandResponse) {
	p := response.Printf

	titleLine := sandbox.Deps.Std.Sprintf("%s  %s", sandbox.Config.ProjectName, sandbox.Config.Version)
	innerW := len(titleLine) + 4
	if innerW < 42 {
		innerW = 42
	}

	p("\n")
	p("  %s╭%s╮%s\n", cyan, sandbox.Deps.Stringsdeps.Repeat("─", innerW), reset)
	p("  %s│%s  %s%s%s%s%s│%s\n",
		cyan, reset, bold+white, titleLine, reset,
		sandbox.Deps.Stringsdeps.Repeat(" ", innerW-2-len(titleLine)), cyan, reset,
	)
	p("  %s╰%s╯%s\n", cyan, sandbox.Deps.Stringsdeps.Repeat("─", innerW), reset)
	p("\n")
}

func printSection(p func(string, ...any) (int, error), title string) {
	p("  %s%s%s\n", bold+cyan, title, reset)
	p("  %s│%s\n", gray, reset)
}

// argTypeLabel spells an arg's type the way the help screen prints it.
func argTypeLabel(arg api.CommandArg) string {
	if arg.End != arg.Start {
		return "string..."
	}
	switch arg.Type {
	case api.IntegerArg:
		return "integer"
	case api.NumberArg:
		return "number"
	case api.UuidArg:
		return "uuid"
	}
	return "string"
}

// flagTypeLabel spells a flag's type the way the help screen prints it.
func flagTypeLabel(kind api.FlagType) string {
	switch kind {
	case api.IntegerFlag:
		return "integer"
	case api.NumberFlag:
		return "number"
	case api.BooleanFlag:
		return "boolean"
	case api.StringArrayFlag:
		return "string, repeatable"
	case api.IntegerArrayFlag:
		return "integer, repeatable"
	}
	return "string"
}
