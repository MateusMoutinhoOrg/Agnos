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
	name_item, _ := project_specs.GetObjectItem("name")
	module_item, _ := project_specs.GetObjectItem("module")
	description_item, _ := project_specs.GetObjectItem("description")

	project_conf := ProjectConf{}
	var err error

	project_conf.Name, err = name_item.GetString()
	if err != nil {
		return nil, sandbox.Deps.Errorf("name is not a string")
	}
	project_conf.Module, err = module_item.GetString()
	if err != nil {
		return nil, sandbox.Deps.Errorf("module is not a string")
	}
	project_conf.Description, err = description_item.GetString()
	if err != nil {
		return nil, sandbox.Deps.Errorf("description is not a string")
	}

	project_conf.Persist = func() error {
		bytes := sandbox.Deps.SerializeLib.SerializeToYaml(project_specs)
		return sandbox.Deps.IoLib.WriteFile(path, []byte(bytes))
	}

	return &project_conf, nil
}
