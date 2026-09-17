package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/goimportsdeps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// CollectPublicApi parses every file of sandbox/api through the Go parser dep
// and returns one rich data map per file, for the {{range .PublicApi}} loop in
// the generated docs/PublicApi/doc.md. sandbox.go comes first (it is the root
// struct, one field per contract), the rest in listing order. Only exported
// declarations are kept: what is unexported is unreachable from a caller, so
// it is not public api.
func CollectPublicApi(sandbox *api.Sandbox, io *smartio.SmartIO) ([]map[string]any, error) {

	files := goFilesOf(sandbox, io, "sandbox/api")
	isRoot := func(file string) bool { return lastSegmentOf(sandbox, file) == "sandbox.go" }
	sandbox.Deps.Sortdeps.SliceStable(files, func(i int, j int) bool {
		return isRoot(files[i]) && !isRoot(files[j])
	})

	var groups []map[string]any
	for _, file := range files {
		parsed, err := parseGoFile(sandbox, io, file)
		if err != nil {
			return nil, err
		}

		group := fileData(sandbox, file, parsed)
		group["Name"] = titleOf(sandbox, sandbox.Deps.Stringsdeps.TrimSuffix(lastSegmentOf(sandbox, file), ".go"))
		group["Page"] = PublicApiPageOf(sandbox, file)
		group["Symbols"] = identifierList(sandbox, declaredSymbols(group))
		groups = append(groups, group)
	}

	return groups, nil
}

// PublicApiPageOf is the file one contract of sandbox/api is documented in,
// named after the source itself ("sandbox/api/actions.go" -> "api.actions.md")
// and relative to the docs/PublicApi directory that indexes it. The index link
// and the generated page are spelled here and nowhere else.
func PublicApiPageOf(sandbox *api.Sandbox, file string) string {
	base := sandbox.Deps.Stringsdeps.TrimSuffix(lastSegmentOf(sandbox, file), ".go")
	return "api." + base + ".md"
}

// DepsApiPageOf is PublicApiPageOf for a dependency contract, named after its
// directory ("stringsdeps" -> "deps.stringsdeps.md").
func DepsApiPageOf(name string) string {
	return "deps." + name + ".md"
}

// declaredSymbols is every exported name one parsed file declares, in the order
// its page prints them. docs/PublicApi's index lists them per page, so which
// page declares a symbol is answered by the index alone.
func declaredSymbols(data map[string]any) []string {
	var names []string
	for _, key := range []string{"Constants", "Variables", "Types", "Functions"} {
		entries, _ := data[key].([]map[string]any)
		for _, entry := range entries {
			if name, _ := entry["Name"].(string); name != "" {
				names = append(names, name)
			}
		}
	}
	return names
}

// parseGoFile reads a project file through the transaction-aware io and parses
// it with the Go parser dep. An unreadable or unparsable file is a hard error:
// `verify` reports both as violations, so a build that reaches here works on
// well-formed sources.
func parseGoFile(sandbox *api.Sandbox, io *smartio.SmartIO, file string) (*goimportsdeps.File, error) {
	content, err := io.ReadFile(file)
	if err != nil {
		return nil, sandbox.Deps.Std.Errorf("could not read %s: %w", file, err)
	}

	parsed, err := sandbox.Deps.Goimportsdeps.Parse(string(content))
	if err != nil {
		return nil, sandbox.Deps.Std.Errorf("could not parse %s: %w", file, err)
	}

	return parsed, nil
}

// fileData converts one parsed file into the map the doc template ranges over.
// Every doc comment is flattened to a single markdown-safe line, because each
// of them ends up inside a table cell.
func fileData(sandbox *api.Sandbox, file string, parsed *goimportsdeps.File) map[string]any {
	return map[string]any{
		"Path":      file,
		"Package":   parsed.Package,
		"Doc":       docLine(sandbox, parsed.Doc),
		"Types":     typesData(sandbox, parsed.Types),
		"Constants": valuesData(sandbox, parsed.Constants),
		"Variables": valuesData(sandbox, parsed.Variables),
		"Functions": functionsData(sandbox, parsed.Functions),
	}
}

// typesData keeps the exported types and describes each one by its kind: a
// struct renders as a field table, an interface as a method table, anything
// else as its underlying expression.
func typesData(sandbox *api.Sandbox, types []goimportsdeps.Type) []map[string]any {
	data := make([]map[string]any, 0, len(types))
	for _, entry := range types {
		if !entry.Exported {
			continue
		}
		fields := fieldsData(sandbox, entry.Fields)
		methods := functionsData(sandbox, entry.Methods)
		data = append(data, map[string]any{
			"Name":        entry.Name,
			"Doc":         docLine(sandbox, entry.Doc),
			"Kind":        entry.Kind,
			"Underlying":  entry.Underlying,
			"IsStruct":    entry.Kind == "struct",
			"IsInterface": entry.Kind == "interface",
			"Fields":      fields,
			"Methods":     methods,
			// A struct whose fields carry no comment of their own is
			// described by the type's own comment, so the doc drops the
			// column instead of printing a blank one per field.
			"FieldsDocumented":  anyDocumented(fields),
			"MethodsDocumented": anyDocumented(methods),
		})
	}
	return data
}

// fieldsData keeps the exported fields of a struct. An embedded field has no
// name of its own, so its type expression is used as its label.
func fieldsData(sandbox *api.Sandbox, fields []goimportsdeps.Field) []map[string]any {
	data := make([]map[string]any, 0, len(fields))
	for _, field := range fields {
		if !field.Exported {
			continue
		}
		name := field.Name
		if name == "" {
			name = field.Type
		}
		data = append(data, map[string]any{
			"Name": name,
			"Type": field.Type,
			"Doc":  docLine(sandbox, field.Doc),
		})
	}
	return data
}

// functionsData keeps the exported functions, methods and interface methods,
// each with the signature it is called by.
func functionsData(sandbox *api.Sandbox, functions []goimportsdeps.Function) []map[string]any {
	data := make([]map[string]any, 0, len(functions))
	for _, function := range functions {
		if !function.Exported {
			continue
		}
		data = append(data, map[string]any{
			"Name":      function.Name,
			"Doc":       docLine(sandbox, function.Doc),
			"Receiver":  function.Receiver,
			"Signature": functionSignature(sandbox, function),
		})
	}
	return data
}

// valuesData keeps the exported names of the const and var blocks, with the
// literal each one is assigned.
func valuesData(sandbox *api.Sandbox, values []goimportsdeps.Value) []map[string]any {
	data := make([]map[string]any, 0, len(values))
	for _, value := range values {
		if !value.Exported {
			continue
		}
		data = append(data, map[string]any{
			"Name":  value.Name,
			"Doc":   docLine(sandbox, value.Doc),
			"Type":  value.Type,
			"Value": value.Value,
		})
	}
	return data
}

// anyDocumented reports whether at least one of the entries carries a doc
// comment.
func anyDocumented(entries []map[string]any) bool {
	for _, entry := range entries {
		if entry["Doc"] != "" {
			return true
		}
	}
	return false
}

// functionSignature renders a parsed function back to the shape a caller
// writes: `Name(param type, param type) (result, result)`.
func functionSignature(sandbox *api.Sandbox, function goimportsdeps.Function) string {
	builder := function.Name + "(" + sandbox.Deps.Stringsdeps.Join(paramList(function.Params), ", ") + ")"

	results := paramList(function.Results)
	switch len(results) {
	case 0:
	case 1:
		builder += " " + results[0]
	default:
		builder += " (" + sandbox.Deps.Stringsdeps.Join(results, ", ") + ")"
	}

	return builder
}

// paramList renders one parameter or result list, keeping the names when the
// source declares them.
func paramList(params []goimportsdeps.Param) []string {
	rendered := make([]string, 0, len(params))
	for _, param := range params {
		if param.Name == "" {
			rendered = append(rendered, param.Type)
			continue
		}
		rendered = append(rendered, param.Name+" "+param.Type)
	}
	return rendered
}

// docLine flattens a doc comment to one markdown table cell: no line breaks,
// no collapsed run of spaces, and no bare pipe to break the column.
func docLine(sandbox *api.Sandbox, doc string) string {
	line := sandbox.Deps.Stringsdeps.Join(sandbox.Deps.Stringsdeps.Fields(doc), " ")
	return sandbox.Deps.Stringsdeps.ReplaceAll(line, "|", "\\|")
}

// goFilesOf lists the .go files directly inside dir, in listing order.
func goFilesOf(sandbox *api.Sandbox, io *smartio.SmartIO, dir string) []string {
	var files []string
	for _, file := range io.ListFiles(dir) {
		if sandbox.Deps.Stringsdeps.HasSuffix(file, ".go") {
			files = append(files, file)
		}
	}
	return files
}

// titleOf upper-cases the first letter of name, the same spelling rule every
// other collector uses.
func titleOf(sandbox *api.Sandbox, name string) string {
	if len(name) == 0 {
		return name
	}
	return sandbox.Deps.Stringsdeps.ToUpper(name[:1]) + name[1:]
}
