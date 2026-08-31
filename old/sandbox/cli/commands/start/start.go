package start

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/cli"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/core"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/sandbox"
)

func NewCommand(sandbox *sandbox.SandBox) cli.CliCommand {
	return cli.CliCommand{

		ValidStartIdentifiers: []string{"start"},
		Category:              "Core Commands",

		Args: []cli.CliArg{
			cli.CliArg{
				Id:          "path",
				Description: "the dir to start the project",
				Examples: []string{
					sandbox.Config.ProjectName + " start . ",
				},
				Defaults:        []string{"."},
				RequiredType:    cli.CliTypeString,
				RequiredMinSize: 0,
				RequiredMaxSize: 1,
			},
		},
		Flags: []cli.Cliflag{
			cli.Cliflag{
				Id:               "quiet",
				ValidIdentifiers: []string{"--quiet", "-q"},
				Description:      "Quiets the cli output",
				Examples: []string{
					sandbox.Config.ProjectName + " start -q",
				},
				Type:             cli.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
			cli.Cliflag{
				Id:               "force",
				ValidIdentifiers: []string{"--force", "-f"},
				Description:      "Forces the creation of the project, overwriting existing files",
				Examples: []string{
					sandbox.Config.ProjectName + " start -f",
				},
				Type:             cli.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
			cli.Cliflag{
				Id:               "module",
				ValidIdentifiers: []string{"--module", "-m"},
				Description:      "Module name for go.mod",
				Examples: []string{
					sandbox.Config.ProjectName + " start -m github.com/user/project",
				},
				Type:             cli.CliTypeString,
				RequiredMinSize:  1,
				RequiredMaxSize:  1,
				RequiredPresence: false,
			},
		},

		Description:     "Initialize a new project in a directory",
		LongDescription: "Scaffolds a new Agnos project in the given directory, creating\nthe required configuration files and folder structure. If no\npath is provided, the current directory is used.",
		Examples: []string{
			sandbox.Config.ProjectName + " start",
			sandbox.Config.ProjectName + " start .",
			sandbox.Config.ProjectName + " start ./my-project",
			sandbox.Config.ProjectName + " start -q",
		},
		Handler: CommandHandler,
	}
}

func CommandHandler(sb any, entries cli.CliEntrys) int {
	sandbox := sb.(*sandbox.SandBox)

	quietFlag := entries.GetFlagById("quiet")
	forceFlag := entries.GetFlagById("force")
	moduleFlag := entries.GetFlagById("module")
	pathArg := entries.GetArgById("path")
	path := pathArg.Values[0].String()

	var module *string
	if moduleFlag.Exist && len(moduleFlag.Values) > 0 {
		modVal := moduleFlag.Values[0].String()
		module = &modVal
	}

	start_error := sandbox.Core.Start(core.StartProps{
		Path:        path,
		ProjectName: sandbox.Config.ProjectName,
		Module:      module,
		Force:       forceFlag.Exist,
	})

	if !quietFlag.Exist && start_error != nil {
		sandbox.Deps.Error(start_error.Error())
	}
	if start_error != nil {
		return cli.ExitFailure
	}
	return cli.ExitOk
}
