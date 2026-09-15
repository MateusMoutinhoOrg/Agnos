package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/templatedeps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// templateFuncs is the FuncMap every template rendered through this package is
// given. It carries one native function, `render`, so a template can embed the
// rendering of another template inside itself:
//
//	{{ render "AgnosConfig/docs/ReadmeHeader.md" }}
//
// `render` takes a project-relative path, reads that file through the
// transaction-aware io (so a file written earlier in the same build is visible),
// renders it as a Go text/template with the *same* vars and the same FuncMap
// (so nesting works to any depth), and returns the result as a string. A
// missing or unparsable file is a hard error, surfaced through the template
// execution error like any other.
//
// The second native function, `copy`, is the raw counterpart of `render`: it
// reads a project-relative file through the same transaction-aware io and
// returns its contents verbatim, without any template rendering. Use it to
// embed a file that is not itself a template (a LICENSE, a fixed snippet). A
// missing file is a hard error.
func templateFuncs(sandbox *api.Sandbox, io *smartio.SmartIO, vars interface{}) map[string]any {
	return map[string]any{
		"render": func(project_path string) (string, error) {
			content, err := io.ReadFile(project_path)
			if err != nil {
				return "", err
			}
			rendered, err := renderTemplate(sandbox, io, baseName(sandbox, project_path), content, vars)
			if err != nil {
				return "", err
			}
			return string(rendered), nil
		},
		"copy": func(project_path string) (string, error) {
			content, err := io.ReadFile(project_path)
			if err != nil {
				return "", err
			}
			return string(content), nil
		},
	}
}

// RenderTemplate renders one template source over vars and returns the result
// without writing anything. It is what a caller needs when it has to read a
// value out of an asset the build is about to write — the doc props of a
// generated doc, say — rather than render that asset to its destination.
func RenderTemplate(sandbox *api.Sandbox, io *smartio.SmartIO, name string, src []byte, vars interface{}) ([]byte, error) {
	return renderTemplate(sandbox, io, name, src, vars)
}

// renderTemplate parses src as a Go text/template named name and executes it
// over vars, with templateFuncs available. It is the single rendering path for
// this package: RenderTemplateToDest and RenderGroup both go through it, so the
// `render` native function is available in every asset template.
func renderTemplate(sandbox *api.Sandbox, io *smartio.SmartIO, name string, src []byte, vars interface{}) ([]byte, error) {
	rendered, err := sandbox.Deps.Templatedeps.Render(templatedeps.RenderProps{
		Name:   name,
		Source: string(src),
		Vars:   vars,
		Funcs:  templateFuncs(sandbox, io, vars),
	})
	if err != nil {
		return nil, err
	}
	return []byte(rendered), nil
}

// baseName is the last slash-separated segment of an asset path, the name a
// template is reported under when it fails to parse or execute.
func baseName(sandbox *api.Sandbox, path string) string {
	segments := sandbox.Deps.Stringsdeps.Split(path, "/")
	return segments[len(segments)-1]
}

// formatIfGo returns content in the canonical form of the Go toolchain when
// dest is a Go file, and content untouched otherwise. Every generated .go file
// passes through here so a regenerated tree diffs to zero against one a
// formatting editor has saved. A render that is not parsable Go is written as
// it came out: the compile step reports it with a real message, which a
// rendering error here would only hide.
func formatIfGo(sandbox *api.Sandbox, dest string, content []byte) []byte {
	if !sandbox.Deps.Stringsdeps.HasSuffix(dest, ".go") {
		return content
	}

	formatted, err := sandbox.Deps.Goimportsdeps.Format(string(content))
	if err != nil {
		sandbox.Deps.Std.Log("could not format %s: %s\n", dest, err.Error())
		return content
	}

	return []byte(formatted)
}

// RenderTemplateToDest renders one asset as a Go text/template over vars and
// writes the result to dest_path. A Go destination is formatted first (see
// formatIfGo).
func RenderTemplateToDest(sandbox *api.Sandbox, io *smartio.SmartIO, template_path string, vars interface{}, dest_path string) error {

	src, err := sandbox.Deps.Embeddeps.ReadFile(template_path)
	if err != nil {
		return err
	}

	content, err := renderTemplate(sandbox, io, baseName(sandbox, template_path), src, vars)
	if err != nil {
		return err
	}

	err = io.WriteFileOverwrite(dest_path, formatIfGo(sandbox, dest_path, content))
	if err != nil {
		return err
	}

	return nil
}

// RenderGroup renders every asset under assets/<group> as a Go text/template
// and writes each result to the path it holds inside the group. An asset at
// assets/all/sandbox/new.go rendered with RenderGroup(deps, io, "all", vars)
// is written to sandbox/new.go. Every file in the group is rendered with the
// same vars, and every file may call the `render` native function (see
// templateFuncs) to embed another template of the target project.
func RenderGroup(sandbox *api.Sandbox, io *smartio.SmartIO, group string, vars interface{}) error {
	return RenderGroupExcept(sandbox, io, group, vars, nil)
}

// RenderGroupExcept renders a group the way RenderGroup does, skipping every
// asset whose group-relative path is listed in except. The dep and adapter
// catalogs carry their declaration at the root of the group — a file that
// describes the package rather than one the package installs — so it is
// rendered nowhere.
func RenderGroupExcept(sandbox *api.Sandbox, io *smartio.SmartIO, group string, vars interface{}, except []string) error {

	files, err := sandbox.Deps.Embeddeps.ListFilesRecursively(group)
	if err != nil {
		return err
	}

	for _, file := range files {
		if containsPath(except, file) {
			continue
		}
		src, err := sandbox.Deps.Embeddeps.ReadFile(group + "/" + file)
		if err != nil {
			return err
		}

		content, err := renderTemplate(sandbox, io, baseName(sandbox, file), src, vars)
		if err != nil {
			return err
		}

		err = io.WriteFileOverwrite(file, formatIfGo(sandbox, file, content))
		if err != nil {
			return err
		}
	}

	return nil
}

// containsPath reports whether paths holds path, the membership test
// RenderGroupExcept skips an asset with.
func containsPath(paths []string, path string) bool {
	for _, candidate := range paths {
		if candidate == path {
			return true
		}
	}
	return false
}
