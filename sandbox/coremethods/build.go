package coremethods

type BuildProps struct {
	Path    string
	Project string
	Force   bool
}

func (self *CoreMethods) Build(props BuildProps) error {
	self.Sandbox.Deps.Printf("build started with path %s \n", props.Path)
	return nil
}
