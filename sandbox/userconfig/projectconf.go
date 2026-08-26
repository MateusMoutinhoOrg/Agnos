package userconfig

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	serializibles "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/serializebles"
)

type ProjectConf struct {
	SetName        func(name string)
	SetModule      func(module string)
	SetDescription func(description string)

	GetName        func() string
	GetModule      func() string
	GetDescription func() string

	Persist func() error
}

func NewProjectConf(sandbox *api.SandBox, path string) ProjectConf {

	var project_specs serializibles.SerializibleObject
	if sandbox.Deps.IoLib.IsFile(path) {
		content_bytes, _ := sandbox.Deps.IoLib.ReadFile(path)
		project_specs, _ = sandbox.Deps.SerializeLib.ParseYaml(string(content_bytes))

	} else {
		project_specs = sandbox.Deps.SerializeLib.CreateObject()
	}
	// check for types

	return ProjectConf{
		SetName: func(name string) {
			project_specs.AddItemToObject("name", name)
		},
		SetModule: func(module string) {
			project_specs.AddItemToObject("module", module)
		},
		SetDescription: func(description string) {
			project_specs.AddItemToObject("description", description)
		},

		GetName: func() string {
			name_obj := project_specs.GetObjectItem("name")
			name_str, _ := name_obj.GetString()
			return name_str
		},
		GetModule: func() string {
			module_obj := project_specs.GetObjectItem("module")
			module_str, _ := module_obj.GetString()
			return module_str
		},
		GetDescription: func() string {
			description_obj := project_specs.GetObjectItem("description")
			description_str, _ := description_obj.GetString()
			return description_str
		},

		Persist: func() error {
			bytes := sandbox.Deps.SerializeLib.SerializeToYaml(project_specs)
			return sandbox.Deps.IoLib.WriteFile(path, []byte(bytes))
		},
	}
}
