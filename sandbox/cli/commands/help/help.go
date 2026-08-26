package help

import (
	"fmt"
	"slices"
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

func NewCommand(sandbox *api.SandBox) api.CliCommand {
	return api.CliCommand{
		ValidStartIdentifiers: []string{"help", "--help"},
		Args: []api.CliArg{
			api.CliArg{
				Id:              "command",
				Description:     "The command to get help for",
				RequiredType:    api.CliTypeString,
				RequiredMinSize: 0,
				RequiredMaxSize: 1,
			},
		},

		Description: "Shows Help of a command",
		Examples: []string{
			sandbox.ProjectName + " --help",
			sandbox.ProjectName + " help",
			sandbox.ProjectName + " help <command>",
		},
		Handler: CommandHandler,
	}
}

func CommandHandler(sandbox *api.SandBox, entries api.CliEntrys) int {
	command := entries.GetArgById("command")
	if len(command.Values) == 0 {
		printGeneralHelp(sandbox)
		return api.ExitOk
	}
	chosen := command.Values[0].String()
	for _, c := range sandbox.Commands {
		if slices.Contains(c.ValidStartIdentifiers, chosen) {
			printCommandHelp(sandbox, &c)
			return api.ExitOk
		}
	}

	sandbox.Deps.Printf("Unknown command: %s\n", chosen)
	sandbox.Deps.Printf("Run '%s help' to see available commands.\n", sandbox.ProjectName)
	return api.ExitUsage
}

// ─── General help ────────────────────────────────────────────────────────────

func printGeneralHelp(sandbox *api.SandBox) {
	p := sandbox.Deps.Printf

	p("%s %s\n", sandbox.ProjectName, sandbox.Version)
	p("\n")
	p("Usage:\n")
	p("  %s <command> [flags] [args]\n", sandbox.ProjectName)
	p("\n")
	p("Available Commands:\n")

	// Measure longest command name for alignment.
	maxLen := 0
	for _, cmd := range sandbox.Commands {
		if len(cmd.ValidStartIdentifiers) > 0 {
			name := cmd.ValidStartIdentifiers[0]
			if len(name) > maxLen {
				maxLen = len(name)
			}
		}
	}

	for _, cmd := range sandbox.Commands {
		if len(cmd.ValidStartIdentifiers) == 0 {
			continue
		}
		name := cmd.ValidStartIdentifiers[0]
		aliases := ""
		if len(cmd.ValidStartIdentifiers) > 1 {
			aliases = " (" + strings.Join(cmd.ValidStartIdentifiers[1:], ", ") + ")"
		}
		p("  %-*s  %s%s\n", maxLen, name, cmd.Description, aliases)
	}

	p("\n")
	p("Use \"%s help <command>\" for more information about a command.\n", sandbox.ProjectName)
}

// ─── Per-command help ────────────────────────────────────────────────────────

func printCommandHelp(sandbox *api.SandBox, cmd *api.CliCommand) {
	p := sandbox.Deps.Printf

	// ── Header ───────────────────────────────────────────────────────
	p("%s\n", cmd.Description)
	p("\n")

	// ── Usage line ───────────────────────────────────────────────────
	p("Usage:\n")
	usage := fmt.Sprintf("  %s %s", sandbox.ProjectName, cmd.ValidStartIdentifiers[0])
	if len(cmd.Flags) > 0 {
		usage += " [flags]"
	}
	for _, arg := range cmd.Args {
		if arg.RequiredMinSize > 0 {
			usage += fmt.Sprintf(" <%s>", arg.Id)
		} else {
			usage += fmt.Sprintf(" [%s]", arg.Id)
		}
	}
	p("%s\n", usage)
	p("\n")

	// ── Aliases ──────────────────────────────────────────────────────
	if len(cmd.ValidStartIdentifiers) > 1 {
		p("Aliases:\n")
		p("  %s\n", strings.Join(cmd.ValidStartIdentifiers, ", "))
		p("\n")
	}

	// ── Arguments ────────────────────────────────────────────────────
	if len(cmd.Args) > 0 {
		p("Arguments:\n")
		maxArgLen := 0
		for _, arg := range cmd.Args {
			if len(arg.Id) > maxArgLen {
				maxArgLen = len(arg.Id)
			}
		}
		for _, arg := range cmd.Args {
			req := "optional"
			if arg.RequiredMinSize > 0 {
				req = "required"
			}
			typeName := cliTypeName(arg.RequiredType)
			p("  %-*s  %s  (%s, %s)\n", maxArgLen, arg.Id, arg.Description, typeName, req)
			if len(arg.Defaults) > 0 {
				p("  %-*s  default: %s\n", maxArgLen, "", strings.Join(arg.Defaults, ", "))
			}
		}
		p("\n")
	}

	// ── Flags ────────────────────────────────────────────────────────
	if len(cmd.Flags) > 0 {
		p("Flags:\n")
		// Build the display identifiers, measure max width.
		flagLabels := make([]string, len(cmd.Flags))
		maxFlagLen := 0
		for i, flag := range cmd.Flags {
			flagLabels[i] = strings.Join(flag.ValidIdentifiers, ", ")
			if len(flagLabels[i]) > maxFlagLen {
				maxFlagLen = len(flagLabels[i])
			}
		}
		for i, flag := range cmd.Flags {
			req := "optional"
			if flag.RequiredPresence {
				req = "required"
			}
			typeName := cliTypeName(flag.Type)
			p("  %-*s  %s  (%s, %s)\n", maxFlagLen, flagLabels[i], flag.Description, typeName, req)
			if len(flag.Defaults) > 0 {
				p("  %-*s  default: %s\n", maxFlagLen, "", strings.Join(flag.Defaults, ", "))
			}
		}
		p("\n")
	}

	// ── Examples ─────────────────────────────────────────────────────
	if len(cmd.Examples) > 0 {
		p("Examples:\n")
		for _, ex := range cmd.Examples {
			p("  $ %s\n", ex)
		}
		p("\n")
	}
}

// cliTypeName returns a human-readable label for a CLI type constant.
func cliTypeName(t int) string {
	switch t {
	case api.CliTypeInt:
		return "int"
	case api.CliTypeFloat:
		return "float"
	case api.CliTypeBool:
		return "bool"
	default:
		return "string"
	}
}
