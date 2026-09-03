package help

import (
	"fmt"
	"slices"
	"strings"

	"{{.Module}}/sandbox/deps"
	"{{.Module}}/sandbox/internal/config"
)

// help is a command like any other — entries.yaml, generated entries.go, and
// this handler.go — except that `agnos build` writes all three instead of the
// user writing two of them. This file is regenerated from every
// sandbox/internal/commands/<name>/entries.yaml: the command metadata is baked
// into helpCommands below, the rendering code is fixed.

const (
	exitOk    = 0
	exitUsage = 2
)

// binaryName is the executable's name as a user types it: the configured
// project name, lowercased. Usage lines show what to type, not the display
// name of the project.
func binaryName() string {
	return strings.ToLower(config.ProjectName)
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

// ─── Baked command metadata ─────────────────────────────────────────────────

type helpField struct {
	Identifiers []string // empty for a positional argument
	Name        string   // set for a positional argument
	Description string
	Examples    []string
	Type        string
	Default     string
	Required    bool
}

type helpCommand struct {
	Identifiers     []string
	Category        string
	Description     string
	LongDescription string
	Examples        []string
	Hidden          bool
	Flags           []helpField
	Args            []helpField
}

var helpCommands = []helpCommand{
{{- range .Commands}}
	{
		Identifiers:     []string{ {{.IdentifiersGo}} },
		Category:        {{printf "%q" .Category}},
		Description:     {{printf "%q" .Help}},
		LongDescription: {{printf "%q" .LongDescription}},
		Examples:        []string{ {{range .Examples}}{{printf "%q" .}}, {{end}} },
		Hidden:          {{.Hidden}},
		Flags: []helpField{
{{- range .Flags}}
			{Identifiers: []string{ {{.IdentifiersGo}} }, Description: {{printf "%q" .Description}}, Examples: []string{ {{range .Examples}}{{printf "%q" .}}, {{end}} }, Type: {{printf "%q" .Type}}, Default: {{printf "%q" .Default}}, Required: {{.Required}}},
{{- end}}
		},
		Args: []helpField{
{{- range .Args}}
			{Name: {{printf "%q" .Key}}, Description: {{printf "%q" .Description}}, Examples: []string{ {{range .Examples}}{{printf "%q" .}}, {{end}} }, Type: {{printf "%q" .Type}}, Default: {{printf "%q" .Default}}, Required: {{.Required}}},
{{- end}}
		},
	},
{{- end}}
}

// ─── Entry points ───────────────────────────────────────────────────────────

// CommandHandler backs the `help` / `--help` verb: with no argument it prints
// the general help screen, with a command name it prints that command's
// detailed help.
func CommandHandler(deps *deps.Deps, entries *Entries) int {
	name := entries.Command
	if name == "" {
		PrintGeneralHelp(deps)
		return exitOk
	}

	for i := range helpCommands {
		if slices.Contains(helpCommands[i].Identifiers, name) {
			printCommandHelp(deps, &helpCommands[i])
			return exitOk
		}
	}

	e := deps.Std.Error
	e("\n")
	e("  %s%s✘%s Unknown command: %s%s%s\n", bold, red, reset, bold+white, name, reset)
	e("  %sRun '%s help' to see available commands.%s\n", dim, binaryName(), reset)
	e("\n")
	return exitUsage
}

// ─── General help ──────────────────────────────────────────────────────────

// PrintGeneralHelp lists every command grouped by category. It is also the
// usage screen shown when the binary is run with no arguments.
func PrintGeneralHelp(deps *deps.Deps) {
	p := deps.Std.Printf

	printBanner(deps)

	p("  %s%sUSAGE%s\n", bold, cyan, reset)
	p("  %s│%s\n", gray, reset)
	p("  %s│%s  %s$%s %s %s<command>%s %s[flags]%s %s[args]%s\n",
		gray, reset, dim, reset, binaryName(),
		green, reset, yellow, reset, dim, reset,
	)
	p("  %s│%s\n", gray, reset)
	p("\n")

	categoryOrder := []string{}
	categorized := map[string][]helpCommand{}
	for _, cmd := range helpCommands {
		if cmd.Hidden {
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
	for _, cmd := range helpCommands {
		if cmd.Hidden || len(cmd.Identifiers) == 0 {
			continue
		}
		if n := len(cmd.Identifiers[0]); n > maxNameLen {
			maxNameLen = n
		}
	}

	for _, cat := range categoryOrder {
		p("  %s%s%s%s\n", bold, cyan, strings.ToUpper(cat), reset)
		p("  %s│%s\n", gray, reset)
		for _, cmd := range categorized[cat] {
			if len(cmd.Identifiers) == 0 {
				continue
			}
			name := cmd.Identifiers[0]

			aliasTag := ""
			if len(cmd.Identifiers) > 1 {
				aliasTag = fmt.Sprintf("  %s[%s]%s", dim, strings.Join(cmd.Identifiers[1:], ", "), reset)
			}

			dotsNeeded := (maxNameLen + 20) - len(name)
			if dotsNeeded < 4 {
				dotsNeeded = 4
			}
			dots := " " + strings.Repeat("·", dotsNeeded-2) + " "

			p("  %s│%s  %s%s%s%s%s%s%s%s\n",
				gray, reset, green+bold, name, reset, gray, dots, reset, cmd.Description, aliasTag,
			)
		}
		p("  %s│%s\n", gray, reset)
		p("\n")
	}

	p("  %s%s─── %sTip%s%s ──────────────────────────────%s\n",
		dim, gray, italic, reset+dim+gray, gray, reset,
	)
	p("  %sRun %s%s help <command>%s%s for detailed info on any command.%s\n",
		dim, reset+cyan, binaryName(), reset, dim, reset,
	)
	p("\n")
}

// ─── Per-command help ──────────────────────────────────────────────────────

func printCommandHelp(deps *deps.Deps, cmd *helpCommand) {
	p := deps.Std.Printf

	name := cmd.Identifiers[0]

	titleLine := fmt.Sprintf("%s %s", binaryName(), name)
	innerW := len(titleLine) + 4
	if w := len(cmd.Description) + 4; w > innerW {
		innerW = w
	}
	if innerW < 42 {
		innerW = 42
	}

	p("\n")
	p("  %s╭%s╮%s\n", cyan, strings.Repeat("─", innerW), reset)
	p("  %s│%s  %s%s%s%s%s│%s\n",
		cyan, reset, bold+white, titleLine, reset,
		strings.Repeat(" ", innerW-2-len(titleLine)), cyan, reset,
	)
	p("  %s│%s  %s%s%s%s%s│%s\n",
		cyan, reset, dim, cmd.Description, reset,
		strings.Repeat(" ", innerW-2-len(cmd.Description)), cyan, reset,
	)
	p("  %s╰%s╯%s\n", cyan, strings.Repeat("─", innerW), reset)
	p("\n")

	if cmd.LongDescription != "" {
		for _, line := range strings.Split(cmd.LongDescription, "\n") {
			p("  %s%s%s\n", dim, line, reset)
		}
		p("\n")
	}

	printSection(p, "USAGE")
	usage := fmt.Sprintf("  %s$%s %s %s", dim, reset, binaryName(), name)
	flagPart := ""
	if len(cmd.Flags) > 0 {
		flagPart = fmt.Sprintf(" %s[flags]%s", yellow, reset)
	}
	argPart := ""
	for _, arg := range cmd.Args {
		if arg.Required {
			argPart += fmt.Sprintf(" %s%s<%s>%s", bold, green, arg.Name, reset)
		} else {
			argPart += fmt.Sprintf(" %s[%s]%s", dim, arg.Name, reset)
		}
	}
	p("  %s│%s%s%s%s\n", gray, reset, usage, flagPart, argPart)
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

	if len(cmd.Args) > 0 {
		printSection(p, "ARGUMENTS")
		for i, arg := range cmd.Args {
			printField(p, arg.Name, arg.Description, arg.Type, arg.Default, arg.Required, arg.Examples)
			if i < len(cmd.Args)-1 {
				p("  %s│%s\n", gray, reset)
			}
		}
		p("  %s│%s\n", gray, reset)
		p("\n")
	}

	if len(cmd.Flags) > 0 {
		printSection(p, "FLAGS")
		for i, flag := range cmd.Flags {
			label := strings.Join(flag.Identifiers, gray+", "+reset+yellow+bold)
			printField(p, label, flag.Description, flag.Type, flag.Default, flag.Required, flag.Examples)
			if i < len(cmd.Flags)-1 {
				p("  %s│%s\n", gray, reset)
			}
		}
		p("  %s│%s\n", gray, reset)
		p("\n")
	}

	if len(cmd.Examples) > 0 {
		printSection(p, "EXAMPLES")
		for _, ex := range cmd.Examples {
			p("  %s│%s  %s$%s %s %s\n", gray, reset, dim, reset, binaryName(), ex)
		}
		p("  %s│%s\n", gray, reset)
		p("\n")
	}
}

// ─── Helpers ───────────────────────────────────────────────────────────────

func printField(p func(string, ...any) (int, error), label, description, kind, def string, required bool, examples []string) {
	reqLabel := dim + "optional" + reset
	if required {
		reqLabel = yellow + bold + "required" + reset
	}

	p("  %s│%s  %s%s%s\n", gray, reset, green+bold, label, reset)
	p("  %s│%s    %s\n", gray, reset, description)
	p("  %s│%s    %s%s%s %s│%s %s\n",
		gray, reset, magenta, typeLabel(kind), reset, gray, reset, reqLabel,
	)
	if def != "" {
		p("  %s│%s    %sdefault:%s %s%s%s\n", gray, reset, dim, reset, white+bold, def, reset)
	}
	for _, ex := range examples {
		p("  %s│%s    %s$ %s%s\n", gray, reset, dim, ex, reset)
	}
}

func printBanner(deps *deps.Deps) {
	p := deps.Std.Printf

	titleLine := fmt.Sprintf("%s  %s", config.ProjectName, config.Version)
	innerW := len(titleLine) + 4
	if innerW < 42 {
		innerW = 42
	}

	p("\n")
	p("  %s╭%s╮%s\n", cyan, strings.Repeat("─", innerW), reset)
	p("  %s│%s  %s%s%s%s%s│%s\n",
		cyan, reset, bold+white, titleLine, reset,
		strings.Repeat(" ", innerW-2-len(titleLine)), cyan, reset,
	)
	p("  %s╰%s╯%s\n", cyan, strings.Repeat("─", innerW), reset)
	p("\n")
}

func printSection(p func(string, ...any) (int, error), title string) {
	p("  %s%s%s\n", bold+cyan, title, reset)
	p("  %s│%s\n", gray, reset)
}

func typeLabel(kind string) string {
	switch kind {
	case "int":
		return "int"
	case "float":
		return "float"
	case "boolean":
		return "bool"
	default:
		return "string"
	}
}
