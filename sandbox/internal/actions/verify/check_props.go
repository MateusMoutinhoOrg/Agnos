package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// configOwnFields are the fields api.Config declares itself, beside the parts
// it embeds.
var configOwnFields = []string{"ProjectName", "Version"}

// CheckProps reports every field two parts of one generated aggregate both
// declare — api.Config, api.Sandbox, RouteProps, CommandProps — and every part
// field that shadows one the aggregate declares itself. Each aggregate embeds
// one struct per file, so the parts are written by different hands (the
// project, backoffice-init), and Go only refuses the ambiguous selector where a
// handler reads it.
func CheckProps(sandbox *api.Sandbox, io *smartio.SmartIO) []string {
	var violations []string

	check := func(aggregate string, dir string, accept func(name string) bool, own []string) {
		if !io.IsDir(dir) {
			return
		}
		structs, err := utils.CollectEmbeddedStructs(sandbox, io, dir, accept, nil)
		if err != nil {
			violations = append(violations, err.Error())
			return
		}
		violations = append(violations, utils.EmbeddedCollisions(sandbox, aggregate, structs, own)...)
	}

	check("api.Config", "sandbox/api", utils.HasFilePrefix(sandbox, utils.UserConfigPrefix), configOwnFields)
	check("api.Sandbox", "sandbox/api", utils.HasFilePrefix(sandbox, utils.UserSandboxPrefix), sandboxOwnFields(sandbox, io))
	check("routeprops.RouteProps", utils.RoutePropsDir, utils.AllBut(utils.RoutePropsFile), nil)
	check("commandprops.CommandProps", utils.CommandPropsDir, utils.AllBut(utils.CommandPropsFile), nil)

	return violations
}

// sandboxOwnFields are the fields api.Sandbox declares itself: Deps, and one
// per contract of sandbox/api/.
func sandboxOwnFields(sandbox *api.Sandbox, io *smartio.SmartIO) []string {
	own := []string{"Deps"}
	for _, file := range io.ListFiles("sandbox/api") {
		name := lastSegment(sandbox, file)
		if !sandbox.Deps.Stringsdeps.HasSuffix(name, ".go") || utils.IsConstructorExempt(sandbox, name) {
			continue
		}
		base := sandbox.Deps.Stringsdeps.TrimSuffix(name, ".go")
		own = append(own, sandbox.Deps.Stringsdeps.ToUpper(base[:1])+base[1:])
	}
	return own
}
