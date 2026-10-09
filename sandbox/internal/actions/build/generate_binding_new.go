package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// GenerateBindingNewFiles renders assets/templates/binding_new.go once per
// declared binding into adapters/bindings/<name>/generated.new.go — the New() that
// binds exactly the adapters that binding.yaml selects. A binding with no
// declaration is a hand-written mix and is not touched.
func GenerateBindingNewFiles(sandbox *api.Sandbox, io *stagedfs.StagedFS, bindings []map[string]any, module string) error {
	for _, binding := range bindings {
		name, _ := binding["BindingName"].(string)
		if name == "" {
			continue
		}

		vars := map[string]any{"Module": module}
		for key, value := range binding {
			vars[key] = value
		}

		dest := utils.BindingDir(name) + "/" + utils.GeneratedFile(sandbox, "new.go")
		if err := utils.RenderTemplateToDest(sandbox, io, "templates/binding_new.go", vars, dest); err != nil {
			return err
		}
	}
	return nil
}
