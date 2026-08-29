package core

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core/build"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core/extensionhelp"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core/install"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core/list"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core/start"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core/uninstall"
)

func SetCoreMethods(sandbox *lib.SandBox) {
	sandbox.Start = start.StartFactory(sandbox)
	sandbox.Build = build.BuildFactory(sandbox)
	sandbox.Install = install.InstallFactory(sandbox)
	sandbox.Uninstall = uninstall.UninstallFactory(sandbox)
	sandbox.List = list.ListFactory(sandbox)
	sandbox.ExtensionHelp = extensionhelp.ExtensionHelpFactory(sandbox)
}
