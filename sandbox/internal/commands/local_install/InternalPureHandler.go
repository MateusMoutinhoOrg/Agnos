package local_install

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/rundeps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	response.Printf("Building project...\n")
	if err := buildAction.Build(sandbox, api.BuildProps{Path: props.Path, Runtime: "go"}); err != nil {
		response.Error("build failed: %s\n", err.Error())
		return cliio.Fail(sandbox, api.ExitFailure, "", "")
	}

	response.Printf("Installing locally...\n")

	// Get GOEXE
	result, err := sandbox.Deps.Rundeps.Run(rundeps.RunProps{
		Dir:     props.Path,
		Program: "go",
		Args:    []string{"env", "GOEXE"},
	})
	if err != nil {
		response.Error("failed to run go env GOEXE: %s\n", err.Error())
		return cliio.Fail(sandbox, api.ExitFailure, "", "")
	}
	goexe := sandbox.Deps.Stringsdeps.TrimSpace(result.Output)

	binName := sandbox.Deps.Stringsdeps.ToLower(sandbox.Config.ProjectName) + goexe

	var outPath string
	if sandbox.Deps.Std.Goos() == "windows" {
		home, err := sandbox.Deps.Iodeps.UserHomeDir()
		if err != nil {
			response.Error("failed to get user home dir: %s\n", err.Error())
			return cliio.Fail(sandbox, api.ExitFailure, "", "")
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
		return cliio.Fail(sandbox, api.ExitFailure, "", "")
	}
	if result.ExitCode != 0 {
		response.Error("go build failed: %s\n", result.Output)
		return cliio.Fail(sandbox, api.ExitFailure, "", "")
	}

	response.Printf("Installed successfully at %s\n", outPath)
	response.SetStatus(api.ExitOk)
	return nil
}
