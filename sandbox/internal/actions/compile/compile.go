package compile

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
)

// Compile runs `build` over the project at props.Path and then cross-compiles
// its ./cmd/main entrypoint once per requested target into release/. The
// target names are resolved before the build runs, so an unknown target fails
// fast; the build runs before any binary is produced, so every release comes
// from a freshly rendered, compilable tree.
func Compile(deps *deps.Deps, props api.CompileProps) error {
	names, err := resolveTargets(deps, props.Targets)
	if err != nil {
		return err
	}

	deps.Std.Log("compile started with path %s \n", props.Path)

	if err := buildAction.Build(deps, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo}); err != nil {
		return err
	}

	return CompileInternal(deps, props.Path, names)
}
