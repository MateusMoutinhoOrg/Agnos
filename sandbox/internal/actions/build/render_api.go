package build

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/utils"
)

func Render_sandbox_api_sandbox_go(deps *deps.Deps, io *smartio.SmartIO) error {
	files := io.ListFiles("sandbox/api")

	var constructors []string
	for _, file := range files {
		// Extract base name from the path (in case it returns full paths or relative paths)
		parts := strings.Split(file, "/")
		name := parts[len(parts)-1]

		if name == "sandbox.go" || !strings.HasSuffix(name, ".go") {
			continue
		}
		
		baseName := strings.TrimSuffix(name, ".go")
		if len(baseName) > 0 {
			title := strings.ToUpper(baseName[:1]) + baseName[1:]
			constructors = append(constructors, title)
		}
	}

	vars := map[string]interface{}{
		"Constructors": constructors,
	}
	err := utils.RenderTemplateToDest(deps, io, "sandbox/api/sandbox.go", vars, "sandbox/api/sandbox.go")
	if err != nil {
		return err
	}
	return nil
}
