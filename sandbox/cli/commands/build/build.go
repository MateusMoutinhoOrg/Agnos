package build

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/cli"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/core"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/sandbox"
)

func NewCommand(sandbox *sandbox.SandBox) cli.CliCommand {
	return cli.CliCommand{

		ValidStartIdentifiers: []string{"build"},
		Category:              "Core Commands",

		Args: []cli.CliArg{
			cli.CliArg{
				Id:          "path",
				Description: "the dir to build the project",
				Examples: []string{
					sandbox.Config.ProjectName + " build . ",
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
					sandbox.Config.ProjectName + " build -q",
				},
				Type:             cli.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
		},

		Description:     "Build the project in a directory",
		LongDescription: "Builds the project in the given directory, compiling\nthe source code into the output artifacts. If no\npath is provided, the current directory is used.",
		Examples: []string{
			sandbox.Config.ProjectName + " build",
			sandbox.Config.ProjectName + " build .",
			sandbox.Config.ProjectName + " build ./my-project",
			sandbox.Config.ProjectName + " build -q",
		},
		Handler: CommandHandler,
	}
}

func CommandHandler(sb any, entries cli.CliEntrys) int {
	sandbox := sb.(*sandbox.SandBox)

	quietFlag := entries.GetFlagById("quiet")
	pathArg := entries.GetArgById("path")
	path := pathArg.Values[0].String()

	build_error := sandbox.Core.Build(core.BuildProps{
		Path:    path,
		Project: sandbox.Config.ProjectName,
	})

	if !quietFlag.Exist && build_error != nil {
		sandbox.Deps.Error(build_error.Error())
	}
	if build_error != nil {
		return cli.ExitFailure
	}
	return cli.ExitOk
}
