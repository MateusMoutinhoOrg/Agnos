package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// apiPartsDir is the directory, inside a group and inside the project, holding
// the parts api.Sandbox and api.Config embed.
const apiPartsDir = "sandbox/api"

// GenerateApiParts renders, ahead of every other asset, each part of
// api.Sandbox or api.Config an enabled group carries — cli's
// clisandbox.go, server's serversandbox.go — and returns the paths it
// wrote. sandbox/api/sandbox.go and config.go are rendered from the parts the
// build collects, and the collection lists disk: on a project's first build
// the parts the groups write in the same transaction would be missed without
// this.
//
// A part is rendered over vars, which is less than the full set: it is a
// struct of contracts and needs no more than the module. The group renders
// the same asset again later over the full set, to the same bytes.
func GenerateApiParts(sandbox *api.Sandbox, io *stagedfs.StagedFS, groups []string, vars map[string]interface{}) ([]string, error) {
	var written []string
	for _, group := range groups {
		files, err := sandbox.Deps.EmbedDeps.ListFilesRecursively(group)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			name := lastSegmentOf(sandbox, file)
			if file != apiPartsDir+"/"+name {
				continue
			}
			if !utils.IsSandboxPart(sandbox, name) && !utils.IsConfigPart(sandbox, name) {
				continue
			}
			if err := utils.RenderTemplateToDest(sandbox, io, group+"/"+file, vars, file); err != nil {
				return nil, err
			}
			written = append(written, file)
		}
	}
	return written, nil
}
