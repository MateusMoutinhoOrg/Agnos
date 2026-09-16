package {{.Package}}

import (
	api "{{.Module}}/sandbox/api"
	{{.Package}} "{{.Module}}/sandbox/internal/{{.Package}}"
)

// Constructor fills Sandbox.{{.Name}}, building it with the
// New{{.Name}} of sandbox/internal/{{.Package}}. sandbox/new.go calls it
// once, along with the Constructor of every other package under
// sandbox/constructors/.
//
// Written once by `{{.GeneratorName}} build` and then yours: wrap the
// implementation, decorate the contract, or build a different one entirely.
// No build rewrites this file once it is there.
func Constructor(sandbox *api.Sandbox) {
	sandbox.{{.Name}} = {{.Package}}.New{{.Name}}(sandbox)
}
