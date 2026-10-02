package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// legacyApiParts pairs the prefix a part of api.Sandbox or api.Config was
// named with before the parts were read by suffix with the suffix it is named
// with now: usersandbox_<x>.go is <x>sandbox.go, userconfig_<x>.go is
// <x>config.go.
var legacyApiParts = [][2]string{
	{"usersandbox_", utils.SandboxPartSuffix},
	{"userconfig_", utils.ConfigPartSuffix},
}

// MigrateLegacyApiParts renames every part of sandbox/api/ an older build or
// mechanic wrote under the old prefix — userconfig_backoffice.go — to the
// name the aggregates embed it under now — backofficeconfig.go. The file is
// the project's, so it is moved as it is; one whose new name is already taken
// is left alone. It returns the paths it wrote: the listing that collects the
// parts reads disk, where a moved file is not yet.
//
// It runs before anything reads sandbox/api/, so the part is never taken for
// a contract of the Sandbox, nor dropped from the aggregate it belongs to.
func MigrateLegacyApiParts(sandbox *api.Sandbox, io *smartio.SmartIO) ([]string, error) {
	var moved []string
	for _, file := range io.ListFiles(legacyPropsDir) {
		name := lastSegmentOf(sandbox, file)
		for _, pair := range legacyApiParts {
			if !sandbox.Deps.Stringsdeps.HasPrefix(name, pair[0]) || !sandbox.Deps.Stringsdeps.HasSuffix(name, ".go") {
				continue
			}
			unit := sandbox.Deps.Stringsdeps.TrimSuffix(sandbox.Deps.Stringsdeps.TrimPrefix(name, pair[0]), ".go")
			dest := legacyPropsDir + "/" + unit + pair[1]
			if unit == "" || io.IsFile(dest) {
				continue
			}

			content, err := io.ReadFile(file)
			if err != nil {
				return nil, err
			}
			if err := io.WriteFileOverwrite(dest, content); err != nil {
				return nil, err
			}
			io.RemoveDir(file)
			moved = append(moved, dest)
		}
	}
	return moved, nil
}
