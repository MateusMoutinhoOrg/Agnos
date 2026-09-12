package publish

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/rundeps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	compileAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/compile"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	if entries.Publisher != "gh" {
		sandbox.Deps.Std.Error("Unsupported publisher %q. The only available publisher is \"gh\".\n", entries.Publisher)
		return api.ExitFailure
	}

	io := smartio.New(sandbox, entries.Path, config.ProjectName)

	releaseName := entries.ReleaseName
	if releaseName == "" {
		rel := config.ProjectName + "Config/project.yaml"
		content, err := io.ReadFile(rel)
		if err != nil {
			sandbox.Deps.Std.Error("could not read %s to determine release name: %s\n", rel, err.Error())
			return api.ExitFailure
		}

		releaseName = versionOf(sandbox, string(content))
		if releaseName == "" {
			sandbox.Deps.Std.Error("could not find version in %s and no --release-name provided\n", rel)
			return api.ExitFailure
		}
	}

	sandbox.Deps.Std.Printf("Building project...\n")
	if err := buildAction.Build(sandbox, api.BuildProps{Path: entries.Path, Runtime: "go"}); err != nil {
		sandbox.Deps.Std.Error("build failed: %s\n", err.Error())
		return api.ExitFailure
	}

	sandbox.Deps.Std.Printf("Compiling targets...\n")
	// If targets is empty or "all", we pass "all".
	targets := entries.Target
	if targets == "" {
		targets = "all"
	}
	if err := compileAction.Compile(sandbox, api.CompileProps{Path: entries.Path, Targets: []string{targets}}); err != nil {
		sandbox.Deps.Std.Error("compile failed: %s\n", err.Error())
		return api.ExitFailure
	}

	sandbox.Deps.Std.Printf("Gathering compiled binaries...\n")
	releaseDir := sandbox.Deps.Iodeps.Join(entries.Path, "release")
	if !sandbox.Deps.Iodeps.IsDir(releaseDir) {
		sandbox.Deps.Std.Error("could not read release directory: %s\n", releaseDir)
		return api.ExitFailure
	}

	args := []string{"release", "create", releaseName}
	if entries.Draft {
		args = append(args, "--draft")
	}

	args = append(args, "--title", sandbox.Deps.Std.Sprintf("Release %s", releaseName))

	// Default notes can be added, or leave to gh defaults

	args = append(args, sandbox.Deps.Iodeps.ListFiles(releaseDir)...)

	sandbox.Deps.Std.Printf("Creating release %s with gh...\n", releaseName)
	result, err := sandbox.Deps.Rundeps.Run(rundeps.RunProps{
		Dir:     entries.Path,
		Program: "gh",
		Args:    args,
	})

	if err != nil {
		sandbox.Deps.Std.Error("gh execution failed: %s\n", err.Error())
		return api.ExitFailure
	}
	if result.ExitCode != 0 {
		sandbox.Deps.Std.Error("gh release create failed: %s\n", result.Output)
		return api.ExitFailure
	}

	sandbox.Deps.Std.Printf("Release %s created and published successfully!\n", releaseName)
	return api.ExitOk
}

// versionOf reads the `version:` field out of a project.yaml, returning "" when
// the file declares none. The whole file is scanned line by line rather than
// parsed: publish runs before the config is otherwise needed, and the one field
// it wants is a plain `key: value` at the start of a line.
func versionOf(sandbox *api.Sandbox, content string) string {
	for _, line := range sandbox.Deps.Stringsdeps.Split(content, "\n") {
		if !sandbox.Deps.Stringsdeps.HasPrefix(line, "version:") {
			continue
		}
		fields := sandbox.Deps.Stringsdeps.Fields(sandbox.Deps.Stringsdeps.TrimPrefix(line, "version:"))
		if len(fields) > 0 {
			return fields[0]
		}
	}
	return ""
}
