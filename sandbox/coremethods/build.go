package coremethods

import "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/core"

func (self *CoreMethods) Build(props core.BuildProps) error {
	self.Sandbox.Deps.Printf("build started with path %s \n", props.Path)
	return nil
}
