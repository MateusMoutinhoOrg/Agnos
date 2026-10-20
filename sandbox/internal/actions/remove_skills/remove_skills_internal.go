package remove_skills

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveSkillsInternal removes the generator's skill directory from every
// root of utils.SkillRoots, whole, then each directory the removal left
// empty — .agents/ and .claude/ only when nothing else is in them, since
// .claude/ holds a tool's own settings too. On a project that never had the
// skill it changes nothing but the key.
//
// Another skill under the same roots is left where it is: only the one named
// after the generator is this mechanic's.
func RemoveSkillsInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, path string) error {
	sandbox.Deps.StdDeps.Logf("remove-skills started with path %s \n", path)

	if err := utils.SetExtension(sandbox, io, utils.ExtensionSkills, false); err != nil {
		return err
	}

	removed := []string{}
	for _, dir := range utils.SkillDirs(sandbox) {
		if !io.IsDir(dir) {
			continue
		}
		removed = append(removed, io.ListAllRecursively(dir)...)
		removed = append(removed, dir)
	}
	utils.RemoveTree(sandbox, io, removed)
	return nil
}
