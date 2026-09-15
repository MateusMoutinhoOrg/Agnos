package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/rundeps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// runtimeStep is one program the runtime hands the rendered project to, in
// the order they run.
type runtimeStep struct {
	Program string
	Args    []string
}

// compilableDirs are the top-level directories of the harness schema that
// hold real Go packages. They are what the go runtime compiles — never
// `./...`, because assets/ holds Go *templates* ({{.Module}} import paths),
// which are not parsable Go and never meant to be.
var compilableDirs = []string{"cmd", "sandbox", "adapters"}

// RunRuntime hands the rendered project to a real toolchain and reports
// whether that toolchain accepted it. Rendering templates only proves the
// files were written; without this step `build` and `verify` would report
// success over a project that does not compile.
//
// It runs after the transaction is persisted — the toolchain reads the disk,
// not the pending writes.
func RunRuntime(sandbox *api.Sandbox, path string, runtime string) error {
	steps, err := runtimeSteps(sandbox, path, runtime)
	if err != nil {
		return err
	}

	for _, step := range steps {
		command := step.Program + " " + sandbox.Deps.Stringsdeps.Join(step.Args, " ")
		sandbox.Deps.Std.Log("runtime %s: %s \n", runtime, command)

		result, err := sandbox.Deps.Rundeps.Run(rundeps.RunProps{
			Dir:     path,
			Program: step.Program,
			Args:    step.Args,
		})
		if err != nil {
			return sandbox.Deps.Std.Errorf("runtime %s: could not run `%s`: %w", runtime, command, err)
		}
		if result.ExitCode != 0 {
			return sandbox.Deps.Std.Errorf("runtime %s: `%s` failed:\n%s", runtime, command, result.Output)
		}
	}

	return nil
}

// runtimeSteps maps a runtime name onto the steps it runs. An unknown name is
// a usage error rather than a silent skip.
func runtimeSteps(sandbox *api.Sandbox, path string, runtime string) ([]runtimeStep, error) {
	switch runtime {
	case api.RuntimeNone, "":
		return nil, nil
	case api.RuntimeGo:
		return goRuntimeSteps(sandbox, path), nil
	default:
		return nil, sandbox.Deps.Std.Errorf("unknown runtime %q (use %q or %q)", runtime, api.RuntimeGo, api.RuntimeNone)
	}
}

// goRuntimeSteps resolve the module graph, then compile every package of the
// schema that the project actually has. `go mod tidy` is the step that writes
// go.sum, so a freshly scaffolded project is left in a state `go build`
// accepts; the compile step is skipped when nothing is there to compile yet.
func goRuntimeSteps(sandbox *api.Sandbox, path string) []runtimeStep {
	steps := []runtimeStep{
		{Program: "go", Args: []string{"mod", "tidy"}},
	}

	io := smartio.New(sandbox, path, config.ProjectName)
	packages := []string{"build"}
	for _, dir := range compilableDirs {
		if io.IsDir(dir) {
			packages = append(packages, "./"+dir+"/...")
		}
	}
	if len(packages) > 1 {
		steps = append(steps, runtimeStep{Program: "go", Args: packages})
	}

	return steps
}
