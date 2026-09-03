package local_install

import (
	"path/filepath"
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/rundeps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	deps.Std.Printf("Building project...\n")
	if err := buildAction.Build(deps, api.BuildProps{Path: entries.Path, Runtime: "go"}); err != nil {
		deps.Std.Error("build failed: %s\n", err.Error())
		return api.ExitFailure
	}

	deps.Std.Printf("Installing locally...\n")
	
	// Get GOBIN
	result, err := deps.Rundeps.Run(rundeps.RunProps{
		Dir:     entries.Path,
		Program: "go",
		Args:    []string{"env", "GOBIN"},
	})
	if err != nil {
		deps.Std.Error("failed to run go env GOBIN: %s\n", err.Error())
		return api.ExitFailure
	}
	gobin := strings.TrimSpace(result.Output)

	if gobin == "" {
		// Fallback to GOPATH/bin
		result, err = deps.Rundeps.Run(rundeps.RunProps{
			Dir:     entries.Path,
			Program: "go",
			Args:    []string{"env", "GOPATH"},
		})
		if err != nil {
			deps.Std.Error("failed to run go env GOPATH: %s\n", err.Error())
			return api.ExitFailure
		}
		gopath := strings.TrimSpace(result.Output)
		if gopath == "" {
			deps.Std.Error("GOPATH is empty, cannot determine install location\n")
			return api.ExitFailure
		}
		gobin = filepath.Join(gopath, "bin")
	}

	// Get GOEXE
	result, err = deps.Rundeps.Run(rundeps.RunProps{
		Dir:     entries.Path,
		Program: "go",
		Args:    []string{"env", "GOEXE"},
	})
	if err != nil {
		deps.Std.Error("failed to run go env GOEXE: %s\n", err.Error())
		return api.ExitFailure
	}
	goexe := strings.TrimSpace(result.Output)

	binName := strings.ToLower(config.ProjectName) + goexe
	outPath := filepath.Join(gobin, binName)

	deps.Std.Log("building to %s\n", outPath)
	result, err = deps.Rundeps.Run(rundeps.RunProps{
		Dir:     entries.Path,
		Program: "go",
		Args:    []string{"build", "-o", outPath, "./cmd/main"},
	})
	if err != nil {
		deps.Std.Error("failed to run go build: %s\n", err.Error())
		return api.ExitFailure
	}
	if result.ExitCode != 0 {
		deps.Std.Error("go build failed: %s\n", result.Output)
		return api.ExitFailure
	}

	deps.Std.Printf("Installed successfully at %s\n", outPath)
	return api.ExitOk
}
