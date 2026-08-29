package parsables

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/sandbox"
)

type ProjectConf struct {
	Name        string
	Version     string
	Description string

	Render func() string
}

func addProjectConfMethods(sandbox *sandbox.SandBox, project_conf *ProjectConf) {
	project_conf.Render = func() string {
		obj := sandbox.Deps.SerializeLib.CreateObject()
		obj.AddItemToObject("name", project_conf.Name)
		obj.AddItemToObject("version", project_conf.Version)
		obj.AddItemToObject("description", project_conf.Description)

		return sandbox.Deps.SerializeLib.SerializeToYaml(obj)
	}
}

func NewProjectConf(sandbox *sandbox.SandBox, content string) (*ProjectConf, error) {

	if content == "" {
		return nil, sandbox.Deps.Errorf("content cannot be empty, use NewProjectConfEmpty instead")
	}

	specs, parse_error := sandbox.Deps.SerializeLib.ParseYaml(content)
	if parse_error != nil {
		return nil, parse_error
	}
	project_specs := specs

	if !project_specs.IsObject() {
		return nil, sandbox.Deps.Errorf("project_specs is not an object")
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
			return nil, sandbox.Deps.Errorf("name is not a string")
		}
	}

	if description_item != nil && !description_item.IsNull() {
		project_conf.Description, err = description_item.GetString()
		if err != nil {
			return nil, sandbox.Deps.Errorf("description is not a string")
		}
	}

	if version_item != nil && !version_item.IsNull() {
		project_conf.Version, err = version_item.GetString()
		if err != nil {
			return nil, sandbox.Deps.Errorf("version is not a string")
		}
	}

	addProjectConfMethods(sandbox, project_conf)
	return project_conf, nil
}

func NewProjectConfEmpty(sandbox *sandbox.SandBox) *ProjectConf {
	project_conf := &ProjectConf{
		Version: "v0.0.0",
	}
	addProjectConfMethods(sandbox, project_conf)
	return project_conf
}
