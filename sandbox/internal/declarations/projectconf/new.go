package projectconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func New(sandbox *api.Sandbox, content string) (*ProjectConf, error) {

	if content == "" {
		return nil, sandbox.Deps.StdDeps.Errorf("content cannot be empty, use NewProjectConfEmpty instead")
	}

	specs, parse_error := sandbox.Deps.SerializableDeps.ParseYaml(content)
	if parse_error != nil {
		return nil, parse_error
	}
	project_specs := specs

	if !project_specs.IsObject() {
		return nil, sandbox.Deps.StdDeps.Errorf("project_specs is not an object")
	}
	name_item, _ := project_specs.GetObjectItem("project-name")
	version_item, _ := project_specs.GetObjectItem("version")

	// Version defaults to empty: `version: null` (or absent) means "no version yet".
	project_conf := &ProjectConf{}
	var err error

	if name_item != nil && !name_item.IsNull() {
		project_conf.ProjectName, err = name_item.GetString()
		if err != nil {
			return nil, sandbox.Deps.StdDeps.Errorf("project-name is not a string")
		}
	}

	if version_item != nil && !version_item.IsNull() {
		project_conf.Version, err = version_item.GetString()
		if err != nil {
			return nil, sandbox.Deps.StdDeps.Errorf("version is not a string")
		}
	}

	BindMethods(sandbox, project_conf)
	return project_conf, nil
}
