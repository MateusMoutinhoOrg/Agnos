package userconfig

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	serializibles "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/serializebles"
)

type ProjectConf struct {
	Name        string
	Module      string
	Description string

	Persist func() error
}

func NewProjectConf(sandbox *api.SandBox, path string) (*ProjectConf, error) {

	var project_specs *serializibles.SerializibleObject
	if sandbox.Deps.IoLib.IsFile(path) {
		content_bytes, fileerror := sandbox.Deps.IoLib.ReadFile(path)
		if fileerror != nil {
			return nil, fileerror
		}
		specs, parse_error := sandbox.Deps.SerializeLib.ParseYaml(string(content_bytes))
		if parse_error != nil {
			return nil, parse_error
		}
		project_specs = specs

	} else {
		project_specs = sandbox.Deps.SerializeLib.CreateObject()
	}

	if !project_specs.IsObject() {
		return nil, sandbox.Deps.Errorf("project_specs is not an object")
	}
	name_item := project_specs.GetObjectItem("name")
	module_item := project_specs.GetObjectItem("module")
	description_item := project_specs.GetObjectItem("description")

	if !name_item.IsString() || !module_item.IsString() || !description_item.IsString() {
		return nil, sandbox.Deps.Errorf("project_specs is not an object")
	}

	project_conf := &ProjectConf{}
	project_conf.Name, _ = name_item.GetString()
	project_conf.Module, _ = module_item.GetString()
	project_conf.Description, _ = description_item.GetString()

	project_conf.Persist = func() error {
		bytes := sandbox.Deps.SerializeLib.SerializeToYaml(project_specs)
		return sandbox.Deps.IoLib.WriteFile(path, []byte(bytes))
	}

	return project_conf
}
