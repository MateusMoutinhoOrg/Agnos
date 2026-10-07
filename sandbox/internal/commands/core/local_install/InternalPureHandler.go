package local_install

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/rundeps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	// The binary is named after the project at --path, never after the cli
	// running this command: sandbox.Config is agnos's own config, and naming the
	// output from it installs every project over agnos itself.
	project_conf, err := utils.LoadProjectConf(sandbox, smartio.New(sandbox, props.Path, sandbox.Config.ProjectName))
	if err != nil {
		response.Error("%s\n", err.Error())
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", "")
	}
	if err := utils.ValidateProjectName(sandbox, project_conf.Name); err != nil {
		response.Error("cannot install: %s\n", err.Error())
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", "")
	}

	// Progress goes to Log: a Printf answers the command line with ExitOk, and
	// a failure after it could no longer change the exit code.
	response.Log("Building project...\n")
	if err := buildAction.Build(sandbox, api.BuildProps{Path: props.Path, Runtime: "go"}); err != nil {
		response.Error("build failed: %s\n", err.Error())
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", "")
	}

	response.Log("Installing locally...\n")

	// Get GOEXE
	result, err := sandbox.Deps.Rundeps.Run(rundeps.RunProps{
		Dir:     props.Path,
		Program: "go",
		Args:    []string{"env", "GOEXE"},
	})
	if err != nil {
		response.Error("failed to run go env GOEXE: %s\n", err.Error())
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", "")
	}
	goexe := sandbox.Deps.Stringsdeps.TrimSpace(result.Output)

	binName := sandbox.Deps.Stringsdeps.ToLower(project_conf.Name) + goexe

	var outPath string
	if sandbox.Deps.Std.Goos() == "windows" {
		home, err := sandbox.Deps.Iodeps.UserHomeDir()
		if err != nil {
			response.Error("failed to get user home dir: %s\n", err.Error())
			return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", "")
		}
		outPath = sandbox.Deps.Iodeps.Join(home, ".local", "bin", binName)
	} else {
		outPath = sandbox.Deps.Iodeps.Join("/usr/local/bin", binName)
	}

	sandbox.Deps.Iodeps.CreateDir(sandbox.Deps.Iodeps.Dir(outPath))

	sandbox.Deps.Std.Log("building to %s\n", outPath)
	result, err = sandbox.Deps.Rundeps.Run(rundeps.RunProps{
		Dir:     props.Path,
		Program: "go",
		Args:    []string{"build", "-o", outPath, "./cmd/main"},
	})
	if err != nil {
		response.Error("failed to run go build: %s\n", err.Error())
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", "")
	}
	if result.ExitCode != 0 {
		response.Error("go build failed: %s\n", result.Output)
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", "")
	}

	response.Printf("Installed successfully at %s\n", outPath)
	response.SetStatus(api.ExitOk)
	return nil
}
