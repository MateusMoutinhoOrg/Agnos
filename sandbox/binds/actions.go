package binds

import (
	api "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	deps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/build"
	startAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/start"
	enableDepsAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/enable_deps"
	removeDepsAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/remove_deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

func ActionsBind(deps *deps.Deps, sandbox *api.Sandbox) {
	sandbox.Actions.Build = func(path string) error {
		io := smartio.New(deps, path, config.ProjectName)
		err := buildAction.Build(deps, io, path)
		if err != nil {
			return err
		}
		return io.Persist()
	}
	sandbox.Actions.Start = func(props api.StartProps) error {
		io := smartio.New(deps, props.Path, props.ProjectName)
		err := startAction.Start(deps, io, props)
		if err == nil {
			err = io.Persist()
		}
		if err == nil {
			err = buildAction.Build(deps, io, props.Path)
		}
		if err != nil {
			return err
		}
		return io.Persist()
	}
	sandbox.Actions.EnableDeps = func(path string) error {
		io := smartio.New(deps, path, config.ProjectName)
		err := enableDepsAction.EnableDeps(deps, io, path)
		if err == nil {
			err = io.Persist()
		}
		if err == nil {
			err = buildAction.Build(deps, io, path)
		}
		if err != nil {
			return err
		}
		return io.Persist()
	}
	sandbox.Actions.RemoveDeps = func(path string) error {
		io := smartio.New(deps, path, config.ProjectName)
		err := removeDepsAction.RemoveDeps(deps, io, path)
		if err == nil {
			err = io.Persist()
		}
		if err == nil {
			err = buildAction.Build(deps, io, path)
		}
		if err != nil {
			return err
		}
		return io.Persist()
	}
}
