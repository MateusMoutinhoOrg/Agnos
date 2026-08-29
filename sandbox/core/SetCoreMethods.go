package core

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core/build"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core/extensionhelp"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core/install"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core/list"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core/start"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core/uninstall"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/sandbox"
)

func SetCoreMethods(sandbox *sandbox.SandBox) {
	sandbox.Core.Start = start.StartFactory(sandbox)
	sandbox.Core.Build = build.BuildFactory(sandbox)
	sandbox.Core.Install = install.InstallFactory(sandbox)
	sandbox.Core.Uninstall = uninstall.UninstallFactory(sandbox)
	sandbox.Core.List = list.ListFactory(sandbox)
	sandbox.Core.ExtensionHelp = extensionhelp.ExtensionHelpFactory(sandbox)
}
