package drop_skills

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// DropSkillsInternal turns the skills extension on in the project's
// declaration. It writes nothing else: the skill is the generator's commands
// as the binary running the build knows them, so the build that follows
// renders it — and every later one renders it again, which is what keeps it
// in step with the generator. Run on a project that has it on already, it
// changes nothing but that build.
func DropSkillsInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, path string) error {
	sandbox.Deps.StdDeps.Logf("drop-skills started with path %s \n", path)
	return utils.SetExtension(sandbox, io, utils.ExtensionSkills, true)
}
