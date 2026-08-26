package userconfig

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	serializibles "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/serializebles"
)

type ProjectConf struct {
	Name        string
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
	description_item, _ := project_specs.GetObjectItem("description")

	project_conf := ProjectConf{}
	var err error

	if name_item != nil && !name_item.IsNull() {
		project_conf.Name, err = name_item.GetString()
		if err != nil {
			return nil, sandbox.Deps.Errorf("name is not a string")
		}
	}

	if description_item != nil && !description_item.IsNull() {
		project_conf.Description, err = description_item.GetString()
		if err != nil {
			return nil, sandbox.Deps.Errorf("description is not a string")
		}
	}

	project_conf.Persist = func() error {
		project_specs.AddItemToObject("name", project_conf.Name)
		project_specs.AddItemToObject("description", project_conf.Description)

		bytes := sandbox.Deps.SerializeLib.SerializeToYaml(project_specs)
		return sandbox.Deps.IoLib.WriteFile(path, []byte(bytes))
	}

	return &project_conf, nil
}
