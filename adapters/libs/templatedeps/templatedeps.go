package templatedeps

import (
	"bytes"
	"text/template"

	templatedeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/templatedeps"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

// render fills templatedeps.Lib.Render, parsing the source as a Go
// text/template with the given native functions registered and executing it
// over the given vars.
func render(props templatedeps.RenderProps) (string, error) {
	parsed, err := template.New(props.Name).Funcs(template.FuncMap(props.Funcs)).Parse(props.Source)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	if err := parsed.Execute(&buffer, props.Vars); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// Bind fills deps.Deps.Templatedeps with the standard library's text/template.
func Bind(deps *deps.Deps) {
	deps.Templatedeps = templatedeps.Lib{
		Render: func(props templatedeps.RenderProps) (string, error) {
			return render(props)
		},
	}
}
