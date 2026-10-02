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

	check("api.Config", "sandbox/api", utils.ConfigParts(sandbox), configOwnFields)
	check("api.Sandbox", "sandbox/api", utils.SandboxParts(sandbox), sandboxOwnFields)
	check("routeprops.RouteProps", utils.RoutePropsDir, utils.AllBut(utils.RoutePropsFile), nil)
	check("commandprops.CommandProps", utils.CommandPropsDir, utils.AllBut(utils.CommandPropsFile), nil)

	return violations
}

// sandboxOwnFields are the fields api.Sandbox declares itself, beside the
// parts it embeds: the two every project has.
var sandboxOwnFields = []string{"Deps", "Config"}

// CheckContractFields reports every contract of sandbox/api/ a package builds
// — sandbox/internal/<x>/new.go, or the one level down utils.ConstructorSource
// finds — that no part of api.Sandbox declares a field for. Its constructor
// writes sandbox.<X>, and api.Sandbox holds no field per contract any more:
// each is declared by the part of whoever owns it, so without this the project
// fails in the compiler on a name it never declared.
//
// A contract a mechanic ships (sandbox-cli's cli.go) is skipped: the mechanic
// ships the part declaring its field beside it (clisandbox.go), and a build
// renders both — the check would otherwise refuse the very build that writes
// that part into a tree older than it.
func CheckContractFields(sandbox *api.Sandbox, io *smartio.SmartIO) []string {
	if !io.IsDir("sandbox/api") {
		return nil
	}
	structs, err := utils.CollectEmbeddedStructs(sandbox, io, "sandbox/api", utils.SandboxParts(sandbox), nil)
	if err != nil {
		return []string{err.Error()}
	}

	declared := map[string]bool{}
	for _, name := range sandboxOwnFields {
		declared[name] = true
	}
	for _, part := range structs {
		for _, field := range part.Fields {
			declared[field] = true
		}
	}

	var violations []string
	for _, file := range io.ListFiles("sandbox/api") {
		name := lastSegment(sandbox, file)
		if !sandbox.Deps.Stringsdeps.HasSuffix(name, ".go") || utils.IsConstructorExempt(sandbox, name) {
			continue
		}
		base := sandbox.Deps.Stringsdeps.TrimSuffix(name, ".go")
		if !io.IsFile(utils.ConstructorSource(io, base)+"/new.go") || shippedByGroup(sandbox, name) {
			continue
		}
		field := sandbox.Deps.Stringsdeps.ToUpper(base[:1]) + base[1:]
		if declared[field] {
			continue
		}
		violations = append(violations, file+" is a contract no part of api.Sandbox declares a field for; add `"+
			field+" "+field+"` to UserSandbox in sandbox/api/"+utils.UserSandboxFile)
	}
	return violations
}

// shippedByGroup reports a sandbox/api/ file some code group of assets/
// renders: a contract the generator owns, not the project.
func shippedByGroup(sandbox *api.Sandbox, name string) bool {
	for _, group := range utils.AssetGroups() {
		if !group.Code {
			continue
		}
		if _, err := sandbox.Deps.Embeddeps.ReadFile(group.Name + "/sandbox/api/" + name); err == nil {
			return true
		}
	}
	return false
}
