package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// databasePagesDir is the doc whose doc.md indexes one page per database.
const databasePagesDir = utils.DocsDir + "/Databases"

// GenerateDatabasePages is the database layer's GenerateRoutePages: it renders
// assets/templates/database_page.md once per declared database into
// docs/Databases/<db>.md, the page docs/Databases' own doc.md links to. A
// database is named by the package that declares it, so the page of a database
// is found from its directory and the other way round.
//
// A page whose database is gone is removed here, so the directory holds the
// databases that are declared now and no page nothing links to.
func GenerateDatabasePages(sandbox *api.Sandbox, io *smartio.SmartIO, docs []DatabaseDoc) error {
	written := map[string]bool{}

	for _, doc := range docs {
		file, err := databasePageFile(sandbox, doc.Package)
		if err != nil {
			return err
		}

		if err := utils.RenderTemplateToDest(sandbox, io, "templates/database_page.md",
			map[string]any{"Database": doc, "GeneratorName": generatorName(sandbox)}, file); err != nil {
			return err
		}
		written[file] = true
	}

	removeStaleDocPages(sandbox, io, databasePagesDir, written)
	return nil
}

// databasePageFile is the page a database is written to. A name that spells
// one of the two reserved names of a doc directory is a hard error rather than
// a page silently overwriting the index it is linked from.
func databasePageFile(sandbox *api.Sandbox, name string) (string, error) {
	file := name + docPageExt

	if file == utils.DocFile || file == utils.DocIndexFile {
		return "", sandbox.Deps.Std.Errorf(
			"database %s cannot be documented: its page would be %s/%s, which is the doc's own %s",
			name, databasePagesDir, file, file)
	}

	return databasePagesDir + "/" + file, nil
}
