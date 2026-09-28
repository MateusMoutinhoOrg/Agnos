package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// Build re-renders every generated file of the project at props.Path and then
// hands the result to props.Runtime, so a build only reports success when the
// toolchain accepts what was rendered.
func Build(sandbox *api.Sandbox, props api.BuildProps) error {
	io := smartio.New(sandbox, props.Path, sandbox.Config.ProjectName)
	err := BuildInternal(sandbox, io, props.Path)
	if err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return RunRuntime(sandbox, props.Path, props.Runtime)
}

// PersistAndBuild is the follow-up every editing action ends on: it persists
// io, then builds the project at props.Path. When the render refuses what was
// persisted — a declaration no generated file can be spelled from — what io
// persisted is undone before the error is returned, so a command that fails
// never leaves its change behind. A failure of the runtime itself keeps the
// change: the render is consistent, and the compile error may be the
// project's own code.
func PersistAndBuild(sandbox *api.Sandbox, io *smartio.SmartIO, props api.BuildProps) error {
	if err := io.Persist(); err != nil {
		return err
	}
	rendered := smartio.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := BuildInternal(sandbox, rendered, props.Path); err != nil {
		if undo_err := io.Undo(); undo_err != nil {
			return sandbox.Deps.Std.Errorf("%w (and the change could not be undone: %s)", err, undo_err.Error())
		}
		return sandbox.Deps.Std.Errorf("%w (nothing was changed)", err)
	}
	if err := rendered.Persist(); err != nil {
		return err
	}
	return RunRuntime(sandbox, props.Path, props.Runtime)
}
