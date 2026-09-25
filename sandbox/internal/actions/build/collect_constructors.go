package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// Constructor is one contract of sandbox/api/ that the sandbox carries a field
// for: the field and its type on sandbox/api/sandbox.go, and the package under
// sandbox/internal/ whose New<Name> fills it in sandbox/new.go.
type Constructor struct {
	// Name is the contract title-cased ("Cli"), both the field and its type.
	Name string
	// Package is the sandbox/internal/ directory that builds it ("cli").
	Package string
	// Source is the project-relative package whose New<Name> builds it —
	// sandbox/internal/<Package>, or the one level down
	// utils.ConstructorSource finds.
	Source string
	// HasNew reports that sandbox/internal/<Package>/new.go is there to be
	// called. A contract with none is a field of the Sandbox that nothing
	// fills — an api published to be installed elsewhere, say — so
	// sandbox/new.go leaves it alone rather than naming a package that is
	// not written yet.
	HasNew bool
}

// CollectConstructors lists sandbox/api and returns one entry per .go file
// other than sandbox.go, for the {{range .Constructors}} loops in
// sandbox/api/sandbox.go and sandbox/new.go. command.go and route.go are
// skipped along with it: what they declare is the shape of one command and of
// one route, not a field of the sandbox — each belongs to the contract whose
// New<Name> builds the slice of them.
//
// sandbox/api/sandbox.go takes every entry — the field is the contract, whether
// or not this repo fills it — while sandbox/new.go takes the ones HasNew marks.
func CollectConstructors(sandbox *api.Sandbox, io *smartio.SmartIO) []Constructor {

	files := io.ListFiles("sandbox/api")

	var constructors []Constructor
	for _, file := range files {
		parts := sandbox.Deps.Stringsdeps.Split(file, "/")
		name := parts[len(parts)-1]

		if name == "sandbox.go" || name == "command.go" || name == "route.go" || !sandbox.Deps.Stringsdeps.HasSuffix(name, ".go") {
			continue
		}

		baseName := sandbox.Deps.Stringsdeps.TrimSuffix(name, ".go")
		if len(baseName) == 0 {
			continue
		}

		title := sandbox.Deps.Stringsdeps.ToUpper(baseName[:1]) + baseName[1:]
		source := utils.ConstructorSource(io, baseName)
		constructors = append(constructors, Constructor{
			Name:    title,
			Package: baseName,
			Source:  source,
			HasNew:  io.IsFile(source + "/new.go"),
		})
	}

	return constructors
}
