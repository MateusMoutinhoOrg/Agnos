package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
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
func MigrateLegacyProps(sandbox *api.Sandbox, io *smartio.SmartIO, hasCli bool, hasServer bool) error {
	if hasServer {
		if err := migrateLegacyProps(sandbox, io, utils.RoutePropsFile, utils.RoutePropsDir); err != nil {
			return err
		}
	}
	if hasCli {
		if err := migrateLegacyProps(sandbox, io, utils.CommandPropsFile, utils.CommandPropsDir); err != nil {
			return err
		}
	}
	return nil
}

// migrateLegacyProps moves sandbox/api/<file> to <dir>/<file>, its package
// clause renamed to the last segment of dir.
func migrateLegacyProps(sandbox *api.Sandbox, io *smartio.SmartIO, file string, dir string) error {
	legacy := legacyPropsDir + "/" + file
	dest := dir + "/" + file
	if !io.IsFile(legacy) || io.IsFile(dest) {
		return nil
	}

	content, err := io.ReadFile(legacy)
	if err != nil {
		return err
	}

	lines := sandbox.Deps.Stringsdeps.Split(string(content), "\n")
	for i, line := range lines {
		if sandbox.Deps.Stringsdeps.TrimSpace(line) == "package api" {
			lines[i] = "package " + lastSegmentOf(sandbox, dir)
			break
		}
	}

	if err := io.WriteFileOverwrite(dest, []byte(sandbox.Deps.Stringsdeps.Join(lines, "\n"))); err != nil {
		return err
	}
	io.RemoveDir(legacy)
	return nil
}
