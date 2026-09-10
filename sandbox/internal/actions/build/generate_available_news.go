package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// GenerateAvailableNews renders assets/templates/available_new.go once per
// declared available into adapters/availables/<name>/new.go — the New() that
// binds exactly the adapters that available.yaml selects. An available with no
// declaration is a hand-written mix and is not touched.
func GenerateAvailableNews(deps *deps.Deps, io *smartio.SmartIO, availables []map[string]any, module string) error {
	for _, available := range availables {
		name, _ := available["Name"].(string)
		if name == "" {
			continue
		}

		vars := map[string]any{"Module": module}
		for key, value := range available {
			vars[key] = value
		}

		dest := utils.AvailableDir(name) + "/new.go"
		if err := utils.RenderTemplateToDest(deps, io, "templates/available_new.go", vars, dest); err != nil {
			return err
		}
	}
	return nil
}
