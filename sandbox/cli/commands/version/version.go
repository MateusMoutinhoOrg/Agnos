package version

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/cli"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/sandbox"
)

func NewCommand(sandbox *sandbox.SandBox) cli.CliCommand {
	return cli.CliCommand{
		ValidStartIdentifiers: []string{"version", "--version"},
		Category:              "Info",
		Description:           "Print the installed version",
		LongDescription:       "Prints the current version of the installed binary and exits.",
		Examples: []string{
			sandbox.Config.ProjectName + " --version",
			sandbox.Config.ProjectName + " version",
		},
		Handler: CommandHandler,
	}
}

func CommandHandler(sb any, entries cli.CliEntrys) int {
	sandbox := sb.(*sandbox.SandBox)
	sandbox.Deps.Printf("Version:%s\n", sandbox.Config.Version)
	return cli.ExitOk
}
