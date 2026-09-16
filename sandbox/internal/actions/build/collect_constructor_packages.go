package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// CollectConstructorPackages lists the packages under sandbox/constructors/ —
// the directory names, each of which is also its package clause and the import
// sandbox/new.go names — for the {{range .ConstructorPackages}} loop of
// sandbox/new.go, sorted so the call order is the same on every build.
//
// It is the union of two sets, because neither alone is the whole surface. The
// directories on disk carry the ones the project wrote itself, which is what
// makes sandbox/new.go an open list rather than a mirror of sandbox/api/. The
// entries of constructors carry the ones this same build is about to write
// (GenerateConstructors), which a listing cannot see: SmartIO buffers writes
// until Persist while List* reads disk, so a project's first build would
// otherwise render a new.go calling nothing.
func CollectConstructorPackages(sandbox *api.Sandbox, io *smartio.SmartIO, constructors []Constructor) []string {

	seen := map[string]bool{}
	var packages []string

	add := func(name string) {
		if name == "" || seen[name] {
			return
		}
		seen[name] = true
		packages = append(packages, name)
	}

	for _, dir := range io.ListDirs(utils.ConstructorsDir) {
		parts := sandbox.Deps.Stringsdeps.Split(dir, "/")
		add(parts[len(parts)-1])
	}

	for _, constructor := range constructors {
		if constructor.HasNew {
			add(constructor.Package)
		}
	}

	sandbox.Deps.Sortdeps.Strings(packages)
	return packages
}
