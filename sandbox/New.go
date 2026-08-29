package sandbox

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/cli"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/cli/commands/help"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/cli/commands/version"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core/build"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core/extensionhelp"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core/install"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core/list"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core/start"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core/uninstall"
)

func New(d deps.Deps) sandbox.SandBox {
	sandbox := sandbox.SandBox{Deps: d}

	sandbox.Config.Version = "v0.0.12"
	sandbox.Config.ProjectName = "Agnos"

	sandbox.Cli.Commands = []cli.CliCommand{
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

	sandbox.Core.Start = start.StartFactory(sandbox)
	sandbox.Core.Build = build.BuildFactory(sandbox)
	sandbox.Core.Install = install.InstallFactory(sandbox)
	sandbox.Core.Uninstall = uninstall.UninstallFactory(sandbox)
	sandbox.Core.List = list.ListFactory(sandbox)
	sandbox.Core.ExtensionHelp = extensionhelp.ExtensionHelpFactory(sandbox)

	return sandbox
}
