// This file is an illustrative sample, not part of the build. It shows the two
// files of one action in a single listing; in the tree they are separate.

// ─── sandbox/internal/actions/dep_remove/dep_remove.go ──────────────────────

package dep_remove

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

// DepRemove deletes what dep installed, persists, and rebuilds with the none
// runtime — dropping a contract may leave sandbox code referring to it.
func DepRemove(deps *deps.Deps, path string, dep string) error {
	io := smartio.New(deps, path, config.ProjectName)
	if err := DepRemoveInternal(deps, io, path, dep); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(deps, api.BuildProps{Path: path, Runtime: api.RuntimeNone})
}

// ─── sandbox/internal/actions/dep_remove/dep_remove_internal.go ─────────────

// DepRemoveInternal removes, inside the open transaction, every file that
// assets/deplist/<dep> installs, plus any directory the removal emptied.
func DepRemoveInternal(deps *deps.Deps, io *smartio.SmartIO, path string, dep string) error {
	deps.Std.Log("dep-remove started with path %s dep %s \n", path, dep)

	files, err := deps.Embeddeps.ListFilesRecursively("deplist/" + dep)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return deps.Std.Errorf("unknown dep %q", dep)
	}
	for _, file := range files {
		io.RemoveDir(file) // project-relative: the path the asset holds inside the dep
	}
	return nil
}
