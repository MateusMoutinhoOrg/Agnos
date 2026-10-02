package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// EmbeddedStruct is one struct a generated aggregate embeds: api.Config embeds
// every struct of sandbox/api/userconfig*.go, api.Sandbox every one of
// sandbox/api/usersandbox*.go, and RouteProps and CommandProps every one their
// own package declares beside them. Each part is a file of its own, so a
// mechanic adds what it needs by adding a file — never by editing one the
// project already wrote in.
type EmbeddedStruct struct {
	// Name is the struct's type name, the embedded field's name too.
	Name string
	// File is the project-relative file that declares it.
	File string
	// Fields are the names of the struct's own fields — the ones the
	// aggregate promotes, so two parts may not share one.
	Fields []string
}

// CollectEmbeddedStructs parses every .go file directly under dir that accept
// takes — by its base name — and returns each exported struct they declare,
// sorted by name so the aggregate renders the same bytes whatever order the
// disk lists them in.
//
// pending are files of dir written earlier in this same transaction: the
// listing reads disk, so a file a migration just moved would be missed
// without them.
func CollectEmbeddedStructs(sandbox *api.Sandbox, io *smartio.SmartIO, dir string, accept func(name string) bool, pending []string) ([]EmbeddedStruct, error) {
	files := io.ListFiles(dir)
	for _, file := range pending {
		if !containsPath(files, file) {
			if _, err := io.ReadFile(file); err == nil {
				files = append(files, file)
			}
		}
	}

	var structs []EmbeddedStruct
	for _, file := range files {
		name := baseName(sandbox, file)
		if !sandbox.Deps.Stringsdeps.HasSuffix(name, ".go") || !accept(name) {
			continue
		}

		content, err := io.ReadFile(file)
		if err != nil {
			return nil, err
		}
		parsed, err := sandbox.Deps.Goimportsdeps.Parse(string(content))
		if err != nil {
			return nil, sandbox.Deps.Std.Errorf("%s is not parsable Go: %s", file, err.Error())
		}

		for _, entry := range parsed.Types {
			if entry.Kind != "struct" || !entry.Exported {
				continue
			}
			var fields []string
			for _, field := range entry.Fields {
				if field.Name != "" {
					fields = append(fields, field.Name)
				}
			}
			structs = append(structs, EmbeddedStruct{Name: entry.Name, File: file, Fields: fields})
		}
	}

	sandbox.Deps.Sortdeps.Slice(structs, func(i int, j int) bool {
		return structs[i].Name < structs[j].Name
	})
	return structs, nil
}

// EmbeddedCollisions returns one line per field name two parts of one
// aggregate both declare, or that a part shares with a field the aggregate
// declares itself (own). Go only refuses such an ambiguous selector where it
// is read, so without this the project compiles until the day a handler
// reads it.
func EmbeddedCollisions(sandbox *api.Sandbox, aggregate string, structs []EmbeddedStruct, own []string) []string {
	owner := map[string]string{}
	for _, name := range own {
		owner[name] = aggregate
	}

	var collisions []string
	for _, part := range structs {
		if _, taken := owner[part.Name]; taken {
			collisions = append(collisions, part.File+" declares "+part.Name+", a name "+owner[part.Name]+" already holds; "+aggregate+" embeds it under that name")
		}
		owner[part.Name] = part.File
		for _, field := range part.Fields {
			if previous, taken := owner[field]; taken {
				collisions = append(collisions, part.File+" declares field "+field+", which "+previous+" declares too; "+aggregate+" embeds both, so reading it is ambiguous")
				continue
			}
			owner[field] = part.File
		}
	}
	return collisions
}

// HasFilePrefix returns an accept func for CollectEmbeddedStructs taking the
// files whose base name starts with prefix: userconfig.go and
// userconfig_backoffice.go for "userconfig".
func HasFilePrefix(sandbox *api.Sandbox, prefix string) func(name string) bool {
	return func(name string) bool {
		return sandbox.Deps.Stringsdeps.HasPrefix(name, prefix)
	}
}

// AllBut returns an accept func for CollectEmbeddedStructs taking every file
// except the aggregate itself.
func AllBut(aggregate string) func(name string) bool {
	return func(name string) bool {
		return name != aggregate
	}
}
