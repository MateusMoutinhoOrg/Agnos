package apishape

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	goimportsdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/goimportsdeps"
)

// New parses one sandbox/api/ package. A file that does not parse is a hard
// error: everything downstream — the shape rule, the copy, the shim — reads
// the declarations, so there is nothing to say about a package half of which
// is unreadable.
func New(sandbox *api.Sandbox, sources []Source) (*Api, error) {

	shape := &Api{ByName: map[string]goimportsdeps.Type{}, Mechanic: map[string]bool{}}

	for _, source := range sources {
		parsed, err := sandbox.Deps.GoimportsDeps.Parse(source.Content)
		if err != nil {
			return nil, sandbox.Deps.StdDeps.Errorf("%s is not parsable Go: %w", source.Name, err)
		}

		if shape.Package == "" {
			shape.Package = parsed.Package
		}

		shape.Files = append(shape.Files, ParsedFile{
			Name:    source.Name,
			Content: source.Content,
			Parsed:  parsed,
		})

		for _, entry := range parsed.Types {
			shape.Types = append(shape.Types, entry)
			shape.ByName[entry.Name] = entry
		}
	}

	markMechanic(sandbox, shape)

	return shape, nil
}

// markMechanic fills Api.Mechanic: first every alias whose target is
// qualified by an opinionated lib's contract import, then — until nothing more
// is found — every struct whose fields are all mechanic, which is how a part
// holding nothing but a mechanic's surface joins it.
func markMechanic(sandbox *api.Sandbox, shape *Api) {
	for _, file := range shape.Files {
		qualifiers := opinionatedQualifiers(sandbox, file.Parsed.Imports)
		for _, entry := range file.Parsed.Types {
			if entry.Kind != "alias" {
				continue
			}
			parts := sandbox.Deps.StringsDeps.Split(entry.Underlying, ".")
			if len(parts) == 2 && qualifiers[parts[0]] {
				shape.Mechanic[entry.Name] = true
			}
		}
	}

	for found := true; found; {
		found = false
		for _, entry := range shape.Types {
			if entry.Kind != "struct" || len(entry.Fields) == 0 || shape.Mechanic[entry.Name] {
				continue
			}
			mechanic := true
			for _, field := range entry.Fields {
				if !shape.Mechanic[field.Type] {
					mechanic = false
					break
				}
			}
			if mechanic {
				shape.Mechanic[entry.Name] = true
				found = true
			}
		}
	}
}

// opinionatedQualifiers is the set of names one file refers to an opinionated lib's
// contract by: its explicit alias, or the lower-cased last segment of its
// path, which is the package clause every such contract is written with.
func opinionatedQualifiers(sandbox *api.Sandbox, imports []goimportsdeps.Import) map[string]bool {
	qualifiers := map[string]bool{}
	for _, imp := range imports {
		segments := sandbox.Deps.StringsDeps.Split(imp.Path, "/")
		last := segments[len(segments)-1]
		if !sandbox.Deps.StringsDeps.HasPrefix(last, OpinionatedPrefix) {
			continue
		}
		if imp.Alias != "" {
			qualifiers[imp.Alias] = true
			continue
		}
		qualifiers[sandbox.Deps.StringsDeps.ToLower(last)] = true
	}
	return qualifiers
}
