package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// commandPagesDir is the doc whose doc.md indexes one page per command.
const commandPagesDir = utils.DocsDir + "/Commands"

// GenerateCommandPages renders assets/templates/command_page.md once per
// visible command into docs/Commands/<identifier>.md, the page docs/Commands'
// own doc.md links to. A page is an asset of that doc directory, not a sub-doc:
// CollectDocTree walks directories, so a plain .md beside doc.md is ignored by
// the index and by `verify`, exactly as the generated Index.md is.
//
// One page per command is what keeps a lookup cheap: reading what `add-flag`
// declares costs that command's page, not every command of the project.
//
// A page whose command is gone is removed here, so the directory holds the
// commands that are declared now and no page nothing links to.
func GenerateCommandPages(sandbox *api.Sandbox, io *smartio.SmartIO, groups []CommandDocGroup, name string) error {
	written := map[string]bool{}

	for _, group := range groups {
		for _, command := range group.Commands {
			file, err := commandPageFile(sandbox, command.Identifier)
			if err != nil {
				return err
			}

			vars := map[string]any{
				"Name":     name,
				"Category": group.Category,
				"Command":  command,
			}

			if err := utils.RenderTemplateToDest(sandbox, io, "templates/command_page.md", vars, file); err != nil {
				return err
			}
			written[file] = true
		}
	}

	removeStaleDocPages(sandbox, io, commandPagesDir, written)
	return nil
}

// commandPageFile is the page a command is written to. An identifier that
// spells one of the two reserved names of a doc directory is a hard error
// rather than a page silently overwriting the index it is linked from.
func commandPageFile(sandbox *api.Sandbox, identifier string) (string, error) {
	file := identifier + docPageExt

	if file == utils.DocFile || file == utils.DocIndexFile {
		return "", sandbox.Deps.Std.Errorf(
			"command %s cannot be documented: its page would be %s/%s, which is the doc's own %s",
			identifier, commandPagesDir, file, file)
	}

	return commandPagesDir + "/" + file, nil
}
