package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

func RenderTemplateToDest(deps *deps.Deps, io *smartio.SmartIO, template_path string, vars interface{}, dest_path string) error {

	content, err := deps.Embeddeps.RenderTemplate(template_path, vars)
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
// same vars.
func RenderGroup(deps *deps.Deps, io *smartio.SmartIO, group string, vars interface{}) error {

	files, err := deps.Embeddeps.ListFilesRecursively(group)
	if err != nil {
		return err
	}

	for _, file := range files {
		content, err := deps.Embeddeps.RenderTemplate(group+"/"+file, vars)
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
