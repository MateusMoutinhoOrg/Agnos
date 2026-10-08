package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// commandPagesDir is the doc whose doc.md indexes one page per command.
const commandPagesDir = utils.DocsDir + "/Commands"

// GenerateCommandPages renders assets/templates/command_page.md once per
// visible command — and once per middleware — into docs/Commands/<identifier>.md, the page docs/Commands'
// own doc.md links to. A page is an asset of that doc directory, not a sub-doc:
// CollectDocTree walks directories, so a plain .md beside doc.md is ignored by
// the index and by `verify`, exactly as the generated Index.md is.
//
// One page per command is what keeps a lookup cheap: reading what `add-flag`
// declares costs that command's page, not every command of the project.
//
// A page whose command is gone is removed here, so the directory holds the
// commands that are declared now and no page nothing links to.
func GenerateCommandPages(sandbox *api.Sandbox, io *stagedfs.StagedFS, docs CommandDocs, name string) error {
	written := map[string]bool{}

	render := func(category string, command CommandDoc) error {
		file, err := commandPageFile(sandbox, command.Page)
		if err != nil {
			return err
		}
		vars := map[string]any{
			"ProjectName": name,
			"Category":    category,
			"Command":     command,
		}
		if err := utils.RenderTemplateToDest(sandbox, io, "templates/command_page.md", vars, file); err != nil {
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

	removeStaleDocPages(sandbox, io, commandPagesDir, written)
	return nil
}

// commandPageFile is the page a command is written to. An identifier that
// spells one of the two reserved names of a doc directory is a hard error
// rather than a page silently overwriting the index it is linked from.
func commandPageFile(sandbox *api.Sandbox, file string) (string, error) {
	identifier := sandbox.Deps.StringsDeps.TrimSuffix(file, docPageExt)

	if file == utils.DocFile || file == utils.DocIndexFile {
		return "", sandbox.Deps.StdDeps.Errorf(
			"command %s cannot be documented: its page would be %s/%s, which is the doc's own %s",
			identifier, commandPagesDir, file, file)
	}

	return commandPagesDir + "/" + file, nil
}
