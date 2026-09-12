package local_install

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/rundeps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	sandbox.Deps.Std.Printf("Building project...\n")
	if err := buildAction.Build(sandbox, api.BuildProps{Path: entries.Path, Runtime: "go"}); err != nil {
		sandbox.Deps.Std.Error("build failed: %s\n", err.Error())
		return api.ExitFailure
	}

	sandbox.Deps.Std.Printf("Installing locally...\n")

	// Get GOEXE
	result, err := sandbox.Deps.Rundeps.Run(rundeps.RunProps{
		Dir:     entries.Path,
		Program: "go",
		Args:    []string{"env", "GOEXE"},
	})
	if err != nil {
		sandbox.Deps.Std.Error("failed to run go env GOEXE: %s\n", err.Error())
		return api.ExitFailure
	}
	goexe := sandbox.Deps.Stringsdeps.TrimSpace(result.Output)

	binName := sandbox.Deps.Stringsdeps.ToLower(config.ProjectName) + goexe

	var outPath string
	if sandbox.Deps.Std.Goos() == "windows" {
		home, err := sandbox.Deps.Iodeps.UserHomeDir()
		if err != nil {
			sandbox.Deps.Std.Error("failed to get user home dir: %s\n", err.Error())
			return api.ExitFailure
		}
		outPath = sandbox.Deps.Iodeps.Join(home, ".local", "bin", binName)
	} else {
		outPath = sandbox.Deps.Iodeps.Join("/usr/local/bin", binName)
	}

	sandbox.Deps.Iodeps.CreateDir(sandbox.Deps.Iodeps.Dir(outPath))

	sandbox.Deps.Std.Log("building to %s\n", outPath)
	result, err = sandbox.Deps.Rundeps.Run(rundeps.RunProps{
		Dir:     entries.Path,
		Program: "go",
		Args:    []string{"build", "-o", outPath, "./cmd/main"},
	})
	if err != nil {
		sandbox.Deps.Std.Error("failed to run go build: %s\n", err.Error())
		return api.ExitFailure
	}
	if result.ExitCode != 0 {
		sandbox.Deps.Std.Error("go build failed: %s\n", result.Output)
		return api.ExitFailure
	}

	sandbox.Deps.Std.Printf("Installed successfully at %s\n", outPath)
	return api.ExitOk
}
