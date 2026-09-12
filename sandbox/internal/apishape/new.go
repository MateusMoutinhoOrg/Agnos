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

	shape := &Api{ByName: map[string]goimportsdeps.Type{}}

	for _, source := range sources {
		parsed, err := sandbox.Deps.Goimportsdeps.Parse(source.Content)
		if err != nil {
			return nil, sandbox.Deps.Std.Errorf("%s is not parsable Go: %w", source.Name, err)
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

	return shape, nil
}
