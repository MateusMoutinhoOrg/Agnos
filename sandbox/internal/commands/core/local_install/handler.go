package local_install

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/rundeps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	// The binary is named after the project at --path, never after the cli
	// running this command: sandbox.Config is agnos's own config, and naming the
	// output from it installs every project over agnos itself.
	project_conf, err := utils.LoadProjectConf(sandbox, stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName))
	if err != nil {
		response.Eprintf("%s\n", err.Error())
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", "")
	}
	if err := utils.ValidateProjectName(sandbox, project_conf.ProjectName); err != nil {
		response.Eprintf("cannot install: %s\n", err.Error())
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", "")
	}

	// Progress goes to Logf: a Printf answers the command line with ExitOk, and
	// a failure after it could no longer change the exit code.
	response.Logf("Building project...\n")
	if err := buildAction.Build(sandbox, api.BuildProps{Path: props.Path, Runtime: "go"}); err != nil {
		response.Eprintf("build failed: %s\n", err.Error())
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", "")
	}

	response.Logf("Installing locally...\n")

	// Get GOEXE
	result, err := sandbox.Deps.RunDeps.Run(rundeps.RunProps{
		Dir:     props.Path,
		Program: "go",
		Args:    []string{"env", "GOEXE"},
	})
	if err != nil {
		response.Eprintf("failed to run go env GOEXE: %s\n", err.Error())
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", "")
	}
	goexe := sandbox.Deps.StringsDeps.TrimSpace(result.Output)

	binName := sandbox.Deps.StringsDeps.ToLower(project_conf.ProjectName) + goexe

	var outPath string
	if sandbox.Deps.StdDeps.GOOS() == "windows" {
		home, err := sandbox.Deps.IoDeps.UserHomeDir()
		if err != nil {
			response.Eprintf("failed to get user home dir: %s\n", err.Error())
			return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", "")
		}
		outPath = sandbox.Deps.IoDeps.Join(home, ".local", "bin", binName)
	} else {
		outPath = sandbox.Deps.IoDeps.Join("/usr/local/bin", binName)
	}

	sandbox.Deps.IoDeps.CreateDir(sandbox.Deps.IoDeps.Dir(outPath))

	sandbox.Deps.StdDeps.Logf("building to %s\n", outPath)
	result, err = sandbox.Deps.RunDeps.Run(rundeps.RunProps{
		Dir:     props.Path,
		Program: "go",
		Args:    []string{"build", "-o", outPath, "./cmd/main"},
	})
	if err != nil {
		response.Eprintf("failed to run go build: %s\n", err.Error())
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", "")
	}
	if result.ExitCode != 0 {
		response.Eprintf("go build failed: %s\n", result.Output)
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", "")
	}

	response.Printf("Installed successfully at %s\n", outPath)
	response.SetStatus(api.ExitOk)
	return nil
}
