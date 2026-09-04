package utils

import (
	"bytes"
	"path"
	"text/template"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
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
func templateFuncs(deps *deps.Deps, io *smartio.SmartIO, vars interface{}) template.FuncMap {
	return template.FuncMap{
		"render": func(project_path string) (string, error) {
			content, err := io.ReadFile(project_path)
			if err != nil {
				return "", err
			}
			rendered, err := renderTemplate(deps, io, path.Base(project_path), content, vars)
			if err != nil {
				return "", err
			}
			return string(rendered), nil
		},
	}
}

// renderTemplate parses src as a Go text/template named name and executes it
// over vars, with templateFuncs available. It is the single rendering path for
// this package: RenderTemplateToDest and RenderGroup both go through it, so the
// `render` native function is available in every asset template.
func renderTemplate(deps *deps.Deps, io *smartio.SmartIO, name string, src []byte, vars interface{}) ([]byte, error) {
	t, err := template.New(name).Funcs(templateFuncs(deps, io, vars)).Parse(string(src))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, vars); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func RenderTemplateToDest(deps *deps.Deps, io *smartio.SmartIO, template_path string, vars interface{}, dest_path string) error {

	src, err := deps.Embeddeps.ReadFile(template_path)
	if err != nil {
		return err
	}

	content, err := renderTemplate(deps, io, path.Base(template_path), src, vars)
	if err != nil {
		return err
	}

	err = io.WriteFileOverwrite(dest_path, content)
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
func RenderGroup(deps *deps.Deps, io *smartio.SmartIO, group string, vars interface{}) error {

	files, err := deps.Embeddeps.ListFilesRecursively(group)
	if err != nil {
		return err
	}

	for _, file := range files {
		src, err := deps.Embeddeps.ReadFile(group + "/" + file)
		if err != nil {
			return err
		}

		content, err := renderTemplate(deps, io, path.Base(file), src, vars)
		if err != nil {
			return err
		}

		err = io.WriteFileOverwrite(file, content)
		if err != nil {
			return err
		}
	}

	return nil
}
