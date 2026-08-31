package help

import (
	"fmt"
	"slices"
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/cli"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/sandbox"
)

// ─── ANSI escape sequences ──────────────────────────────────────────────────

const (
	bold    = "\033[1m"
	dim     = "\033[2m"
	italic  = "\033[3m"
	underln = "\033[4m"
	reset   = "\033[0m"
	cyan    = "\033[36m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	magenta = "\033[35m"
	white   = "\033[97m"
	gray    = "\033[90m"
	red     = "\033[31m"
)

func NewCommand(sandbox *sandbox.SandBox) cli.CliCommand {
	return cli.CliCommand{
		ValidStartIdentifiers: []string{"help", "--help"},
		Category:              "Info",
		Args: []cli.CliArg{
			{
				Id:              "command",
				Description:     "The command to get help for",
				RequiredType:    cli.CliTypeString,
				RequiredMinSize: 0,
				RequiredMaxSize: 1,
			},
		},

		Description:     "Display help for a command",
		LongDescription: "When called without arguments, lists every available command\ngrouped by category. When called with a command name, shows\ndetailed usage, arguments, flags, and examples for that command.",
		Examples: []string{
			sandbox.Config.ProjectName + " --help",
			sandbox.Config.ProjectName + " help",
			sandbox.Config.ProjectName + " help start",
		},
		Handler: CommandHandler,
	}
}

func CommandHandler(sb any, entries cli.CliEntrys) int {
	sandbox := sb.(*sandbox.SandBox)
	command := entries.GetArgById("command")
	if len(command.Values) == 0 {
		printGeneralHelp(sandbox)
		return cli.ExitOk
	}
	chosen := command.Values[0].String()
	for _, c := range sandbox.Cli.Commands {
		if slices.Contains(c.ValidStartIdentifiers, chosen) {
			printCommandHelp(sandbox, &c)
			return cli.ExitOk
		}
	}

	p := sandbox.Deps.Printf
	p("\n")
	p("  %s%s✘%s Unknown command: %s%s%s\n", bold, red, reset, bold+white, chosen, reset)
	p("  %sRun '%s help' to see available commands.%s\n", dim, sandbox.Config.ProjectName, reset)
	p("\n")
	return cli.ExitUsage
}

// ═════════════════════════════════════════════════════════════════════════════
//  General help
// ═════════════════════════════════════════════════════════════════════════════

func printGeneralHelp(sandbox *sandbox.SandBox) {
	p := sandbox.Deps.Printf

	// ── Banner ───────────────────────────────────────────────────────
	printBanner(sandbox)

	// ── Usage ────────────────────────────────────────────────────────
	p("  %s%sUSAGE%s\n", bold, cyan, reset)
	p("  %s│%s\n", gray, reset)
	p("  %s│%s  %s$%s %s %s<command>%s %s[flags]%s %s[args]%s\n",
		gray, reset,
		dim, reset,
		sandbox.Config.ProjectName,
		green, reset,
		yellow, reset,
		dim, reset,
	)
	p("  %s│%s\n", gray, reset)
	p("\n")

	// ── Collect categories ───────────────────────────────────────────
	categoryOrder := []string{}
	categorized := map[string][]cli.CliCommand{}
	for _, cmd := range sandbox.Cli.Commands {
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

	// ── Measure widths ───────────────────────────────────────────────
	maxNameLen := 0
	for _, cmd := range sandbox.Cli.Commands {
		if cmd.Hidden {
			continue
		}
		if len(cmd.ValidStartIdentifiers) > 0 {
			n := len(cmd.ValidStartIdentifiers[0])
			if n > maxNameLen {
				maxNameLen = n
			}
		}
	}

	// ── Render categories ────────────────────────────────────────────
	for _, cat := range categoryOrder {
		cmds := categorized[cat]
		p("  %s%s%s%s\n", bold, cyan, strings.ToUpper(cat), reset)
		p("  %s│%s\n", gray, reset)
		for _, cmd := range cmds {
			if len(cmd.ValidStartIdentifiers) == 0 {
				continue
			}
			name := cmd.ValidStartIdentifiers[0]

			aliasTag := ""
			if len(cmd.ValidStartIdentifiers) > 1 {
				aliasTag = fmt.Sprintf("  %s[%s]%s",
					dim,
					strings.Join(cmd.ValidStartIdentifiers[1:], ", "),
					reset,
				)
			}

			// Dotted leader
			dotsNeeded := (maxNameLen + 20) - len(name)
			if dotsNeeded < 4 {
				dotsNeeded = 4
			}
			dots := " " + strings.Repeat("·", dotsNeeded-2) + " "

			p("  %s│%s  %s%s%s%s%s%s%s%s\n",
				gray, reset,
				green+bold, name, reset,
				gray, dots, reset,
				cmd.Description,
				aliasTag,
			)
		}
		p("  %s│%s\n", gray, reset)
		p("\n")
	}

	// ── Footer ───────────────────────────────────────────────────────
	p("  %s%s─── %sTip%s%s ──────────────────────────────────────────────────%s\n",
		dim, gray, italic, reset+dim+gray, gray, reset,
	)
	p("  %sRun %s%s help <command>%s%s for detailed info on any command.%s\n",
		dim, reset+cyan, sandbox.Config.ProjectName, reset, dim, reset,
	)
	p("\n")
}

// ═════════════════════════════════════════════════════════════════════════════
//  Per-command help
// ═════════════════════════════════════════════════════════════════════════════

func printCommandHelp(sandbox *sandbox.SandBox, cmd *cli.CliCommand) {
	p := sandbox.Deps.Printf

	name := cmd.ValidStartIdentifiers[0]

	// ── Header box ───────────────────────────────────────────────────
	titleLine := fmt.Sprintf("%s %s", sandbox.Config.ProjectName, name)
	boxW := len(titleLine) + 6
	if boxW < 44 {
		boxW = 44
	}
	innerW := boxW - 2 // inside the box walls

	p("\n")
	p("  %s╭%s╮%s\n", cyan, strings.Repeat("─", innerW), reset)
	p("  %s│%s  %s%s%s%s%s│%s\n",
		cyan, reset,
		bold+white, titleLine,
		reset, strings.Repeat(" ", innerW-2-len(titleLine)),
		cyan, reset,
	)
	p("  %s│%s  %s%s%s%s%s│%s\n",
		cyan, reset,
		dim, cmd.Description,
		reset, strings.Repeat(" ", innerW-2-len(cmd.Description)),
		cyan, reset,
	)
	p("  %s╰%s╯%s\n", cyan, strings.Repeat("─", innerW), reset)
	p("\n")

	// ── Long description ─────────────────────────────────────────────
	if cmd.LongDescription != "" {
		for _, line := range strings.Split(cmd.LongDescription, "\n") {
			p("  %s%s%s\n", dim, line, reset)
		}
		p("\n")
	}

	// ── Usage ────────────────────────────────────────────────────────
	printSection(p, "USAGE")
	usage := fmt.Sprintf("  %s$%s %s %s",
		dim, reset, sandbox.Config.ProjectName, name,
	)
	flagPart := ""
	if len(cmd.Flags) > 0 {
		flagPart = fmt.Sprintf(" %s[flags]%s", yellow, reset)
	}
	argPart := ""
	for _, arg := range cmd.Args {
		if arg.RequiredMinSize > 0 {
			argPart += fmt.Sprintf(" %s%s<%s>%s", bold, green, arg.Id, reset)
		} else {
			argPart += fmt.Sprintf(" %s[%s]%s", dim, arg.Id, reset)
		}
	}
	p("  %s│%s%s%s%s\n", gray, reset, usage, flagPart, argPart)
	p("  %s│%s\n", gray, reset)
	p("\n")

	// ── Aliases ──────────────────────────────────────────────────────
	if len(cmd.ValidStartIdentifiers) > 1 {
		printSection(p, "ALIASES")
		for _, alias := range cmd.ValidStartIdentifiers {
			bullet := gray + "◦" + reset
			if alias == name {
				bullet = green + "●" + reset
			}
			p("  %s│%s  %s %s%s%s\n", gray, reset, bullet, cyan, alias, reset)
		}
		p("  %s│%s\n", gray, reset)
		p("\n")
	}

	// ── Arguments ────────────────────────────────────────────────────
	if len(cmd.Args) > 0 {
		printSection(p, "ARGUMENTS")
		for i, arg := range cmd.Args {
			reqLabel := dim + "optional" + reset
			if arg.RequiredMinSize > 0 {
				reqLabel = yellow + bold + "required" + reset
			}
			typeName := cliTypeName(arg.RequiredType)

			p("  %s│%s  %s%s%s\n",
				gray, reset,
				green+bold, arg.Id, reset,
			)
			p("  %s│%s    %s\n", gray, reset, arg.Description)
			p("  %s│%s    %s%s%s %s│%s %s\n",
				gray, reset,
				magenta, typeName, reset,
				gray, reset,
				reqLabel,
			)
			if len(arg.Defaults) > 0 {
				p("  %s│%s    %sdefault:%s %s%s%s\n",
					gray, reset,
					dim, reset,
					white+bold, strings.Join(arg.Defaults, ", "), reset,
				)
			}
			if len(arg.Examples) > 0 {
				for _, ex := range arg.Examples {
					p("  %s│%s    %s$ %s%s\n", gray, reset, dim, ex, reset)
				}
			}
			if i < len(cmd.Args)-1 {
				p("  %s│%s\n", gray, reset)
			}
		}
		p("  %s│%s\n", gray, reset)
		p("\n")
	}

	// ── Flags ────────────────────────────────────────────────────────
	if len(cmd.Flags) > 0 {
		printSection(p, "FLAGS")
		for i, flag := range cmd.Flags {
			ids := strings.Join(flag.ValidIdentifiers, gray+", "+reset+yellow+bold)
			reqLabel := dim + "optional" + reset
			if flag.RequiredPresence {
				reqLabel = yellow + bold + "required" + reset
			}
			typeName := cliTypeName(flag.Type)

			p("  %s│%s  %s%s%s\n",
				gray, reset,
				yellow+bold, ids, reset,
			)
			p("  %s│%s    %s\n", gray, reset, flag.Description)
			p("  %s│%s    %s%s%s %s│%s %s\n",
				gray, reset,
				magenta, typeName, reset,
				gray, reset,
				reqLabel,
			)
			if len(flag.Defaults) > 0 {
				p("  %s│%s    %sdefault:%s %s%s%s\n",
					gray, reset,
					dim, reset,
					white+bold, strings.Join(flag.Defaults, ", "), reset,
				)
			}
			if len(flag.Examples) > 0 {
				for _, ex := range flag.Examples {
					p("  %s│%s    %s$ %s%s\n", gray, reset, dim, ex, reset)
				}
			}
			if i < len(cmd.Flags)-1 {
				p("  %s│%s\n", gray, reset)
			}
		}
		p("  %s│%s\n", gray, reset)
		p("\n")
	}

	// ── Examples ─────────────────────────────────────────────────────
	if len(cmd.Examples) > 0 {
		printSection(p, "EXAMPLES")
		for _, ex := range cmd.Examples {
			p("  %s│%s  %s$%s %s\n", gray, reset, dim, reset, ex)
		}
		p("  %s│%s\n", gray, reset)
		p("\n")
	}
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

// printBanner renders the top box with project name and version.
func printBanner(sandbox *sandbox.SandBox) {
	p := sandbox.Deps.Printf

	titleLine := fmt.Sprintf("%s  %s", sandbox.Config.ProjectName, sandbox.Config.Version)
	boxW := len(titleLine) + 6
	if boxW < 44 {
		boxW = 44
	}
	innerW := boxW - 2

	p("\n")
	p("  %s╭%s╮%s\n", cyan, strings.Repeat("─", innerW), reset)
	p("  %s│%s  %s%s%s%s%s│%s\n",
		cyan, reset,
		bold+white, titleLine,
		reset, strings.Repeat(" ", innerW-2-len(titleLine)),
		cyan, reset,
	)
	p("  %s╰%s╯%s\n", cyan, strings.Repeat("─", innerW), reset)
	p("\n")
}

// printSection renders a colored section header with a connecting vertical
// bar below it.
func printSection(p func(string, ...any) (int, error), title string) {
	p("  %s%s%s%s\n", bold+cyan, title, reset, "")
	p("  %s│%s\n", gray, reset)
}

// cliTypeName returns a human-readable label for a CLI type constant.
func cliTypeName(t int) string {
	switch t {
	case cli.CliTypeInt:
		return "int"
	case cli.CliTypeFloat:
		return "float"
	case cli.CliTypeBool:
		return "bool"
	default:
		return "string"
	}
}
