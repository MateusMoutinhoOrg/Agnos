package start

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib"
)

func NewCommand(sandbox *lib.SandBox) lib.CliCommand {
	return lib.CliCommand{

		ValidStartIdentifiers: []string{"start"},
		Category:              "Core Commands",

		Args: []lib.CliArg{
			lib.CliArg{
				Id:          "path",
				Description: "the dir to start the project",
				Examples: []string{
					sandbox.ProjectName + " start . ",
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
					sandbox.ProjectName + " start -q",
				},
				Type:             lib.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
			lib.Cliflag{
				Id:               "force",
				ValidIdentifiers: []string{"--force", "-f"},
				Description:      "Forces the creation of the project, overwriting existing files",
				Examples: []string{
					sandbox.ProjectName + " start -f",
				},
				Type:             lib.CliTypeBool,
				RequiredMinSize:  0,
				RequiredMaxSize:  0,
				RequiredPresence: false,
			},
			lib.Cliflag{
				Id:               "module",
				ValidIdentifiers: []string{"--module", "-m"},
				Description:      "Module name for go.mod",
				Examples: []string{
					sandbox.ProjectName + " start -m github.com/user/project",
				},
				Type:             lib.CliTypeString,
				RequiredMinSize:  1,
				RequiredMaxSize:  1,
				RequiredPresence: false,
			},
		},

		Description:     "Initialize a new project in a directory",
		LongDescription: "Scaffolds a new Agnos project in the given directory, creating\nthe required configuration files and folder structure. If no\npath is provided, the current directory is used.",
		Examples: []string{
			sandbox.ProjectName + " start",
			sandbox.ProjectName + " start .",
			sandbox.ProjectName + " start ./my-project",
			sandbox.ProjectName + " start -q",
		},
		Handler: CommandHandler,
	}
}

func CommandHandler(sandbox *lib.SandBox, entries lib.CliEntrys) int {

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

	start_error := sandbox.Start(lib.StartProps{
		Path:        path,
		ProjectName: sandbox.ProjectName,
		Module:      module,
		Force:       forceFlag.Exist,
	})

	if !quietFlag.Exist && start_error != nil {
		sandbox.Deps.Error(start_error.Error())
	}
	if start_error != nil {
		return lib.ExitFailure
	}
	return lib.ExitOk
}
