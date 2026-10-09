package {{.Package}}

import (
	api "{{.Module}}/sandbox/api"
	{{.Package}} "{{.Module}}/{{.Source}}"
)

// Constructor fills Sandbox.{{.ContractName}}, building it with the
// New{{.ContractName}} of {{.Source}}. sandbox/generated.new.go calls it
// once, along with the Constructor of every other package under
// sandbox/constructors/.
//
// Written once by `{{.GeneratorName}} build` and then yours: wrap the
// implementation, decorate the contract, or build a different one entirely.
// No build rewrites this file once it is there.
func Constructor(sandbox *api.Sandbox) {
	sandbox.{{.ContractName}} = {{.Package}}.New{{.ContractName}}(sandbox)
}
