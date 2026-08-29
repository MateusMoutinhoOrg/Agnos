package sandbox

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/core"
	sandbox "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/sandbox"

	// CORE
	start "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core/start"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core/build"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core/extensionhelp"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core/install"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core/list"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core/uninstall"

	// CMD RELATED
	cliMain "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/cli"
	cmdBuild "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/cli/commands/build"
	cmdExtensionhelp "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/cli/commands/extensionhelp"
	cmdHelp "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/cli/commands/help"
	cmdInstall "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/cli/commands/install"
	cmdList "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/cli/commands/list"
	cmdStart "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/cli/commands/start"
	cmdUninstall "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/cli/commands/uninstall"
	cmdVersion "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/cli/commands/version"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
	cliApi "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/cli"
)

func New(d deps.Deps) sandbox.SandBox {
	sandbox := sandbox.SandBox{Deps: d}

	sandbox.Config.Version = "v0.0.12"
	sandbox.Config.ProjectName = "Agnos"

	sandbox.Cli.Commands = []cliApi.CliCommand{
		cmdVersion.NewCommand(&sandbox),
		cmdHelp.NewCommand(&sandbox),
		cmdStart.NewCommand(&sandbox),
		cmdBuild.NewCommand(&sandbox),
		cmdInstall.NewCommand(&sandbox),
		cmdUninstall.NewCommand(&sandbox),
		cmdList.NewCommand(&sandbox),
		cmdExtensionhelp.NewCommand(&sandbox),
	}

	sandbox.Cli.CliMain = cliMain.CliMainFactory(&sandbox)

	sandbox.Core.Start = func(props core.StartProps) error {
		return start.Start(&sandbox, props)
	}
	sandbox.Core.Build = func(props core.BuildProps) error {
		return build.Build(&sandbox, props)
	}
	sandbox.Core.Install = func(props core.InstallProps) error {
		return install.Install(&sandbox, props)
	}
	sandbox.Core.Uninstall = func(props core.UninstallProps) error {
		return uninstall.Uninstall(&sandbox, props)
	}
	sandbox.Core.List = func(props core.ListProps) error {
		return list.List(&sandbox, props)
	}
	sandbox.Core.ExtensionHelp = func(props core.ExtensionHelpProps) error {
		return extensionhelp.ExtensionHelp(&sandbox, props)
	}

	return sandbox
}
