package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// databaseGeneratedFile is one file of a database package and the template it
// is rendered from.
type databaseGeneratedFile struct {
	Name     string
	Template string
}

// databaseGeneratedFiles is the three files of a database package a build
// writes, in render order. methods_custom.go is not here and never will be: it
// is the one file of the package agnos neither reads nor rewrites.
var databaseGeneratedFiles = []databaseGeneratedFile{
	{"api.go", "templates/database_api.go"},
	{"new.go", "templates/database_new.go"},
	{"methods.go", "templates/database_methods.go"},
}

// GenerateDatabaseNew renders the whole generated half of every declared
// database: api.go (the records and the struct of function fields), new.go
// (the database.Props and the wiring) and methods.go (the body of every
// method), each from its own template under assets/templates.
//
// It is the database layer's GenerateRouteNew, only wider: a route declares
// one generated file and a database three, because a database has no generic
// dispatch to read its declaration back at runtime — its methods are typed by
// table, so they are spelled out.
func GenerateDatabaseNew(sandbox *api.Sandbox, io *smartio.SmartIO, databases []map[string]any, module string) error {
	for _, database := range databases {
		pkg, _ := database["Package"].(string)
		if pkg == "" {
			continue
		}

		vars := map[string]any{
			"Module":        module,
			"GeneratorName": generatorName(sandbox),
			"CustomFile":    utils.DatabaseCustomFile,
		}
		for key, value := range database {
			vars[key] = value
		}

		for _, file := range databaseGeneratedFiles {
			dest := utils.DatabasesDir + "/" + pkg + "/" + file.Name
			if err := utils.RenderTemplateToDest(sandbox, io, file.Template, vars, dest); err != nil {
				return err
			}
		}
	}
	return nil
}
