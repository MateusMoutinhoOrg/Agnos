package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// legacyPropsDir is where RouteProps and CommandProps were declared before
// they moved under sandbox/internal: beside the contracts, where no field could
// name a type of the project.
const legacyPropsDir = "sandbox/api"

// MigrateLegacyProps moves a props struct an older build wrote into
// sandbox/api/ to the package that declares it now — routeprops for a project
// with the server layer on, commandprops for one with the cli layer on. The
// file is the project's, so it is moved as it is, only its package clause
// renamed; a struct already at the new home is left alone, and so is the old
// file beside it.
//
// It runs before anything reads sandbox/api/, so the moved file is never taken
// for a contract of the Sandbox.
func MigrateLegacyProps(sandbox *api.Sandbox, io *stagedfs.StagedFS, hasCli bool, hasServer bool) error {
	if hasServer {
		if err := migrateLegacyPropsFile(sandbox, io, utils.RoutePropsFile, utils.RoutePropsDir); err != nil {
			return err
		}
	}
	if hasCli {
		if err := migrateLegacyPropsFile(sandbox, io, utils.CommandPropsFile, utils.CommandPropsDir); err != nil {
			return err
		}
	}
	return nil
}

// migrateLegacyPropsFile moves sandbox/api/<file> to <dir>/<file>, its package
// clause renamed to the last segment of dir.
func migrateLegacyPropsFile(sandbox *api.Sandbox, io *stagedfs.StagedFS, file string, dir string) error {
	legacy := legacyPropsDir + "/" + file
	dest := dir + "/" + file
	if !io.IsFile(legacy) || io.IsFile(dest) {
		return nil
	}

	content, err := io.ReadFile(legacy)
	if err != nil {
		return err
	}

	lines := sandbox.Deps.StringsDeps.Split(string(content), "\n")
	for i, line := range lines {
		if sandbox.Deps.StringsDeps.TrimSpace(line) == "package api" {
			lines[i] = "package " + lastSegmentOf(sandbox, dir)
			break
		}
	}

	if err := io.WriteFile(dest, []byte(sandbox.Deps.StringsDeps.Join(lines, "\n"))); err != nil {
		return err
	}
	io.RemoveDir(legacy)
	return nil
}

// migrateHandWrittenProps moves a <File> without the generated marker to
// project.go, `type <Type> struct` renamed `type Project struct` and its doc
// comment with it. A project.go already there is never overwritten: the two
// are the project's, and which one wins is for it to say.
func migrateHandWrittenProps(sandbox *api.Sandbox, io *stagedfs.StagedFS, props propsAggregate, dest string, project string) error {
	content, err := io.ReadFile(dest)
	if err != nil {
		return nil
	}
	text := string(content)
	if sandbox.Deps.StringsDeps.Contains(text, generatedMarker) {
		return nil
	}
	if _, err := io.ReadFile(project); err == nil {
		return sandbox.Deps.StdDeps.Errorf("%s was written by hand and %s exists too: move the fields of %s into %s, then remove %s — the build generates it", dest, project, props.Type, project, dest)
	}

	text = sandbox.Deps.StringsDeps.ReplaceAll(text, "type "+props.Type+" struct", "type "+projectPropsType+" struct")
	text = sandbox.Deps.StringsDeps.ReplaceAll(text, "// "+props.Type+" is ", "// "+projectPropsType+" is ")
	sandbox.Deps.StdDeps.Logf("build: moved %s to %s; %s now embeds its fields\n", dest, project, props.Type)
	return io.WriteFile(project, []byte(text))
}
