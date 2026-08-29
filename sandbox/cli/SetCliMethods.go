package cli

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/cli/commands/build"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/cli/commands/extensionhelp"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/cli/commands/help"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/cli/commands/install"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/cli/commands/list"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/cli/commands/start"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/cli/commands/uninstall"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/cli/commands/version"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib"
)

func SetCliMethods(sandbox *lib.SandBox) {

	sandbox.Cli.Commands = []lib.CliCommand{
		version.NewCommand(sandbox),
		help.NewCommand(sandbox),
		start.NewCommand(sandbox),
		build.NewCommand(sandbox),
		install.NewCommand(sandbox),
		uninstall.NewCommand(sandbox),
		list.NewCommand(sandbox),
		extensionhelp.NewCommand(sandbox),
	}

	sandbox.Cli.CliMain = CliMainFactory(sandbox)
}
