package publish

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/rundeps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	compileAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/compile"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	if entries.Publisher != "gh" {
		response.Error("Unsupported publisher %q. The only available publisher is \"gh\".\n", entries.Publisher)
		return cliio.Fail(sandbox, api.ExitFailure, "", "")
	}

	io := smartio.New(sandbox, props.Path, sandbox.Config.ProjectName)

	// The release defaults to the version of the project at --path. A project
	// that declares none (`version: null`, what `start` writes) has no release
	// name to offer, so it is a refusal rather than a release called "null".
	releaseName := sandbox.Deps.Stringsdeps.TrimSpace(entries.ReleaseName)
	if releaseName == "" {
		project_conf, err := utils.LoadProjectConf(sandbox, io)
		if err != nil {
			response.Error("%s\n", err.Error())
			return cliio.Fail(sandbox, api.ExitFailure, "", "")
		}
		releaseName = project_conf.Version
		if releaseName == "" {
			response.Error("%s declares no version: set one there or pass --release-name\n", utils.ProjectConfPath(sandbox))
			return cliio.Fail(sandbox, api.ExitFailure, "", "")
		}
	}

	targets := publishTargets(sandbox, entries.Target)
	outputs, err := compileAction.Outputs(sandbox, targets)
	if err != nil {
		response.Error("%s\n", err.Error())
		return cliio.Fail(sandbox, api.ExitFailure, "", "")
	}

	// Progress goes to Log: a Printf answers the command line with ExitOk, and
	// a failure after it could no longer change the exit code.
	response.Log("Building project...\n")
	if err := buildAction.Build(sandbox, api.BuildProps{Path: props.Path, Runtime: "go"}); err != nil {
		response.Error("build failed: %s\n", err.Error())
		return cliio.Fail(sandbox, api.ExitFailure, "", "")
	}

	response.Log("Compiling targets...\n")
	if err := compileAction.Compile(sandbox, api.CompileProps{Path: props.Path, Targets: targets}); err != nil {
		response.Error("compile failed: %s\n", err.Error())
		return cliio.Fail(sandbox, api.ExitFailure, "", "")
	}

	args := []string{"release", "create", releaseName}
	if entries.Draft {
		args = append(args, "--draft")
	}

	args = append(args, "--title", sandbox.Deps.Std.Sprintf("Release %s", releaseName))

	// Only the binaries this run compiled are uploaded: release/ may still hold
	// the outputs of an earlier compile for other targets, stale ones included.
	args = append(args, outputs...)

	response.Log("Creating release %s with gh...\n", releaseName)
	result, err := sandbox.Deps.Rundeps.Run(rundeps.RunProps{
		Dir:     props.Path,
		Program: "gh",
		Args:    args,
	})

	if err != nil {
		response.Error("gh execution failed: %s\n", err.Error())
		return cliio.Fail(sandbox, api.ExitFailure, "", "")
	}
	if result.ExitCode != 0 {
		response.Error("gh release create failed: %s\n", result.Output)
		return cliio.Fail(sandbox, api.ExitFailure, "", "")
	}

	response.Printf("Release %s created and published successfully!\n", releaseName)
	response.SetStatus(api.ExitOk)
	return nil
}

// publishTargets is the --target list compile is handed, as compile reads
// its own: every target when none is given.
func publishTargets(sandbox *api.Sandbox, raws []string) []string {
	targets := []string{}
	for _, raw := range raws {
		if raw = sandbox.Deps.Stringsdeps.TrimSpace(raw); raw != "" {
			targets = append(targets, raw)
		}
	}
	if len(targets) == 0 {
		return []string{"all"}
	}
	return targets
}
