package publish

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/rundeps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	compileAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/compile"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	if entries.Publisher != "gh" {
		deps.Std.Error("Unsupported publisher %q. The only available publisher is \"gh\".\n", entries.Publisher)
		return api.ExitFailure
	}

	io := smartio.New(deps, entries.Path, config.ProjectName)

	releaseName := entries.ReleaseName
	if releaseName == "" {
		rel := config.ProjectName + "Config/project.yaml"
		content, err := io.ReadFile(rel)
		if err != nil {
			deps.Std.Error("could not read %s to determine release name: %s\n", rel, err.Error())
			return api.ExitFailure
		}

		re := regexp.MustCompile(`(?m)^version:\s*([^\s]+)`)
		match := re.FindStringSubmatch(string(content))
		if len(match) > 1 {
			releaseName = match[1]
		} else {
			deps.Std.Error("could not find version in %s and no --release-name provided\n", rel)
			return api.ExitFailure
		}
	}

	deps.Std.Printf("Building project...\n")
	if err := buildAction.Build(deps, api.BuildProps{Path: entries.Path, Runtime: "go"}); err != nil {
		deps.Std.Error("build failed: %s\n", err.Error())
		return api.ExitFailure
	}

	deps.Std.Printf("Compiling targets...\n")
	// If targets is empty or "all", we pass "all". 
	targets := entries.Target
	if targets == "" {
		targets = "all"
	}
	if err := compileAction.Compile(deps, api.CompileProps{Path: entries.Path, Targets: []string{targets}}); err != nil {
		deps.Std.Error("compile failed: %s\n", err.Error())
		return api.ExitFailure
	}

	deps.Std.Printf("Gathering compiled binaries...\n")
	releaseDir := filepath.Join(entries.Path, "release")
	entriesInfo, err := os.ReadDir(releaseDir)
	if err != nil {
		deps.Std.Error("could not read release directory: %s\n", err.Error())
		return api.ExitFailure
	}

	args := []string{"release", "create", releaseName}
	if entries.Draft {
		args = append(args, "--draft")
	}

	args = append(args, "--title", fmt.Sprintf("Release %s", releaseName))
	
	// Default notes can be added, or leave to gh defaults

	for _, e := range entriesInfo {
		if !e.IsDir() {
			args = append(args, filepath.Join(releaseDir, e.Name()))
		}
	}

	deps.Std.Printf("Creating release %s with gh...\n", releaseName)
	result, err := deps.Rundeps.Run(rundeps.RunProps{
		Dir:     entries.Path,
		Program: "gh",
		Args:    args,
	})

	if err != nil {
		deps.Std.Error("gh execution failed: %s\n", err.Error())
		return api.ExitFailure
	}
	if result.ExitCode != 0 {
		deps.Std.Error("gh release create failed: %s\n", result.Output)
		return api.ExitFailure
	}

	deps.Std.Printf("Release %s created and published successfully!\n", releaseName)
	return api.ExitOk
}
