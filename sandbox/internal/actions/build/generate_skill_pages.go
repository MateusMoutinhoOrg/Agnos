package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/extensionsconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// SkillExtension is one row of the extensions table a skill prints: the key,
// whether this project has it on, and the line list-extensions prints for it.
type SkillExtension struct {
	Name    string
	Enabled bool
	Help    string
}

// GenerateSkillPages renders the generator's skill into every root of
// utils.SkillRoots: assets/templates/skill.md into <dir>/SKILL.md, and
// assets/templates/skill_command.md once per visible command — and once per
// middleware — into <dir>/commands/<identifier>.md, the page SKILL.md's index
// links to. Every root gets the same bytes, so whichever tool reads the
// project finds the same skill.
//
// One page per command keeps a lookup cheap, as in docs/Commands: an agent
// reads the index once and the page of the command it is about to type.
//
// A page whose command is gone is removed here, so a skill rendered by a newer
// binary holds no page of a command that binary no longer answers.
func GenerateSkillPages(sandbox *api.Sandbox, io *stagedfs.StagedFS, docs CommandDocs, extensions_conf *extensionsconf.ExtensionsConf, vars map[string]any) error {
	extensions := []SkillExtension{}
	for _, spec := range utils.ExtensionCatalog() {
		extensions = append(extensions, SkillExtension{
			Name:    spec.Name,
			Enabled: extensions_conf.IsEnabled(spec.Name),
			Help:    docCell(sandbox, spec.Help),
		})
	}

	skill_vars := map[string]any{}
	for key, value := range vars {
		skill_vars[key] = value
	}
	skill_vars["Skill"] = docs
	skill_vars["Extensions"] = extensions

	for _, dir := range utils.SkillDirs(sandbox) {
		if err := utils.RenderTemplateToDest(sandbox, io, "templates/skill.md", skill_vars, dir+"/"+utils.SkillFile); err != nil {
			return err
		}

		commands_dir := dir + "/" + utils.SkillCommandsDir
		written := map[string]bool{}

		render := func(category string, command CommandDoc) error {
			file := commands_dir + "/" + command.Page
			page_vars := map[string]any{
				"GeneratorName": vars["GeneratorName"],
				"SkillFile":     utils.SkillFile,
				"Category":      category,
				"Command":       command,
			}
			if err := utils.RenderTemplateToDest(sandbox, io, "templates/skill_command.md", page_vars, file); err != nil {
				return err
			}
			written[file] = true
			return nil
		}

		for _, group := range docs.Groups {
			for _, command := range group.Commands {
				if err := render(group.Category, command); err != nil {
					return err
				}
			}
		}
		for _, middleware := range docs.Middlewares {
			if err := render("Middlewares", middleware); err != nil {
				return err
			}
		}

		removeStaleDocPages(sandbox, io, commands_dir, written)
	}
	return nil
}
