package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// CheckStructure enforces that every item of
// <ProjectName>Config/structure.yaml still names something on disk. The
// declaration is what docs/Structure is rendered from, so an item left behind
// by a deleted file — a ghost spec — would publish a tree the project no
// longer has.
//
// An item marked `dir: true` must be a directory and any other must be a file.
// An item whose path holds a pattern character stands for a family of paths
// ("libs/<lib>/<lib>.go"), so the literal part of its path is checked instead:
// the family may be empty, but the directory it would live in has to exist.
//
// A project with no structure.yaml describes nothing and has nothing to check.
func CheckStructure(deps *deps.Deps, io *smartio.SmartIO) []string {
	var violations []string

	if !io.IsFile(utils.StructureConfPath()) {
		return violations
	}

	structure_conf, err := utils.LoadStructureConf(deps, io)
	if err != nil {
		return append(violations, err.Error())
	}

	for _, node := range utils.FlattenStructure(structure_conf.Items) {
		if utils.IsStructurePattern(node.Path) {
			parent := utils.StructureParentPath(node.Path)
			if parent != "" && !io.IsDir(parent) {
				violations = append(violations, ghostSpec(node.Path,
					"the directory "+parent+" it would live in does not exist"))
			}
			continue
		}

		if node.Dir {
			if !io.IsDir(node.Path) {
				violations = append(violations, ghostSpec(node.Path, "there is no such directory"))
			}
			continue
		}

		if !io.IsFile(node.Path) {
			violations = append(violations, ghostSpec(node.Path, "there is no such file"))
		}
	}

	return violations
}

// ghostSpec words one violation the same way for every kind of missing item.
func ghostSpec(path string, reason string) string {
	return utils.StructureConfPath() + " describes " + path + ", but " + reason +
		" (a ghost spec: drop the item or restore the path)"
}
