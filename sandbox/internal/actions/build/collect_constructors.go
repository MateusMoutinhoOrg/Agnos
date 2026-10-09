package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// Constructor is one contract of sandbox/api/ that the sandbox carries a field
// for: the field and its type on sandbox/api/generated.sandbox.go, and the package under
// sandbox/internal/ whose New<Name> fills it in sandbox/generated.new.go.
type Constructor struct {
	// Name is the contract title-cased ("Cli"), both the field and its type.
	Name string
	// Package is the sandbox/internal/ directory that builds it ("cli").
	Package string
	// Source is the project-relative package whose New<Name> builds it —
	// sandbox/internal/<Package>, or the one level down
	// utils.ConstructorSource finds.
	Source string
	// HasNew reports that sandbox/internal/<Package>/ declares its New<Name>
	// — in a generated.new.go or a hand-written new.go — to be called. A
	// contract with none is a field of the Sandbox that nothing fills — an api published to be installed elsewhere, say — so
	// sandbox/generated.new.go leaves it alone rather than naming a package that is
	// not written yet.
	HasNew bool
}

// CollectConstructors lists sandbox/api and returns one entry per .go file
// other than sandbox.go — each named with utils.GeneratedPrefix taken off, so
// generated.cli.go is the Cli contract — for the {{range .Constructors}} loops
// in sandbox/api/generated.sandbox.go and sandbox/generated.new.go. command.go and route.go are
// skipped along with it: what they declare is the shape of one command and of
// one route, not a field of the sandbox — each belongs to the contract whose
// New<Name> builds the slice of them.
//
// generated.sandbox.go takes every entry — the field is the contract, whether
// or not this repo fills it — while generated.new.go takes the ones HasNew
// marks.
func CollectConstructors(sandbox *api.Sandbox, io *stagedfs.StagedFS) []Constructor {

	files := utils.DropSuperseded(sandbox, io.ListFiles("sandbox/api"))

	var constructors []Constructor
	for _, file := range files {
		parts := sandbox.Deps.StringsDeps.Split(file, "/")
		name := utils.SourceName(sandbox, parts[len(parts)-1])

		if utils.IsConstructorExempt(sandbox, name) || !sandbox.Deps.StringsDeps.HasSuffix(name, ".go") {
			continue
		}

		baseName := sandbox.Deps.StringsDeps.TrimSuffix(name, ".go")
		if len(baseName) == 0 {
			continue
		}

		title := sandbox.Deps.StringsDeps.ToUpper(baseName[:1]) + baseName[1:]
		source := utils.ConstructorSource(sandbox, io, baseName)
		constructors = append(constructors, Constructor{
			Name:    title,
			Package: baseName,
			Source:  source,
			HasNew:  utils.ConstructorNewFile(sandbox, io, source) != "",
		})
	}

	return constructors
}
