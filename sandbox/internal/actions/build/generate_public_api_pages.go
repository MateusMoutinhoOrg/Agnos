package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// publicApiPagesDir is the doc whose doc.md indexes one page per contract.
const publicApiPagesDir = utils.DocsDir + "/PublicApi"

// GeneratePublicApiPages renders assets/templates/public_api_page.md once per
// file of sandbox/api and once per contract of sandbox/deps into
// docs/PublicApi/<page>, the pages docs/PublicApi' own doc.md links to. Like a
// command's page, each one is an asset of that doc directory rather than a
// sub-doc: CollectDocTree walks directories, so a plain .md beside doc.md is
// ignored by the index and by `verify`.
//
// The page a contract lands on is the one its collector named, so the link on
// the index and the file written here can never disagree.
func GeneratePublicApiPages(sandbox *api.Sandbox, io *smartio.SmartIO, public_api []map[string]any, deps_api []map[string]any) error {
	written := map[string]bool{}

	for _, group := range public_api {
		path, _ := group["Path"].(string)
		vars := map[string]any{
			"Title": "`" + path + "`",
			"Files": []map[string]any{group},
		}
		if err := renderPublicApiPage(sandbox, io, group, vars, written); err != nil {
			return err
		}
	}

	for _, contract := range deps_api {
		title, _ := contract["Title"].(string)
		name, _ := contract["Name"].(string)
		files, _ := contract["Files"].([]map[string]any)
		vars := map[string]any{
			"Title":  "`deps." + title + "`",
			"Source": "`sandbox/deps/" + name + "`",
			"Files":  files,
		}
		if err := renderPublicApiPage(sandbox, io, contract, vars, written); err != nil {
			return err
		}
	}

	removeStaleDocPages(sandbox, io, publicApiPagesDir, written)
	return nil
}

// renderPublicApiPage writes one page to the file its collector named and
// records it as written, so a page this build produced is never taken for a
// stale one.
func renderPublicApiPage(sandbox *api.Sandbox, io *smartio.SmartIO, unit map[string]any, vars map[string]any, written map[string]bool) error {
	page, _ := unit["Page"].(string)
	if page == "" {
		return nil
	}

	dest := publicApiPagesDir + "/" + page
	if err := utils.RenderTemplateToDest(sandbox, io, "templates/public_api_page.md", vars, dest); err != nil {
		return err
	}

	written[dest] = true
	return nil
}
