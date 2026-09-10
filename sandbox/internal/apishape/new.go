package apishape

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	goimportsdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/goimportsdeps"
)

// New parses one sandbox/api/ package. A file that does not parse is a hard
// error: everything downstream — the shape rule, the copy, the shim — reads
// the declarations, so there is nothing to say about a package half of which
// is unreadable.
func New(deps *deps.Deps, sources []Source) (*Api, error) {

	api := &Api{ByName: map[string]goimportsdeps.Type{}}

	for _, source := range sources {
		parsed, err := deps.Goimportsdeps.Parse(source.Content)
		if err != nil {
			return nil, deps.Std.Errorf("%s is not parsable Go: %w", source.Name, err)
		}

		if api.Package == "" {
			api.Package = parsed.Package
		}

		api.Files = append(api.Files, ParsedFile{
			Name:    source.Name,
			Content: source.Content,
			Parsed:  parsed,
		})

		for _, entry := range parsed.Types {
			api.Types = append(api.Types, entry)
			api.ByName[entry.Name] = entry
		}
	}

	return api, nil
}
