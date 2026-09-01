package build

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/utils"
)

func Render_sandbox_deps_deps_go(deps *deps.Deps, io *smartio.SmartIO, module string) error {

	dirs := io.ListDirs("sandbox/deps")

	var libs []map[string]string
	for _, dir := range dirs {
		parts := strings.Split(dir, "/")
		name := parts[len(parts)-1]

		if len(name) == 0 {
			continue
		}

		title := strings.ToUpper(name[:1]) + name[1:]
		libs = append(libs, map[string]string{
			"Name":  name,
			"Title": title,
		})
	}

	vars := map[string]interface{}{
		"Module": module,
		"Libs":   libs,
	}
	err := utils.RenderTemplateToDest(deps, io, "sandbox/deps/deps.go", vars, "sandbox/deps/deps.go")
	if err != nil {
		return err
	}
	return nil
}
