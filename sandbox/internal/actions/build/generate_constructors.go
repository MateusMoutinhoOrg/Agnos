package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// GenerateConstructors renders assets/templates/constructor.go into
// sandbox/constructors/<x>/constructor.go once per contract of sandbox/api/
// that has a new.go to call, at utils.ConstructorSource — the Constructor(sandbox)
// that fills Sandbox.<X>.
//
// It is written **once**. A constructor already on disk is left exactly as it
// is, however far it has drifted from what this template renders: how a field
// of the Sandbox is built is the project's to change, and rewriting the file
// every build is precisely what kept it from being. What every build does own
// is sandbox/new.go, which calls whatever packages are there.
func GenerateConstructors(sandbox *api.Sandbox, io *smartio.SmartIO, constructors []Constructor, module string) error {
	for _, constructor := range constructors {
		if !constructor.HasNew {
			continue
		}

		dest := utils.ConstructorPath(constructor.Package)
		if io.IsFile(dest) {
			continue
		}

		vars := map[string]any{
			"Module":        module,
			"Name":          constructor.Name,
			"Package":       constructor.Package,
			"Source":        constructor.Source,
			"GeneratorName": generatorName(sandbox),
		}

		if err := utils.RenderTemplateToDest(sandbox, io, "templates/constructor.go", vars, dest); err != nil {
			return err
		}
	}
	return nil
}
