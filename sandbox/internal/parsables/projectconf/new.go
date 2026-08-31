package projectconf

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func New(deps *deps.Deps, content string) (*ProjectConf, error) {

	if content == "" {
		return nil, deps.Errorf("content cannot be empty, use NewProjectConfEmpty instead")
	}

	specs, parse_error := deps.SerializeLib.ParseYaml(content)
	if parse_error != nil {
		return nil, parse_error
	}
	project_specs := specs

	if !project_specs.IsObject() {
		return nil, deps.Errorf("project_specs is not an object")
	}
	name_item, _ := project_specs.GetObjectItem("name")
	version_item, _ := project_specs.GetObjectItem("version")
	description_item, _ := project_specs.GetObjectItem("description")

	project_conf := &ProjectConf{
		Version: "v0.0.0",
	}
	var err error

	if name_item != nil && !name_item.IsNull() {
		project_conf.Name, err = name_item.GetString()
		if err != nil {
			return nil, deps.Errorf("name is not a string")
		}
	}

	if description_item != nil && !description_item.IsNull() {
		project_conf.Description, err = description_item.GetString()
		if err != nil {
			return nil, deps.Errorf("description is not a string")
		}
	}

	if version_item != nil && !version_item.IsNull() {
		project_conf.Version, err = version_item.GetString()
		if err != nil {
			return nil, deps.Errorf("version is not a string")
		}
	}

	BindMethods(deps, project_conf)
	return project_conf, nil
}
