// This file is an illustrative sample, not part of the build. It shows the
// five files of one parsable in a single listing; in the tree they are separate.

// ─── api.go ─────────────────────────────────────────────────────────────────

package projectconf

type ProjectConf struct {
	Name        string
	Version     string
	Description string

	Render func() string
}

// ─── new.go ─────────────────────────────────────────────────────────────────

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"

func New(deps *deps.Deps, content string) (*ProjectConf, error) {
	root, err := deps.Serializables.ParseYaml(content)
	if err != nil {
		return nil, err
	}
	self := NewEmpty(deps)
	if root.HasKey("name") {
		item, _ := root.GetObjectItem("name")
		if self.Name, err = item.GetString(); err != nil {
			return nil, deps.Std.Errorf("project.yaml: name must be a string")
		}
	}
	// … version, description the same way
	return self, nil
}

// ─── new_empty.go ───────────────────────────────────────────────────────────

func NewEmpty(deps *deps.Deps) *ProjectConf {
	self := &ProjectConf{}
	bindMethods(deps, self)
	return self
}

// ─── bind_methods.go ────────────────────────────────────────────────────────

func bindMethods(deps *deps.Deps, self *ProjectConf) {
	self.Render = func() string { return render(deps, self) }
}

// ─── render.go ──────────────────────────────────────────────────────────────

func render(deps *deps.Deps, self *ProjectConf) string {
	root := deps.Serializables.CreateObject()
	root.AddItemToObject("name", self.Name)
	root.AddItemToObject("version", self.Version)
	root.AddItemToObject("description", self.Description)
	return deps.Serializables.SerializeToYaml(root)
}
