package compile

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
)

// Compile runs `build` over the project at props.Path and then cross-compiles
// its ./cmd/main entrypoint once per requested target into release/. The
// target names are resolved before the build runs, so an unknown target fails
// fast; the build runs before any binary is produced, so every release comes
// from a freshly rendered, compilable tree.
func Compile(sandbox *api.Sandbox, props api.CompileProps) error {
	names, err := resolveTargets(sandbox, props.Targets)
	if err != nil {
		return err
	}

	sandbox.Deps.Std.Log("compile started with path %s \n", props.Path)

	if err := buildAction.Build(sandbox, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo}); err != nil {
		return err
	}

	return CompileInternal(sandbox, props.Path, names)
}
