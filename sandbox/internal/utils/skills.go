package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// SkillRoots are the directories an agent looks for skills in, one per
// convention: .agents/skills is the open one, .claude/skills Claude Code's.
// Every build with skills on writes the same skill into each, so whichever
// tool opens the project finds it.
var SkillRoots = []string{".agents/skills", ".claude/skills"}

// SkillFile is the page an agent loads first: its frontmatter is what the tool
// lists, the rest what it reads once the description matches.
const SkillFile = "SKILL.md"

// SkillCommandsDir is the directory of a skill holding one page per command,
// beside its SKILL.md.
const SkillCommandsDir = "commands"

// SkillDirs is the skill directory under every root, named after the
// generator running the build — agnos for agnos — since the skill teaches
// that generator's commands, never the project's own.
func SkillDirs(sandbox *api.Sandbox) []string {
	name := sandbox.Deps.StringsDeps.ToLower(sandbox.Config.ProjectName)
	dirs := make([]string, 0, len(SkillRoots))
	for _, root := range SkillRoots {
		dirs = append(dirs, root+"/"+name)
	}
	return dirs
}
