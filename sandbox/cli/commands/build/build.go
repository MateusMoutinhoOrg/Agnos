package build

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib"
)

func NewCommand(sandbox *lib.SandBox) lib.CliCommand {
	return lib.CliCommand{

		ValidStartIdentifiers: []string{"build"},
		Category:              "Core Commands",

		Args: []lib.CliArg{
			lib.CliArg{
				Id:          "path",
				Description: "the dir to build the project",
				Examples: []string{
					sandbox.ProjectName + " build . ",
				},
				Defaults:        []string{"."},
				RequiredType:    lib.CliTypeString,
				RequiredMinSize: 0,
				RequiredMaxSize: 1,
			},
		},
		Flags: []lib.Cliflag{
			lib.Cliflag{
				Id:               "quiet",
				ValidIdentifiers: []string{"--quiet", "-q"},
				Description:      "Quiets the cli output",
				Examples: []string{
					sandbox.ProjectName + " build -q",
				},
				Type:             lib.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
		},

		Description:     "Build the project in a directory",
		LongDescription: "Builds the project in the given directory, compiling\nthe source code into the output artifacts. If no\npath is provided, the current directory is used.",
		Examples: []string{
			sandbox.ProjectName + " build",
			sandbox.ProjectName + " build .",
			sandbox.ProjectName + " build ./my-project",
			sandbox.ProjectName + " build -q",
		},
		Handler: CommandHandler,
	}
}

func CommandHandler(sandbox *lib.SandBox, entries lib.CliEntrys) int {

	quietFlag := entries.GetFlagById("quiet")
	pathArg := entries.GetArgById("path")
	path := pathArg.Values[0].String()

	build_error := sandbox.Build(lib.BuildProps{
		Path:    path,
		Project: sandbox.ProjectName,
	})

	if !quietFlag.Exist && build_error != nil {
		sandbox.Deps.Error(build_error.Error())
	}
	if build_error != nil {
		return lib.ExitFailure
	}
	return lib.ExitOk
}
