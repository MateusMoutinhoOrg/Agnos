package build

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/utils"
)

func Render_sandbox_new_go(deps *deps.Deps, io *smartio.SmartIO, module string) error {

	files := io.ListFiles("sandbox/binds")

	var binds []string
	for _, file := range files {
		parts := strings.Split(file, "/")
		name := parts[len(parts)-1]

		if !strings.HasSuffix(name, ".go") {
			continue
		}

		baseName := strings.TrimSuffix(name, ".go")
		if len(baseName) == 0 {
			continue
		}

		title := strings.ToUpper(baseName[:1]) + baseName[1:]
		binds = append(binds, title+"Bind")
	}

	vars := map[string]interface{}{
		"Module":  module,
		"HasDeps": smartio.IsDir(deps, io, "sandbox/deps"),
		"Binds":   binds,
	}
	err := utils.RenderTemplateToDest(deps, io, "sandbox/new.go", vars, "sandbox/new.go")
	if err != nil {
		return err
	}
	return nil
}
